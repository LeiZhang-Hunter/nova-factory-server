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
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
)

type Config struct {
	Enabled                    bool          `mapstructure:"enabled"`
	Address                    string        `mapstructure:"address"`
	Username                   string        `mapstructure:"username"`
	Password                   string        `mapstructure:"password"`
	DBName                     string        `mapstructure:"db_name"`
	APIKey                     string        `mapstructure:"api_key"`
	EnableTLS                  bool          `mapstructure:"enable_tls"`
	DialTimeout                time.Duration `mapstructure:"dial_timeout"`
	CollectionLoadTimeout      time.Duration `mapstructure:"collection_load_timeout"`
	CollectionLoadPollInterval time.Duration `mapstructure:"collection_load_poll_interval"`
	CollectionLoadCacheTTL     time.Duration `mapstructure:"collection_load_cache_ttl"`
	AllowedCollections         []string      `mapstructure:"allowed_collections"`
	LoadFields                 []string      `mapstructure:"load_fields"`
	SkipLoadDynamicField       bool          `mapstructure:"skip_load_dynamic_field"`
}

var (
	client                *milvusclient.Client
	clientMu              sync.RWMutex
	collectionLoadGroup   singleflight.Group
	collectionLoadCache   = make(map[string]time.Time)
	collectionLoadCacheMu sync.RWMutex
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

const (
	defaultCollectionLoadTimeout      = 60 * time.Second
	defaultCollectionLoadPollInterval = time.Second
	defaultCollectionLoadCacheTTL     = 30 * time.Second
)

type collectionLoadSettings struct {
	Timeout              time.Duration
	PollInterval         time.Duration
	CacheTTL             time.Duration
	AllowedCollections   map[string]struct{}
	LoadFields           []string
	SkipLoadDynamicField bool
}

func EnsureCollectionLoaded(ctx context.Context, client *milvusclient.Client, collectionName string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	if client == nil {
		return fmt.Errorf("milvus client is nil")
	}

	collectionName = strings.TrimSpace(collectionName)
	if collectionName == "" {
		return fmt.Errorf("milvus collection name is empty")
	}

	settings, err := loadCollectionSettings()
	if err != nil {
		return err
	}
	if !settings.collectionAllowed(collectionName) {
		return fmt.Errorf("milvus collection %q is not allowed", collectionName)
	}

	loadKey := collectionLoadKey(client, collectionName)
	if isCollectionLoadCached(loadKey, settings.CacheTTL) {
		return nil
	}

	ch := collectionLoadGroup.DoChan(loadKey, func() (any, error) {
		if isCollectionLoadCached(loadKey, settings.CacheTTL) {
			return nil, nil
		}

		loadCtx, cancel := context.WithTimeout(context.Background(), settings.Timeout)
		defer cancel()

		return nil, ensureCollectionLoaded(loadCtx, client, collectionName, settings)
	})

	select {
	case <-ctx.Done():
		return ctx.Err()
	case result := <-ch:
		if result.Err == nil {
			markCollectionLoadCached(loadKey)
		}
		return result.Err
	}
}

func collectionLoadKey(client *milvusclient.Client, collectionName string) string {
	return fmt.Sprintf("%p:%s", client, collectionName)
}

func ensureCollectionLoaded(ctx context.Context, client *milvusclient.Client, collectionName string, settings collectionLoadSettings) error {
	state, err := client.GetLoadState(ctx, milvusclient.NewGetLoadStateOption(collectionName))
	if err != nil {
		return fmt.Errorf("读取 Milvus collection 加载状态失败: %w", err)
	}
	if state.State == entity.LoadStateLoaded {
		return nil
	}
	if state.State == entity.LoadStateLoading {
		return waitCollectionLoaded(ctx, client, collectionName, settings)
	}

	return loadCollectionAndWait(ctx, client, collectionName, settings)
}

func loadCollectionAndWait(ctx context.Context, client *milvusclient.Client, collectionName string, settings collectionLoadSettings) error {
	option := milvusclient.NewLoadCollectionOption(collectionName)
	if len(settings.LoadFields) > 0 {
		option.WithLoadFields(settings.LoadFields...)
	}
	if settings.SkipLoadDynamicField {
		option.WithSkipLoadDynamicField(true)
	}

	if _, err := client.LoadCollection(ctx, option); err != nil {
		if isMilvusCollectionLoadingOrLoadedError(err) {
			return waitCollectionLoaded(ctx, client, collectionName, settings)
		}
		return fmt.Errorf("加载 Milvus collection 失败: %w", err)
	}
	return waitCollectionLoaded(ctx, client, collectionName, settings)
}

func waitCollectionLoaded(ctx context.Context, client *milvusclient.Client, collectionName string, settings collectionLoadSettings) error {
	ticker := time.NewTicker(settings.PollInterval)
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
				return loadCollectionAndWait(ctx, client, collectionName, settings)
			}
		}
	}
}

