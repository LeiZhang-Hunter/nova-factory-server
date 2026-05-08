package mysql

import (
	"fmt"
	"log"
	"os"
	"time"

	applogger "nova-factory-server/app/utils/logger"

	"github.com/baizeplus/sqly"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

//type Sql interface {
//	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
//	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
//	NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqly.Rows, error)
//	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
//	NamedSelectContext(ctx context.Context, dest interface{}, query string, arg interface{}) error
//	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
//	NamedGetContext(ctx context.Context, dest interface{}, query string, args interface{}) error
//	NamedSelectPageContext(ctx context.Context, dest interface{}, total *int64, query string, page sqly.Page) error
//	MustBegin() *sqly.Tx
//}

func NewData() (sqly.SqlyContext, func(), error) {
	type Mysql struct {
		Host         string `mapstructure:"host"`
		User         string `mapstructure:"user"`
		Password     string `mapstructure:"password"`
		DB           string `mapstructure:"dbname"`
		Port         int    `mapstructure:"port"`
		MaxOpenConns int    `mapstructure:"max_open_conns"`
		MaxIdleConns int    `mapstructure:"max_idle_conns"`
	}
	// 把读取到的配置信息反序列化到 Conf 变量中
	var d Mysql
	if err := viper.UnmarshalKey("mysql", &d); err != nil {
		panic(err)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", d.User, d.Password, d.Host, d.Port, d.DB) + "?parseTime=true&loc=Asia%2FShanghai"
	var db *sqly.DB
	var err error
	for {
		db, err = sqly.Connect("mysql", dsn)
		if err != nil {
			zap.L().Error("connect mysql failed", zap.Error(err))
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	db.SetMaxOpenConns(d.MaxOpenConns)
	db.SetMaxIdleConns(d.MaxIdleConns)
	db.SetConnMaxLifetime(time.Minute * 5)
	sqly.SetLog(new(applogger.SqlyLog))

	return db, func() {

	}, err
}

// NewDB .
func NewDB() *gorm.DB {
	type Mysql struct {
		Host          string `mapstructure:"host"`
		User          string `mapstructure:"user"`
		Password      string `mapstructure:"password"`
		DB            string `mapstructure:"dbname"`
		Port          int    `mapstructure:"port"`
		MaxOpenConns  int    `mapstructure:"max_open_conns"`
		MaxIdleConns  int    `mapstructure:"max_idle_conns"`
		LogLevel      string `mapstructure:"log_level"`
		SlowThreshold int    `mapstructure:"slow_threshold"`
	}
	// 把读取到的配置信息反序列化到 Conf 变量中
	var d Mysql
	if err := viper.UnmarshalKey("mysql", &d); err != nil {
		panic(err)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", d.User, d.Password, d.Host, d.Port, d.DB) + "?parseTime=true&loc=Asia%2FShanghai"

	// 解析 GORM 日志级别
	var gormLogLevel gormlogger.LogLevel
	switch d.LogLevel {
	case "silent":
		gormLogLevel = gormlogger.Silent
	case "error":
		gormLogLevel = gormlogger.Error
	case "warn":
		gormLogLevel = gormlogger.Warn
	case "info":
		gormLogLevel = gormlogger.Info
	default:
		gormLogLevel = gormlogger.Silent
	}
	gdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			gormlogger.Config{
				SlowThreshold:             time.Duration(d.SlowThreshold) * time.Millisecond,
				LogLevel:                  gormLogLevel,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
				ParameterizedQueries:      false,
			},
		),
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(d.MaxOpenConns)
	sqlDB.SetMaxIdleConns(d.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	return gdb
}
