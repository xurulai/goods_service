package registry

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

// 定义一个结构体，用于封装 Consul 客户端
type consul struct {
	client *api.Client // Consul 客户端实例
}

// 定义一个全局变量，用于存储注册中心的接口实例
var Reg Register  // Register 是一个接口类型，稍后会定义

// 确保 consul 结构体实现了 Register 接口
var _ Register = (*consul)(nil)

// Init 函数用于初始化 Consul 客户端，并将其绑定到全局变量 Reg
func Init(addr string) (err error) {
	// 创建 Consul 客户端配置
	cfg := api.DefaultConfig()
	// 设置 Consul 服务的地址
	cfg.Address = addr
	// 根据配置创建 Consul 客户端
	c, err := api.NewClient(cfg)
	if err != nil {
		return err // 如果创建失败，返回错误
	}
	// 将初始化好的 Consul 客户端封装到 consul 结构体中，并赋值给全局变量 Reg
	Reg = &consul{c}
	return
}

// RegisterService 方法用于将 gRPC 服务注册到 Consul
func (c *consul) RegisterService(serviceName string, ip string, port int, tags []string) error {
	// 定义健康检查配置
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", ip, port), // gRPC 服务地址，必须是外部可访问的地址
		Timeout:                        "10s",                         // 健康检查超时时间
		Interval:                       "10s",                         // 健康检查间隔
		DeregisterCriticalServiceAfter: "20s",                         // 服务连续失败后自动注销的时间
	}
	// 定义服务注册信息
	srv := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", serviceName, ip, port), // 服务唯一标识，由服务名、IP 和端口组成
		Name:    serviceName,                                    // 服务名称
		Tags:    tags,                                           // 服务标签，用于分类或标识服务
		Address: ip,                                             // 服务地址
		Port:    port,                                           // 服务端口
		Check:   check,                                          // 健康检查配置
	}
	// 调用 Consul 客户端的 Agent().ServiceRegister 方法，将服务注册到 Consul
	return c.client.Agent().ServiceRegister(srv)
}

// ListService 方法用于从 Consul 中获取指定服务的所有实例
func (c *consul) ListService(serviceName string) (map[string]*api.AgentService, error) {
	// 使用 Consul 的过滤功能，根据服务名称获取服务实例
	return c.client.Agent().ServicesWithFilter(fmt.Sprintf("Service==`%s`", serviceName))
}

// Deregister 方法用于从 Consul 中注销指定的服务实例
func (c *consul) Deregister(serviceID string) error {
	// 调用 Consul 客户端的 Agent().ServiceDeregister 方法，注销服务
	return c.client.Agent().ServiceDeregister(serviceID)
}