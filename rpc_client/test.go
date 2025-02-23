package main

import (
	"context"
	"goods_srv/proto"
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	conn   *grpc.ClientConn  // gRPC 客户端连接
	client proto.GoodsClient // gRPC 客户端实例
)

func init() {
	var err error
	conn, err = grpc.Dial(
		"127.0.0.1:8384", // 确保这是 gRPC 服务的端口
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*10)), // 设置最大接收消息大小为 10MB
	)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	client = proto.NewGoodsClient(conn) // 创建 gRPC 客户端实例
}

func TestGetGoodsDetail(wg *sync.WaitGroup, index int) {
	defer wg.Done()

	param := &proto.GetGoodsDetailReq{
		GoodsId: int64(index + 1001), // 假设商品 ID 从 1001 开始
		UserId:  1,                   // 假设用户 ID 为 1
	}
	log.Printf("Sending request with GoodsId: %d", param.GoodsId)

	resp, err := client.GetGoodsDetail(context.Background(), param)
	if err != nil {
		log.Printf("Error calling GetGoodsDetail: %v", err)
	} else {
		log.Printf("Response: %+v", resp)
	}
}

func main() {
	defer conn.Close()    // 程序结束时关闭 gRPC 客户端连接
	var wg sync.WaitGroup // 使用 WaitGroup 等待所有协程完成

	// 启动多个协程测试 GetGoodsDetail 方法
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go TestGetGoodsDetail(&wg, i)
	}
	wg.Wait() // 等待所有协程完成
}

var num int = 100 // 全局变量，初始值为 100

//这是一个简单的并发测试函数，用于模拟多个协程对全局变量 num 的操作。
//num = num - 1：每次调用减 1。由于没有加锁保护，这会导致并发问题（如竞态条件）
