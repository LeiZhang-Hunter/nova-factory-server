package iotdb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var tsDbOnce sync.Once
var tsDbInstance *TsdbStorage

// TSDBStorage 定义时序数据库存储，提供写入和查询能力。
type TSDBStorage interface {
	Appendable
	Queryable

	// Close 关闭存储以及其底层资源。
	Close() error
}

// Appendable 用于创建批量写入器。
type Appendable interface {
	// Appender 返回一个新的存储写入器。
	Appender() Appender
}

// Appender 用于向存储批量写入 MetricSample。
type Appender interface {
	// Append 写入一组 MetricSample。
	Append(s []MetricSample) error
	// Commit 提交本批次 MetricSample 并清理批次状态。
	// 如果 Commit 返回错误，会回滚当前写入器已做出的修改。
	// Commit 调用后，无论成功或失败，都不能继续使用该 Appender。
	Commit() error
}

// Querier 提供固定时间范围内的时序数据查询能力。
type Querier interface {
	// Query 查询匹配指定元数据的序列，并写入 MetricResult。
	// 查询提示可用于传递查询优化信息，具体是否使用由实现决定。
	Query(meta MetricMeta, hints *QueryHints, result MetricResult) error

	// QueryAndClose 与 Query 行为一致，但会在返回前关闭查询器。
	QueryAndClose(meta MetricMeta, hints *QueryHints, result MetricResult) error

	// Close 关闭查询器并释放其资源。
	Close()
}

// Queryable 表示可被查询的存储。
type Queryable interface {
	// Querier 返回指定时间范围内的查询器。
	Querier(startTime, endTime time.Time) (Querier, error)
}

var (
	errTSDBStorageClosed = errors.New("prometheus tsdb storage is closed")
	errTSDBInvalidMetric = errors.New("invalid metric")
)

// tsdbConfig 将本地 `tsdb` 配置段映射为 Prometheus TSDB 选项。
// 时间字段使用 Go time.Duration，因此配置值可以写成 "24h" 这类格式。
type tsdbConfig struct {
	// TSDBPath 是 Prometheus TSDB 存放 WAL 和数据块的目录。
	TSDBPath                      string        `mapstructure:"path"`
	TSDBRetentionDuration         time.Duration `mapstructure:"retention_duration"`
	TSDBStripeSize                int           `mapstructure:"stripe_size"`
	TSDBMaxBytes                  int64         `mapstructure:"max_bytes"`
	TSDBWALSegmentSize            int           `mapstructure:"wal_segment_size"`
	TSDBMaxBlockChunkSegmentSize  int64         `mapstructure:"max_block_chunk_segment_size"`
	TSDBMinBlockDuration          time.Duration `mapstructure:"min_block_duration"`
	TSDBMaxBlockDuration          time.Duration `mapstructure:"max_block_duration"`
	TSDBHeadChunksWriteBufferSize int           `mapstructure:"head_chunks_write_buffer_size"`
	TSDBEnablePromMetrics         bool          `mapstructure:"enable_prom_metrics"`
}

// defaultTSDBConfig 保留 Prometheus 默认值，并补充本项目默认数据目录。
func defaultTSDBConfig() tsdbConfig {
	opt := tsdb.DefaultOptions()
	return tsdbConfig{
		TSDBPath:                      filepath.Join("data", "tsdb"),
		TSDBRetentionDuration:         time.Duration(opt.RetentionDuration) * time.Millisecond,
		TSDBStripeSize:                opt.StripeSize,
		TSDBMaxBytes:                  opt.MaxBytes,
		TSDBWALSegmentSize:            opt.WALSegmentSize,
		TSDBMaxBlockChunkSegmentSize:  opt.MaxBlockChunkSegmentSize,
		TSDBMinBlockDuration:          time.Duration(opt.MinBlockDuration) * time.Millisecond,
		TSDBMaxBlockDuration:          time.Duration(opt.MaxBlockDuration) * time.Millisecond,
		TSDBHeadChunksWriteBufferSize: opt.HeadChunksWriteBufferSize,
	}
}

// loadTSDBConfig 在 Prometheus 默认配置之上叠加 viper 配置。
func loadTSDBConfig() tsdbConfig {
	conf := defaultTSDBConfig()
	if err := viper.UnmarshalKey("tsdb", &conf); err != nil {
		zap.L().Warn("unmarshal tsdb config failed", zap.Error(err))
	}
	return conf
}

