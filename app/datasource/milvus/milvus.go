package milvus

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type Config struct {
	Enabled     bool          `mapstructure:"enabled"`
	Address     string        `mapstructure:"address"`
	Username    string        `mapstructure:"username"`
	Password    string        `mapstructure:"password"`
	DBName      string        `mapstructure:"db_name"`
	APIKey      string        `mapstructure:"api_key"`
	EnableTLS   bool          `mapstructure:"enable_tls"`
	DialTimeout time.Duration `mapstructure:"dial_timeout"`
}

var (
	client   *milvusclient.Client
	clientMu sync.RWMutex
)

func GetClient(ctx context.Context) (*milvusclient.Client, error) {
	clientMu.RLock()
	if client != nil {
		defer clientMu.RUnlock()
		return client, nil
	}
	clientMu.RUnlock()

	cfg, enabled, err := loadConfig()
	if err != nil {
		return nil, err
	}
	if !enabled {
		return nil, fmt.Errorf("milvus is not configured or disabled")
	}

	clientMu.Lock()
	defer clientMu.Unlock()
	if client != nil {
		return client, nil
	}

	if cfg.DialTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.DialTimeout)
		defer cancel()
	}

	c, err := newClient(ctx, cfg, "default")
	if err != nil {
		return nil, fmt.Errorf("init milvus client failed: %w", err)
	}

	if err := ensureClientDatabase(ctx, c, cfg.DBName); err != nil {
		_ = c.Close(ctx)
		return nil, fmt.Errorf("init milvus client failed: %w", err)
	}

	client = c
	return client, nil
}

type milvusDatabaseClient interface {
	ListDatabase(context.Context, milvusclient.ListDatabaseOption, ...grpc.CallOption) ([]string, error)
	CreateDatabase(context.Context, milvusclient.CreateDatabaseOption, ...grpc.CallOption) error
	UseDatabase(context.Context, milvusclient.UseDatabaseOption) error
}

func newClient(ctx context.Context, cfg Config, dbName string) (*milvusclient.Client, error) {
	return milvusclient.New(ctx, &milvusclient.ClientConfig{
		Address:       cfg.Address,
		Username:      cfg.Username,
		Password:      cfg.Password,
		DBName:        dbName,
		APIKey:        cfg.APIKey,
		EnableTLSAuth: cfg.EnableTLS,
	})
}

func ensureClientDatabase(ctx context.Context, client milvusDatabaseClient, dbName string) error {
	dbName = strings.TrimSpace(dbName)
	if dbName == "" || dbName == "default" {
		return nil
	}

	exists, err := milvusDatabaseExists(ctx, client, dbName)
	if err != nil {
		return err
	}
	if exists {
		return useMilvusDatabase(ctx, client, dbName)
	}

	if err := client.CreateDatabase(ctx, milvusclient.NewCreateDatabaseOption(dbName)); err != nil {
		if !isMilvusDatabaseAlreadyExistsError(err) {
			return err
		}
	}

	return useMilvusDatabase(ctx, client, dbName)
}

func milvusDatabaseExists(ctx context.Context, client milvusDatabaseClient, dbName string) (bool, error) {
	dbNames, err := client.ListDatabase(ctx, milvusclient.NewListDatabaseOption())
	if err != nil {
		return false, err
	}
	for _, name := range dbNames {
		if strings.TrimSpace(name) == dbName {
			return true, nil
		}
	}
	return false, nil
}

func useMilvusDatabase(ctx context.Context, client milvusDatabaseClient, dbName string) error {
	return client.UseDatabase(ctx, milvusclient.NewUseDatabaseOption(dbName))
}

func isMilvusDatabaseAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "already exists") || strings.Contains(msg, "database already exists")
}

const collectionLoadPollInterval = 200 * time.Millisecond

func EnsureCollectionLoaded(ctx context.Context, client *milvusclient.Client, collectionName string) error {
	collectionName = strings.TrimSpace(collectionName)
	if collectionName == "" {
		return fmt.Errorf("milvus collection name is empty")
	}

	state, err := client.GetLoadState(ctx, milvusclient.NewGetLoadStateOption(collectionName))
	if err != nil {
		return fmt.Errorf("读取 Milvus collection 加载状态失败: %w", err)
	}
	if state.State == entity.LoadStateLoaded {
		return nil
	}
	if state.State == entity.LoadStateLoading {
		return waitCollectionLoaded(ctx, client, collectionName)
	}

	return loadCollectionAndWait(ctx, client, collectionName)
}

func loadCollectionAndWait(ctx context.Context, client *milvusclient.Client, collectionName string) error {
	task, err := client.LoadCollection(ctx, milvusclient.NewLoadCollectionOption(collectionName))
	if err != nil {
		if isMilvusCollectionLoadingOrLoadedError(err) {
			return waitCollectionLoaded(ctx, client, collectionName)
		}
		return fmt.Errorf("加载 Milvus collection 失败: %w", err)
	}
	if err = task.Await(ctx); err != nil {
		return fmt.Errorf("等待 Milvus collection 加载完成失败: %w", err)
	}
	return nil
}

func waitCollectionLoaded(ctx context.Context, client *milvusclient.Client, collectionName string) error {
	ticker := time.NewTicker(collectionLoadPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			state, err := client.GetLoadState(ctx, milvusclient.NewGetLoadStateOption(collectionName))
			if err != nil {
				return fmt.Errorf("读取 Milvus collection 加载状态失败: %w", err)
			}
			if state.State == entity.LoadStateLoaded {
				return nil
			}
			if state.State != entity.LoadStateLoading {
				return loadCollectionAndWait(ctx, client, collectionName)
			}
		}
	}
}

func isMilvusCollectionLoadingOrLoadedError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "not loaded") {
		return false
	}
	return strings.Contains(msg, "already loaded") ||
		strings.Contains(msg, "already loading") ||
		strings.Contains(msg, "loading") ||
		strings.Contains(msg, "load state")
}

func closeClient() error {
	clientMu.Lock()
	defer clientMu.Unlock()

	if client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := client.Close(ctx)
	client = nil
	return err
}

func loadConfig() (Config, bool, error) {
	if !viper.IsSet("milvus") {
		return Config{}, false, nil
	}

	var cfg Config
	if err := viper.UnmarshalKey("milvus", &cfg); err != nil {
		return Config{}, false, fmt.Errorf("unmarshal milvus config failed: %w", err)
	}

	cfg.Address = strings.TrimSpace(cfg.Address)
	if cfg.DialTimeout <= 0 {
		cfg.DialTimeout = 5 * time.Second
	}

	if cfg.Address == "" || !cfg.Enabled {
		return cfg, false, nil
	}

	return cfg, true, nil
}
