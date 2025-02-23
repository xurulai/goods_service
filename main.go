package main

import (
	"flag"
	"fmt"
	"goods_srv/config"
	"goods_srv/handler"
	"goods_srv/logger"
	"goods_srv/proto"
	"goods_srv/registry"
	"net"
	"os"
	"os/signal"
	"syscall"
	"goods_srv/dao/mysql"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// GoodsServer 是一个 gRPC 服务结构体，实现了 proto.UnimplementedGoodsServer 接口。
type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

func main() {
	// 0.从命令行获取可能的配置文件路径
	// 例如：goods_service -conf="./conf/config_qa.yaml"
	var cfn string
	flag.StringVar(&cfn, "conf", "./conf/config.yaml", "指定配置文件路径")
	flag.Parse()
	//使用 flag 包解析命令行参数，允许用户通过 -conf 参数指定配置文件路径。

	// 1. 加载配置文件
	err := config.Init(cfn)
	if err != nil {
		panic(err) // 程序启动时加载配置文件失败直接退出
	}

	// 2. 加载日志
	err = logger.Init(config.Conf.LogConfig, config.Conf.Mode)
	if err != nil {
		panic(err) // 程序启动时初始化日志模块失败直接退出
	}

	// 3. 初始化 MySQL 数据库连接
	err = mysql.Init(config.Conf.MySQLConfig)
	if err != nil {
		panic(err) // 初始化 MySQL 失败，程序直接退出
	}

	// 4. 初始化 Consul 服务注册中心
	err = registry.Init(config.Conf.ConsulConfig.Addr)
	if err != nil {
		panic(err) // 初始化注册中心失败，程序直接退出
	}

	// 监听端口
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Conf.IP, config.Conf.Port))
	if err != nil {
		panic(err)
	}

	// 创建 gRPC 服务
	s := grpc.NewServer()
	// 注册商品服务到 gRPC 服务
	proto.RegisterGoodsServer(s, &handler.GoodsSrv{})

	// 启动 gRPC 服务
	go func() {
		err = s.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

	// 注册服务到 Consul
	registry.Reg.RegisterService(config.Conf.Name, config.Conf.IP, config.Conf.Port, nil)

	// 打印服务启动日志
	zap.L().Info("service start...")

	// 服务退出时要注销服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit // 正常会 hang 在此处

	// 退出时注销服务
	serviceId := fmt.Sprintf("%s-%s-%d", config.Conf.Name, config.Conf.IP, config.Conf.Port)
	registry.Reg.Deregister(serviceId)
}
