package logger

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var (
	// 全局日志实例
	log     *zap.Logger
	once    sync.Once
	logLock sync.RWMutex
)

type Config struct {
	Level    string `mapstructure:"level"`
	Encoding string `mapstructure:"encoding"`
	Dev      bool   `mapstructure:"dev"`
}

func Init(cfg *Config) {
	once.Do(func() {
		// 默认配置
		if cfg == nil {
			cfg = &Config{
				Level:    "info",
				Encoding: "json",
				Dev:      false,
			}
		}
		// 解析日志级别
		level := zap.InfoLevel
		switch cfg.Level {
		case "debug":
			level = zap.DebugLevel
		case "info":
			level = zap.InfoLevel
		case "warn":
			level = zap.WarnLevel
		case "error":
			level = zap.ErrorLevel
		case "dpanic":
			level = zap.DPanicLevel
		case "panic":
			level = zap.PanicLevel
		case "fatal":
			level = zap.FatalLevel
		}
		// 为生产环境优化的编码器配置
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.MillisDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		// 为开发环境优化编码器配置
		if cfg.Dev {
			encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
			encoderConfig.TimeKey = "time"
			encoderConfig.MessageKey = "msg"
		}
		// 创建 core
		var enc zapcore.Encoder
		if cfg.Encoding == "console" {
			enc = zapcore.NewConsoleEncoder(encoderConfig)
		} else {
			enc = zapcore.NewJSONEncoder(encoderConfig)
		}

		// 输出到 stdout
		core := zapcore.NewCore(
			enc,
			zapcore.AddSync(os.Stdout),
			zap.NewAtomicLevelAt(level),
		)
		// 创建日志实例
		logger := zap.New(
			core,
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
		// 如果是开发模式，启用开发模式日志
		if cfg.Dev {
			logger = logger.WithOptions(zap.Development())
		}
		// 替换全局日志实例
		logLock.Lock()
		log = logger
		logLock.Unlock()

		// 替换 zap 的全局日志实例
		zap.ReplaceGlobals(logger)
	})
}

func GetLogger() *zap.Logger {
	logLock.RLock()
	defer logLock.RUnlock()

	if log == nil {
		// 如果日志未初始化，则使用默认初始化
		logLock.RUnlock()
		Init(nil)
		logLock.RLock()
	}

	return log
}

// WithRequestId 创建带有请求 ID 的日志实例
func WithRequestId(requestId string) *zap.Logger {
	if requestId == "" {
		return GetLogger()
	}
	return GetLogger().With(zap.String("request_id", requestId))
}

// FromContext 从 Echo 上下文中获取日志实例
func FromContext(c echo.Context) *zap.Logger {
	if c == nil {
		return GetLogger()
	}
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	return WithRequestId(requestId)
}

// WithFields 创建带有字段的日志实例
func WithFields(fields ...zap.Field) *zap.Logger {
	return GetLogger().With(fields...)
}

// Debug 记录 debug 级别的日志
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Info 记录 info 级别的日志
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Warn 记录 warn 级别的日志
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Error 记录 error 级别的日志
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Fatal 记录 fatal 级别的日志
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}

// Panic 记录 panic 级别的日志
func Panic(msg string, fields ...zap.Field) {
	GetLogger().Panic(msg, fields...)
}

// Sync 同步日志缓冲区
func Sync() error {
	return GetLogger().Sync()
}

// DebugC 记录带有上下文的调试级别日志
func DebugC(c echo.Context, msg string, fields ...zap.Field) {
	FromContext(c).Debug(msg, fields...)
}

// InfoC 记录带有上下文的信息级别日志
func InfoC(c echo.Context, msg string, fields ...zap.Field) {
	FromContext(c).Info(msg, fields...)
}

// WarnC 记录带有上下文的警告级别日志
func WarnC(c echo.Context, msg string, fields ...zap.Field) {
	FromContext(c).Warn(msg, fields...)
}

// ErrorC 记录带有上下文的错误级别日志
func ErrorC(c echo.Context, msg string, fields ...zap.Field) {
	FromContext(c).Error(msg, fields...)
}

// FatalC 记录带有上下文的致命级别日志
func FatalC(c echo.Context, msg string, fields ...zap.Field) {
	FromContext(c).Fatal(msg, fields...)
}
