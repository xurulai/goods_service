package registry

import "github.com/hashicorp/consul/api"

// Register 是一个接口，定义了注册中心需要实现的方法
type Register interface {
	// RegisterService 方法用于注册服务
	RegisterService(serviceName string, ip string, port int, tags []string) error
	// ListService 方法用于发现服务
	ListService(serviceName string) (map[string]*api.AgentService, error)
	// Deregister 方法用于注销服务
	Deregister(serviceID string) error
}
