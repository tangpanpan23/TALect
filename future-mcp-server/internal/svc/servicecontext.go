package svc

import (
	"github.com/future-mcp/future-mcp-server/internal/config"
	"github.com/future-mcp/future-mcp-server/internal/service"
)

type ServiceContext struct {
	Config          config.Config
	MCPService      *service.MCPService
	MaterialService service.MaterialService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
