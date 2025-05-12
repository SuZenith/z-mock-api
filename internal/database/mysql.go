package database

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"kite/internal/configs"
	KiteLogger "kite/pkg/logger"
	"time"
)

var DB *gorm.DB

func InitDB(cfg *configs.MySQLConfig) error {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.Charset,
	)

	gormConfig := &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		// 简直自动创建外建约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 日志配置
		Logger: logger.New(
			&GormLogWriter{},
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
	}

	// 连接数据库
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}
	// 获取 SQL DB 对象，用于设置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifeConn) * time.Second)

	// 测试数据库连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	KiteLogger.Info("Successfully connected to database",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.Name),
	)

	return nil
}

func GetDB() *gorm.DB {
	return DB
}

func CloseDB() error {
	if DB == nil {
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}

func WithTransaction(txFn func(tx *gorm.DB) error) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return txFn(tx)
	})
}

type GormLogWriter struct{}

func (w *GormLogWriter) Printf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	KiteLogger.Debug("gorm", zap.String("sql", msg))
}