// TsdbStorage 实现 TSDBStorage。
type TsdbStorage struct {
	db *tsdb.DB
}

func GetTSDBStorage() *TsdbStorage {
	tsDbOnce.Do(func() {
		tsDbInstance = newTSDBStorage()
	})
	return tsDbInstance
}

// newTSDBStorage 打开一个本地 Prometheus TSDB 实例。
// 打开失败时直接触发 panic，以保持和项目内其他数据源构造函数一致。
func newTSDBStorage() *TsdbStorage {
	conf := loadTSDBConfig()
	tsdbOpt := tsdb.DefaultOptions()
	tsdbOpt.RetentionDuration = int64(conf.TSDBRetentionDuration / time.Millisecond)
	tsdbOpt.StripeSize = conf.TSDBStripeSize
	tsdbOpt.MaxBytes = conf.TSDBMaxBytes
	tsdbOpt.WALSegmentSize = conf.TSDBWALSegmentSize
	tsdbOpt.MaxBlockChunkSegmentSize = conf.TSDBMaxBlockChunkSegmentSize
	tsdbOpt.MinBlockDuration = int64(conf.TSDBMinBlockDuration / time.Millisecond)
	tsdbOpt.MaxBlockDuration = int64(conf.TSDBMaxBlockDuration / time.Millisecond)
	tsdbOpt.HeadChunksWriteBufferSize = conf.TSDBHeadChunksWriteBufferSize
	// 避免使用 prometheus.tsdb v0.39 及以上版本时出现乱序写入冲突。
	// prometheus.tsdb v0.37 要求同一序列的样本严格按时间递增写入。
	// v0.39 起可通过 OutOfOrderTimeWindow 允许一定窗口内的乱序样本。
	// 该窗口需要和指标序列的采样粒度保持匹配。
	tsdbOpt.OutOfOrderTimeWindow = int64(time.Minute / time.Millisecond)
	zap.L().Debug("ready to start tsdb", zap.String("path", conf.TSDBPath), zap.Any("option", tsdbOpt))

	// 仅在显式启用时注册 Prometheus 自身的 TSDB 指标，避免多存储场景下重复注册采集器。
	var promReg prometheus.Registerer
	if conf.TSDBEnablePromMetrics {
		promReg = prometheus.DefaultRegisterer
	}

	db, err := tsdb.Open(conf.TSDBPath, slog.New(slog.NewTextHandler(io.Discard, nil)), promReg, tsdbOpt, nil)
	if err != nil {
		panic(fmt.Errorf("open prometheus tsdb: %w", err))
	}

	return &TsdbStorage{db: db}
}

// Appender 创建一个 Prometheus 写入事务。批次内所有样本写入后需要调用 Commit。
func (s *TsdbStorage) Appender() Appender {
	if s == nil || s.db == nil {
		return &tsdbAppender{err: errTSDBStorageClosed}
	}
	return &tsdbAppender{appender: s.db.Appender(context.Background())}
}

// Querier 创建固定时间范围的读取器。Prometheus TSDB 使用毫秒时间戳，因此在这里转换时间对象。
func (s *TsdbStorage) Querier(startTime, endTime time.Time) (Querier, error) {
	if s == nil || s.db == nil {
		return nil, errTSDBStorageClosed
	}

	start := timeToMillis(startTime)
	end := timeToMillis(endTime)
	querier, err := s.db.Querier(start, end)
	if err != nil {
		return nil, err
	}

	return &tsdbQuerier{
		querier: querier,
		start:   start,
		end:     end,
	}, nil
}

func (s *TsdbStorage) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

// Delete 按指标类型和 labels 删除 Prometheus TSDB 中匹配的全部序列数据。
func (s *TsdbStorage) Delete(meta MetricMeta) error {
	if s == nil || s.db == nil {
		return errTSDBStorageClosed
	}
	matchers, err := metricMatchers(meta)
	if err != nil {
		return err
	}
	return s.db.Delete(context.Background(), -1<<63, 1<<63-1, matchers...)
}

// tsdbAppender 将项目内的 MetricSample 接口适配到 Prometheus 的 storage.Appender 事务接口。
type tsdbAppender struct {
	appender storage.Appender
	err      error
}

