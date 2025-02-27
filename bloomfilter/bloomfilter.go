package bloomfilter

import (
	"context"
	"fmt"
	"goods_srv/dao/mysql"
	"log"

	"github.com/willf/bloom"
)

var (
	goodsbloomfiltyer *bloom.BloomFilter //本地布隆过滤器实例
)

func InitBloomFilter(ctx context.Context) error {
	// 预计插入的商品ID数量和误判率
	estimatedItems := 1000000 // 预计商品总数
	errorRate := 0.0001       // 误判率 0.01%
	goodsbloomfiltyer = bloom.New(uint(estimatedItems), uint(errorRate))

	//从数据库加载所有商品ID
	goodsIDs, err := mysql.GetAllGoodsIDs(ctx)
	if err != nil {
		log.Fatalf("Failed to load goods IDs from database: %v", err)
	}

	// 将商品ID添加到布隆过滤器中
	for _, goodsID := range goodsIDs {
		goodsbloomfiltyer.Add([]byte(fmt.Sprintf("%d", goodsID)))
	}

	log.Printf("Bloom filter initialized with %d goods IDs", len(goodsIDs))
	return nil
}
