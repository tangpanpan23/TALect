package database

import (
	"context"
	"fmt"
	"time"

	"github.com/future-mcp/future-mcp-server/pkg/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DB 全局数据库实例
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() (*gorm.DB, error) {
	dsn := buildDSN()

	// 配置GORM日志
	gormLogger := gormlogger.New(
		&GormWriter{},
		gormlogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  getGormLogLevel(),
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// 连接数据库
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		PrepareStmt:                              true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(viper.GetInt("database.max_idle_conns"))
	sqlDB.SetMaxOpenConns(viper.GetInt("database.max_open_conns"))
	sqlDB.SetConnMaxLifetime(time.Duration(viper.GetInt("database.conn_max_lifetime")) * time.Second)

	logger.Info("Database connected successfully",
		logger.Field("host", viper.GetString("database.host")),
		logger.Field("port", viper.GetInt("database.port")),
		logger.Field("database", viper.GetString("database.dbname")),
	)

	return DB, nil
}

// buildDSN 构建数据库连接字符串
func buildDSN() string {
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	user := viper.GetString("database.user")
	password := viper.GetString("database.password")
	dbname := viper.GetString("database.dbname")
	sslmode := viper.GetString("database.sslmode")

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	// 添加其他参数
	timezone := viper.GetString("database.timezone")
	if timezone != "" {
		dsn += fmt.Sprintf(" TimeZone=%s", timezone)
	}

	return dsn
}

// getGormLogLevel 获取GORM日志级别
func getGormLogLevel() gormlogger.LogLevel {
	switch viper.GetString("log.level") {
	case "debug":
		return gormlogger.Info
	case "info":
		return gormlogger.Info
	case "warn":
		return gormlogger.Warn
	case "error":
		return gormlogger.Error
	default:
		return gormlogger.Info
	}
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// Close 关闭数据库连接
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// HealthCheck 数据库健康检查
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// Migrate 执行数据库迁移
func Migrate(models ...interface{}) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	for _, model := range models {
		if err := DB.AutoMigrate(model); err != nil {
			logger.Error("Failed to migrate model", logger.Field("model", fmt.Sprintf("%T", model)), logger.Error(err))
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
		logger.Info("Model migrated successfully", logger.Field("model", fmt.Sprintf("%T", model)))
	}

	return nil
}

// Transaction 执行数据库事务
func Transaction(fn func(tx *gorm.DB) error) error {
	return DB.Transaction(fn)
}

// WithContext 使用上下文执行查询
func WithContext(ctx context.Context) *gorm.DB {
	return DB.WithContext(ctx)
}

// GormWriter GORM日志写入器
type GormWriter struct{}

// Write 实现io.Writer接口
func (w *GormWriter) Write(p []byte) (int, error) {
	logger.Info(string(p))
	return len(p), nil
}