// Append 将每个样本的指标元数据转换为标签，并写入时间戳和值。
// 任意写入错误都会回滚整个事务。
func (t *tsdbAppender) Append(samples []MetricSample) error {
	if t.err != nil {
		return t.err
	}
	if t.appender == nil {
		t.err = errTSDBStorageClosed
		return t.err
	}

	for _, sample := range samples {
		metricLabels, err := metricLabels(sample)
		if err != nil {
			_ = t.appender.Rollback()
			t.err = err
			return err
		}

		if _, err := t.appender.Append(0, metricLabels, sample.timestamp(), sample.value()); err != nil {
			_ = t.appender.Rollback()
			t.err = err
			return err
		}
	}

	return nil
}

// Commit 完成 Prometheus 写入事务。调用后必须丢弃该写入器，保持 Prometheus Appender 契约。
func (t *tsdbAppender) Commit() error {
	if t.err != nil {
		return t.err
	}
	if t.appender == nil {
		return errTSDBStorageClosed
	}
	return t.appender.Commit()
}

// tsdbQuerier 保存底层 Prometheus 查询器，以及作为 SelectHints 使用的毫秒级时间边界。
type tsdbQuerier struct {
	querier storage.Querier
	start   int64
	end     int64
}

// Query 选择匹配指标类型和全部指标属性的序列，并把结果整理交给 MetricResult。
func (q *tsdbQuerier) Query(meta MetricMeta, hints *QueryHints, result MetricResult) error {
	if q == nil || q.querier == nil {
		return errTSDBStorageClosed
	}
	if result == nil {
		return fmt.Errorf("%w: nil metric result", errTSDBInvalidMetric)
	}

	matchers, err := metricMatchers(meta)
	if err != nil {
		return err
	}

	seriesSet := q.querier.Select(context.Background(), true, &storage.SelectHints{
		Start: q.start,
		End:   q.end,
	}, matchers...)
	for seriesSet.Next() {
		if err := result.AddSeries(seriesSet.At()); err != nil {
			return err
		}
	}
	if err := seriesSet.Err(); err != nil {
		return err
	}
	if warnings := seriesSet.Warnings(); len(warnings) > 0 {
		zap.L().Warn("query prometheus tsdb with warnings", zap.Any("warnings", warnings))
	}

	return nil
}

func (q *tsdbQuerier) QueryAndClose(meta MetricMeta, hints *QueryHints, result MetricResult) error {
	defer q.Close()
	return q.Query(meta, hints, result)
}

func (q *tsdbQuerier) Close() {
	if q == nil || q.querier == nil {
		return
	}
	if err := q.querier.Close(); err != nil {
		zap.L().Warn("close prometheus tsdb querier failed", zap.Error(err))
	}
}

// metricLabels 将指标类型保存为 Prometheus 约定的 __name__ label，
// 并将 MetricMeta 的属性保存为普通标签。
func metricLabels(meta MetricMeta) (labels.Labels, error) {
	if meta == nil {
		return labels.EmptyLabels(), fmt.Errorf("%w: nil metric meta", errTSDBInvalidMetric)
	}
	if meta.GetKind() == "" {
		return labels.EmptyLabels(), fmt.Errorf("%w: empty metric kind", errTSDBInvalidMetric)
	}

	values := make(map[string]string, len(meta.GetProperties())+1)
	values[metricLabelName] = meta.GetKind()
	for key, value := range meta.GetProperties() {
		if key == "" {
			return labels.EmptyLabels(), fmt.Errorf("%w: empty label name", errTSDBInvalidMetric)
		}
		values[key] = value
	}

	return labels.FromMap(values), nil
}

// metricMatchers 与 metricLabels 保持一致，使用精确匹配器按指标类型和维度查询。
func metricMatchers(meta MetricMeta) ([]*labels.Matcher, error) {
	metricLabels, err := metricLabels(meta)
	if err != nil {
		return nil, err
	}

	var matcherErr error
	matchers := make([]*labels.Matcher, 0, metricLabels.Len())
	metricLabels.Range(func(label labels.Label) {
		if matcherErr != nil {
			return
		}

		matcher, err := labels.NewMatcher(labels.MatchEqual, label.Name, label.Value)
		if err != nil {
			matcherErr = err
			return
		}
		matchers = append(matchers, matcher)
	})
	if matcherErr != nil {
		return nil, fmt.Errorf("%w: %w", errTSDBInvalidMetric, matcherErr)
	}

	return matchers, nil
}

// timeToMillis 将 Go 时间转换为 Prometheus TSDB 使用的毫秒时间戳。
func timeToMillis(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
