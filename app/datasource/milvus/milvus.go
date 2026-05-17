package milvus

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"github.com/spf13/viper"
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

	c, err := milvusclient.New(ctx, &milvusclient.ClientConfig{
		Address:       cfg.Address,
		Username:      cfg.Username,
		Password:      cfg.Password,
		DBName:        cfg.DBName,
		APIKey:        cfg.APIKey,
		EnableTLSAuth: cfg.EnableTLS,
	})
	if err != nil {
		return nil, fmt.Errorf("init milvus client failed: %w", err)
	}

	client = c
	return client, nil
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
