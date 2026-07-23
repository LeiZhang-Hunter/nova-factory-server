package milvus

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type fakeMilvusDatabaseClient struct {
	listDatabases []string
	listErr       error
	createErr     error
	useErr        error

	created []string
	used    []string
}

func (f *fakeMilvusDatabaseClient) ListDatabase(context.Context, milvusclient.ListDatabaseOption, ...grpc.CallOption) ([]string, error) {
	return f.listDatabases, f.listErr
}

func (f *fakeMilvusDatabaseClient) CreateDatabase(_ context.Context, opt milvusclient.CreateDatabaseOption, _ ...grpc.CallOption) error {
	f.created = append(f.created, opt.Request().GetDbName())
	return f.createErr
}

func (f *fakeMilvusDatabaseClient) UseDatabase(_ context.Context, opt milvusclient.UseDatabaseOption) error {
	f.used = append(f.used, opt.DbName())
	return f.useErr
}

func TestEnsureClientDatabaseSkipsDefault(t *testing.T) {
	fake := &fakeMilvusDatabaseClient{}

	if err := ensureClientDatabase(context.Background(), fake, "default"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fake.created) != 0 || len(fake.used) != 0 {
		t.Fatalf("unexpected calls: created=%v used=%v", fake.created, fake.used)
	}
}

func TestEnsureClientDatabaseCreatesMissingDatabase(t *testing.T) {
	fake := &fakeMilvusDatabaseClient{listDatabases: []string{"default"}}

	if err := ensureClientDatabase(context.Background(), fake, "nova"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fake.created) != 1 || fake.created[0] != "nova" {
		t.Fatalf("unexpected create calls: %v", fake.created)
	}
	if len(fake.used) != 1 || fake.used[0] != "nova" {
		t.Fatalf("unexpected use calls: %v", fake.used)
	}
}

func TestEnsureClientDatabaseUsesExistingDatabase(t *testing.T) {
	fake := &fakeMilvusDatabaseClient{listDatabases: []string{"default", "nova"}}

	if err := ensureClientDatabase(context.Background(), fake, "nova"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fake.created) != 0 {
		t.Fatalf("unexpected create calls: %v", fake.created)
	}
	if len(fake.used) != 1 || fake.used[0] != "nova" {
		t.Fatalf("unexpected use calls: %v", fake.used)
	}
}

func TestIsMilvusDatabaseAlreadyExistsError(t *testing.T) {
	if !isMilvusDatabaseAlreadyExistsError(errors.New("database already exists")) {
		t.Fatal("expected already-exists error to match")
	}
	if isMilvusDatabaseAlreadyExistsError(errors.New("database not found")) {
		t.Fatal("unexpected match for not-found error")
	}
}