func loadCollectionSettings() (collectionLoadSettings, error) {
	cfg, _, err := loadConfig()
	if err != nil {
		return collectionLoadSettings{}, err
	}
	return newCollectionLoadSettings(cfg), nil
}

func newCollectionLoadSettings(cfg Config) collectionLoadSettings {
	timeout := cfg.CollectionLoadTimeout
	if timeout <= 0 {
		timeout = defaultCollectionLoadTimeout
	}

	pollInterval := cfg.CollectionLoadPollInterval
	if pollInterval <= 0 {
		pollInterval = defaultCollectionLoadPollInterval
	}

	cacheTTL := cfg.CollectionLoadCacheTTL
	if cacheTTL == 0 {
		cacheTTL = defaultCollectionLoadCacheTTL
	}
	if cacheTTL < 0 {
		cacheTTL = 0
	}

	return collectionLoadSettings{
		Timeout:              timeout,
		PollInterval:         pollInterval,
		CacheTTL:             cacheTTL,
		AllowedCollections:   normalizeAllowedCollections(cfg.AllowedCollections),
		LoadFields:           normalizeMilvusLoadFields(cfg.LoadFields),
		SkipLoadDynamicField: cfg.SkipLoadDynamicField,
	}
}

func (s collectionLoadSettings) collectionAllowed(collectionName string) bool {
	if len(s.AllowedCollections) == 0 {
		return true
	}
	_, ok := s.AllowedCollections[collectionName]
	return ok
}

func normalizeAllowedCollections(collections []string) map[string]struct{} {
	if len(collections) == 0 {
		return nil
	}

	allowed := make(map[string]struct{}, len(collections))
	for _, collection := range collections {
		collection = strings.TrimSpace(collection)
		if collection == "" {
			continue
		}
		allowed[collection] = struct{}{}
	}
	if len(allowed) == 0 {
		return nil
	}
	return allowed
}

func isCollectionLoadCached(loadKey string, ttl time.Duration) bool {
	if ttl <= 0 {
		return false
	}

	now := time.Now()

	collectionLoadCacheMu.RLock()
	loadedAt, ok := collectionLoadCache[loadKey]
	collectionLoadCacheMu.RUnlock()
	if !ok {
		return false
	}
	if now.Sub(loadedAt) < ttl {
		return true
	}

	collectionLoadCacheMu.Lock()
	if currentLoadedAt, ok := collectionLoadCache[loadKey]; ok && now.Sub(currentLoadedAt) >= ttl {
		delete(collectionLoadCache, loadKey)
	}
	collectionLoadCacheMu.Unlock()
	return false
}

func markCollectionLoadCached(loadKey string) {
	collectionLoadCacheMu.Lock()
	collectionLoadCache[loadKey] = time.Now()
	collectionLoadCacheMu.Unlock()
}

func resetCollectionLoadCache() {
	collectionLoadCacheMu.Lock()
	collectionLoadCache = make(map[string]time.Time)
	collectionLoadCacheMu.Unlock()
}

func normalizeMilvusLoadFields(fields []string) []string {
	if len(fields) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(fields))
	normalized := make([]string, 0, len(fields))
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		if _, ok := seen[field]; ok {
			continue
		}
		seen[field] = struct{}{}
		normalized = append(normalized, field)
	}
	return normalized
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
		strings.Contains(msg, "collection has been loaded") ||
		strings.Contains(msg, "collection is loading")
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
	resetCollectionLoadCache()
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
