package main

import (
	"flag"
	"fmt"

	"github.com/future-mcp/future-mcp-server/internal/config"
	"github.com/future-mcp/future-mcp-server/internal/handler"
	"github.com/future-mcp/future-mcp-server/internal/repository"
	"github.com/future-mcp/future-mcp-server/internal/service"
	"github.com/future-mcp/future-mcp-server/internal/svc"
	"github.com/future-mcp/future-mcp-server/pkg/logger"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*config.ConfigFile, &c)

	// 初始化日志
	if err := logger.Init(); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	logger.Info("Starting TALink MCP Server...")

	// 初始化服务上下文
	ctx := svc.NewServiceContext(c)

	// 初始化缓存服务 (暂时使用内存实现)
	cacheService := service.NewMemoryCacheService()

	// 初始化存储库 (暂时使用内存实现)
	materialRepo := repository.NewMemoryMaterialRepository()

	// 初始化素材服务
	materialService := service.NewMaterialService(materialRepo, cacheService)

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

	// 更新服务上下文中的服务
	ctx.MCPService = mcpService
	ctx.MaterialService = materialService

	// 创建服务器
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 注册路由
	handler.RegisterHandlers(server, ctx)

	logger.Info("Server starting", logger.Any("port", c.Port))

	// 启动服务器
	server.Start()
}

