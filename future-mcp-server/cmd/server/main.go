package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/future-mcp/future-mcp-server/internal/handler"
	"github.com/future-mcp/future-mcp-server/internal/middleware"
	"github.com/future-mcp/future-mcp-server/internal/repository"
	"github.com/future-mcp/future-mcp-server/internal/service"
	"github.com/future-mcp/future-mcp-server/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// 初始化配置
	if err := initConfig(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// 初始化日志
	if err := logger.Init(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	logger.Info("Starting TALink MCP Server...")

	// 初始化缓存服务 (暂时使用内存实现)
	cacheService := service.NewMemoryCacheService()

	// 初始化存储库 (暂时使用内存实现)
	materialRepo := repository.NewMemoryMaterialRepository()
	// TODO: 实现其他仓库
	repos := &repository.Repositories{
		Material: materialRepo,
	}

	// 认证服务暂时未实现
	// authService := auth.NewService(viper.GetString("auth.jwt_secret"))

	// 初始化素材服务
	materialService := service.NewMaterialService(repos.Material, cacheService)

	// 初始化MCP服务
	mcpService := service.NewMCPService(&service.MCPServiceConfig{
		MaterialService: materialService,
	})

	// 初始化工具服务
	toolService := service.NewToolService(mcpService)

	// 设置工具服务的MCP引用
	if ts, ok := toolService.(*service.ToolServiceImpl); ok {
		ts.SetMCPService(mcpService)
	}

	// 初始化Gin路由
	r := setupRouter(mcpService)

	// 获取服务器配置
	host := viper.GetString("server.host")
	port := viper.GetInt("server.port")
	addr := fmt.Sprintf("%s:%d", host, port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// 启动服务器
	go func() {
		logger.Info("Server starting", logger.Any("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", logger.Any("error", err))
		}
	}()

	// 等待中断信号优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// 上下文用于通知服务器它有5秒的时间完成当前正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", logger.Any("error", err))
	}

	logger.Info("Server exited")
}

func initConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// 设置默认值
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Warn("Config file not found, using defaults")
		} else {
			return fmt.Errorf("failed to read config: %w", err)
		}
	}

	return nil
}

func setDefaults() {
	// 服务器配置
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "release")

	// 数据库配置
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "future_mcp")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.dbname", "future_mcp")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)

	// Redis配置
	viper.SetDefault("redis.host", "localhost:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)

	// 认证配置
	viper.SetDefault("auth.jwt_secret", "your-jwt-secret-key")
	viper.SetDefault("auth.jwt_expire", 86400)

	// 向量搜索配置
	viper.SetDefault("vector_search.provider", "pinecone")
	viper.SetDefault("vector_search.api_key", "")
	viper.SetDefault("vector_search.environment", "us-east-1")
	viper.SetDefault("vector_search.index_name", "future-materials")
	viper.SetDefault("vector_search.dimension", 768)

	// 存储配置
	viper.SetDefault("storage.provider", "local")
	viper.SetDefault("storage.bucket", "future-materials")
	viper.SetDefault("storage.region", "us-east-1")
	viper.SetDefault("storage.endpoint", "")
	viper.SetDefault("storage.access_key", "")
	viper.SetDefault("storage.secret_key", "")

	// 日志配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "stdout")
}

func setupRouter(mcpService *service.MCPService) *gin.Engine {
	if viper.GetString("server.mode") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 全局中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())

	// 健康检查
	r.GET("/health", handler.HealthCheck)
	r.GET("/ready", handler.ReadinessCheck)

	// API v1 路由组 (暂时简化)
	// v1 := r.Group("/api/v1")
	// v1.Use(middleware.AuthRequired(authService)) // 暂时移除认证

	// 素材相关路由 (暂时简化)
	// materials := v1.Group("/materials")
	// {
	// 	materials.POST("/search", handler.SearchMaterials(services.Material))
	// }

	// MCP协议路由
	mcpGroup := r.Group("/mcp")
	{
		mcpGroup.POST("/jsonrpc", handler.MCPHandler(mcpService))
		mcpGroup.GET("/sse", handler.MCPSSEHandler(mcpService))
	}

	// WebSocket路由（用于实时通信）
	// r.GET("/ws", handler.WebSocketHandler(services)) // 暂时移除

	return r
}
