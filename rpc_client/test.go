package main

import (
	"context"
	"goods_srv/proto" // 引入定义了 gRPC 服务的 proto 文件
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"                      // 引入 gRPC 客户端库
	"google.golang.org/grpc/credentials/insecure" // 引入用于不安全连接的证书
)

var (
	conn   *grpc.ClientConn  // gRPC 客户端连接
	client proto.GoodsClient // gRPC 客户端实例
)

// 初始化 gRPC 客户端连接和客户端实例
func init() {
	var err error
	conn, err = grpc.Dial(
		"127.0.0.1:8391", // gRPC 服务地址
		grpc.WithTransportCredentials(insecure.NewCredentials()), // 使用不安全的连接（仅用于测试环境）
	)
	if err != nil {
		panic(err) // 如果连接失败，直接 panic
	}
	client = proto.NewGoodsClient(conn) // 创建 gRPC 客户端对象
}

// 测试 GetGoodsDetail 方法的函数
func TestGetGoodsDetail(wg *sync.WaitGroup, index int) {
	defer wg.Done() // 在函数结束时通知 WaitGroup 当前协程已完成

	// 构造请求参数
	param := &proto.GetGoodsDetailReq{
		GoodsId: int64(index + 1001), // 商品 ID 从 1001 开始
		UserId:  1,                   // 用户 ID 为 1
	}
	log.Printf("Sending request with GoodsId: %d", param.GoodsId) // 记录发送的请求

	// 调用 gRPC 服务的 GetGoodsDetail 方法
	startTime := time.Now()
	resp, err := client.GetGoodsDetail(context.Background(), param)
	duration := time.Since(startTime)
	if err != nil {
		log.Printf("Error calling GetGoodsDetail: %v", err) // 如果调用失败，记录错误
	} else {
		log.Printf("Response: %+v", resp) // 如果调用成功，记录响应
		log.Printf("Request for GoodsId: %d took %v", param.GoodsId, duration)
	}
}

// 测试 UpdateGoodsDetail 方法的函数
func TestUpdateGoodsDetail(wg *sync.WaitGroup, index int) {
	defer wg.Done()

	// 构造请求参数
	param := &proto.UpdateGoodsDetailReq{
		GoodsId: int64(index + 1001), // 商品 ID 从 1001 开始
		Price:   int64(index * 200),  // 更新后的销售价格（单位：分）
	}
	log.Printf("Sending request to update GoodsId: %d with new price: %d", param.GoodsId, param.Price) // 记录发送的请求

	// 调用 gRPC 服务的 UpdateGoodsDetail 方法
	resp, err := client.UpdateGoodsDetail(context.Background(), param)
	if err != nil {
		log.Printf("Error calling UpdateGoodsDetail: %v", err) // 如果调用失败，记录错误
	} else {
		log.Printf("Update response: %+v", resp) // 如果调用成功，记录响应
	}
}
func main() {
	defer conn.Close()    // 程序结束时关闭 gRPC 客户端连接
	var wg sync.WaitGroup // 使用 WaitGroup 等待所有协程完成

	// 启动多个协程测试 GetGoodsDetail 方法
	for i := 0; i < 5; i++ {
		wg.Add(1) // 增加 WaitGroup 的计数
		//go TestGetGoodsDetail(&wg, i) // 启动协程调用 TestGetGoodsDetail 函数
		go TestUpdateGoodsDetail(&wg, i)
	}
	wg.Wait() // 等待所有协程完成
}

var num int = 100 // 全局变量，初始值为 100（未在代码中使用）
//这是一个简单的并发测试函数，用于模拟多个协程对全局变量 num 的操作�?
//num = num - 1：每次调用减 1。由于没有加锁保护，这会导致并发问题（如竞态条件）
