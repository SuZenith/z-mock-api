package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"kite/internal/api/handlers"
	"kite/internal/api/handlers/fund_pay"
	"kite/internal/api/routes"
	"kite/internal/api/validators"
	"kite/internal/configs"
	"kite/internal/database"
	KiteLogger "kite/pkg/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Echo            *echo.Echo
	MySQLConnection *database.MySQLConnection
	WithdrawHandler *fund_pay.WithdrawHandler
}

func NewServer(echo *echo.Echo, withdrawHandler *fund_pay.WithdrawHandler, mysqlConnection *database.MySQLConnection) *Server {
	return &Server{echo, mysqlConnection, withdrawHandler}
}

func main() {
	// 加载配置
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 初始化日志
	KiteLogger.Init(&KiteLogger.Config{
		Level:    cfg.Log.Level,
		Encoding: cfg.Log.Encoding,
		Dev:      cfg.Log.Dev,
	})
	defer func() {
		err := KiteLogger.Sync()
		if err != nil {
			log.Fatalf("Failed to sync logger: %v\n", err)
		}
	}()

	server, err := InitializeApp(&cfg.Database, echo.New())
	if err != nil {
		return
	}

	// 注册全局中间件
	registerGlobalMiddlewares(server.Echo)
	// 注册自定义错误处理器
	server.Echo.HTTPErrorHandler = handlers.CustomHTTPErrorHandler
	// 注册自定义验证器
	server.Echo.Validator = validators.NewCustomValidator()

	// 注册路由
	routes.RegisterRoutes(server.Echo, server.WithdrawHandler)

	// 创建通道接受关机信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	// 启动 HTTP 服务
	go func() {
		KiteLogger.Info("Starting server", zap.String("addr", serverAddr))
		if err := server.Echo.Start(serverAddr); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				KiteLogger.Info("Shutting down the server")
			} else {
				KiteLogger.Error("Failed to start server", zap.Error(err))
			}
		}
	}()

	// 等待关机信号
	sig := <-quit
	KiteLogger.Info("Received signal", zap.String("signal", sig.String()))

	// 设置关机超时等待时间
	serverTimeout := time.Duration(cfg.Server.ShutdownTimeout) * time.Second
	if serverTimeout == 0 {
		serverTimeout = 10 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), serverTimeout)
	defer cancel()

	// 优雅关机
	if err := server.Echo.Shutdown(ctx); err != nil {
		KiteLogger.Error("Server shutdown failed:", zap.Error(err))
	}

	// 关闭数据库连接
	err = server.MySQLConnection.CloseDB()
	if err != nil {
		KiteLogger.Error("Failed to close database", zap.Error(err))
		os.Exit(1)
	}
	KiteLogger.Info("Server gracefully stopped")
}

func registerGlobalMiddlewares(e *echo.Echo) {
	// 生成请求ID
	e.Use(middleware.RequestID())
	// 记录请求日志
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		Skipper: func(c echo.Context) bool {
			return c.Request().URL.Path == "/health"
		},
		LogStatus: true,
		LogURI:    true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			KiteLogger.InfoC(c, "request", zap.String("URI", v.URI), zap.Int("status", v.Status))
			return nil
		},
	}))
	// 捕获 panic 转换为 500 错误
	e.Use(middleware.Recover())
}