func TestIsMilvusCollectionLoadingOrLoadedError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{name: "already loaded", err: errors.New("collection already loaded"), want: true},
		{name: "already loading", err: errors.New("collection already loading"), want: true},
		{name: "has been loaded", err: errors.New("collection has been loaded"), want: true},
		{name: "is loading", err: errors.New("collection is loading"), want: true},
		{name: "not loaded", err: errors.New("collection not loaded"), want: false},
		{name: "generic loading failure", err: errors.New("loading collection failed: out of memory"), want: false},
		{name: "load state failure", err: errors.New("failed to read load state"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isMilvusCollectionLoadingOrLoadedError(tt.err); got != tt.want {
				t.Fatalf("isMilvusCollectionLoadingOrLoadedError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCollectionLoadSettingsDefaults(t *testing.T) {
	settings := newCollectionLoadSettings(Config{})

	if settings.Timeout != defaultCollectionLoadTimeout {
		t.Fatalf("unexpected timeout: %v", settings.Timeout)
	}
	if settings.PollInterval != defaultCollectionLoadPollInterval {
		t.Fatalf("unexpected poll interval: %v", settings.PollInterval)
	}
	if settings.CacheTTL != defaultCollectionLoadCacheTTL {
		t.Fatalf("unexpected cache ttl: %v", settings.CacheTTL)
	}
	if len(settings.AllowedCollections) != 0 {
		t.Fatalf("unexpected allowed collections: %v", settings.AllowedCollections)
	}
	if len(settings.LoadFields) != 0 {
		t.Fatalf("unexpected load fields: %v", settings.LoadFields)
	}
	if settings.SkipLoadDynamicField {
		t.Fatal("unexpected skip dynamic field")
	}
}

func TestNewCollectionLoadSettingsNormalizesFields(t *testing.T) {
	settings := newCollectionLoadSettings(Config{
		CollectionLoadTimeout:      2 * time.Minute,
		CollectionLoadPollInterval: 1500 * time.Millisecond,
		CollectionLoadCacheTTL:     45 * time.Second,
		AllowedCollections:         []string{" product_vectors ", "", "archive_vectors"},
		LoadFields:                 []string{" product_id ", "", "vector", "product_id", " title "},
		SkipLoadDynamicField:       true,
	})

	if settings.Timeout != 2*time.Minute {
		t.Fatalf("unexpected timeout: %v", settings.Timeout)
	}
	if settings.PollInterval != 1500*time.Millisecond {
		t.Fatalf("unexpected poll interval: %v", settings.PollInterval)
	}
	if settings.CacheTTL != 45*time.Second {
		t.Fatalf("unexpected cache ttl: %v", settings.CacheTTL)
	}
	if !settings.collectionAllowed("product_vectors") || !settings.collectionAllowed("archive_vectors") {
		t.Fatalf("expected configured collections to be allowed: %v", settings.AllowedCollections)
	}
	if settings.collectionAllowed("other_vectors") {
		t.Fatalf("unexpected allowed collection: %v", settings.AllowedCollections)
	}
	wantFields := []string{"product_id", "vector", "title"}
	if !reflect.DeepEqual(settings.LoadFields, wantFields) {
		t.Fatalf("unexpected load fields: got=%v want=%v", settings.LoadFields, wantFields)
	}
	if !settings.SkipLoadDynamicField {
		t.Fatal("expected skip dynamic field")
	}
}

func TestLoadConfigParsesCollectionLoadSettings(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(strings.NewReader(`
milvus:
  enabled: true
  address: 127.0.0.1:19530
  db_name: nova
  dial_timeout: 3s
  collection_load_timeout: 90s
  collection_load_poll_interval: 2s
  collection_load_cache_ttl: 15s
  allowed_collections:
    - product_vectors
  load_fields:
    - product_id
    - vector
  skip_load_dynamic_field: true
`)); err != nil {
		t.Fatalf("read config failed: %v", err)
	}

	cfg, enabled, err := loadConfig()
	if err != nil {
		t.Fatalf("load config failed: %v", err)
	}
	if !enabled {
		t.Fatal("expected milvus config enabled")
	}
	if cfg.DialTimeout != 3*time.Second {
		t.Fatalf("unexpected dial timeout: %v", cfg.DialTimeout)
	}
	if cfg.CollectionLoadTimeout != 90*time.Second {
		t.Fatalf("unexpected collection load timeout: %v", cfg.CollectionLoadTimeout)
	}
	if cfg.CollectionLoadPollInterval != 2*time.Second {
		t.Fatalf("unexpected collection load poll interval: %v", cfg.CollectionLoadPollInterval)
	}
	if cfg.CollectionLoadCacheTTL != 15*time.Second {
		t.Fatalf("unexpected collection load cache ttl: %v", cfg.CollectionLoadCacheTTL)
	}
	if !reflect.DeepEqual(cfg.AllowedCollections, []string{"product_vectors"}) {
		t.Fatalf("unexpected allowed collections: %v", cfg.AllowedCollections)
	}
	if !reflect.DeepEqual(cfg.LoadFields, []string{"product_id", "vector"}) {
		t.Fatalf("unexpected load fields: %v", cfg.LoadFields)
	}
	if !cfg.SkipLoadDynamicField {
		t.Fatal("expected skip dynamic field")
	}
}

func TestCollectionLoadCache(t *testing.T) {
	resetCollectionLoadCache()
	t.Cleanup(resetCollectionLoadCache)

	loadKey := "client:product_vectors"
	if isCollectionLoadCached(loadKey, time.Minute) {
		t.Fatal("unexpected cached collection")
	}

	markCollectionLoadCached(loadKey)
	if !isCollectionLoadCached(loadKey, time.Minute) {
		t.Fatal("expected cached collection")
	}
	if isCollectionLoadCached(loadKey, -time.Second) {
		t.Fatal("expected disabled cache for non-positive ttl")
	}
	if isCollectionLoadCached(loadKey, time.Nanosecond) {
		t.Fatal("expected expired cache")
	}
}

func TestNewCollectionLoadSettingsAllowsDisablingCache(t *testing.T) {
	settings := newCollectionLoadSettings(Config{CollectionLoadCacheTTL: -time.Second})
	if settings.CacheTTL != 0 {
		t.Fatalf("unexpected cache ttl: %v", settings.CacheTTL)
	}
}
