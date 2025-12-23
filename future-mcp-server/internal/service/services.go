package service

// Services 服务集合
type Services struct {
	Material MaterialService
	Tool     ToolService
	Resource ResourceService
	User     UserService
}

// ServiceDeps 服务依赖
type ServiceDeps struct {
	Repos       interface{}
	Cache       interface{}
	AuthService interface{}
}

// NewServices 创建服务集合
func NewServices(deps ServiceDeps) *Services {
	return &Services{}
}
