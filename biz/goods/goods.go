package goods

import (
	"context"
	"encoding/json"
	"fmt"
	"goods_srv/dao/mysql"
	"goods_srv/dao/redis"
	"goods_srv/errno"
	"goods_srv/proto"
	"log"
	"math/rand"
	"time"
)

// biz层业务代码
// biz -> dao

// GetRoomGoodsListProto 根据直播间 ID 查询直播间绑定的所有商品信息，并组装成 protobuf 响应对象返回
func GetGoodsByRoom(ctx context.Context, roomId int64) (*proto.GoodsListResp, error) {
	// 1. 先去 xx_room_goods 表，根据 room_id 查询出所有的 goods_id
	objList, err := mysql.GetGoodsByRoomId(ctx, roomId)
	if err != nil {
		return nil, err // 如果查询失败，直接返回错误
	}

	// 处理数据
	// 1. 拿出所有的商品 ID
	// 2. 记住当前正在讲解的商品 ID
	var (
		currGoodsId int64                            // 当前正在讲解的商品 ID
		idList      = make([]int64, 0, len(objList)) // 存储所有商品 ID 的切片
	)

	// 遍历查询结果，提取商品 ID 和当前讲解的商品 ID
	for _, obj := range objList {
		fmt.Printf("obj:%#v\n", obj)         // 打印当前对象信息（调试用）
		idList = append(idList, obj.GoodsId) // 将商品 ID 添加到 idList 中
		if obj.IsCurrent == 1 {              // 如果当前对象是正在讲解的商品
			currGoodsId = obj.GoodsId // 记录当前正在讲解的商品 ID
		}
	}

	// 2. 再拿上面获取到的 goods_id 去 xx_goods 表查询所有的商品详细信息
	goodsList, err := mysql.GetGoodsByIdList(ctx, idList)
	if err != nil {
		return nil, err // 如果查询失败，直接返回错误
	}

	// 拼装响应数据
	data := make([]*proto.GoodsInfo, 0, len(goodsList)) // 创建一个存储商品信息的切片
	for _, goods := range goodsList {
		data = append(data, &proto.GoodsInfo{ // 创建一个 GoodsInfo 对象并添加到 data 切片中
			GoodsId:     goods.GoodsId,                                       // 商品 ID
			CategoryId:  goods.CategoryId,                                    // 商品分类 ID
			Status:      int32(goods.Status),                                 // 商品状态
			Title:       goods.Title,                                         // 商品标题
			MarketPrice: fmt.Sprintf("%.2f", float64(goods.MarketPrice/100)), // 商品市场价（单位转换为元）
			Price:       fmt.Sprintf("%.2f", float64(goods.Price/100)),       // 商品售价（单位转换为元）
			Brief:       goods.Brief,                                         // 商品简介
		})
	}

	// 创建并返回 protobuf 响应对象
	resp := &proto.GoodsListResp{
		CurrentGoodsId: currGoodsId, // 当前正在讲解的商品 ID
		Data:           data,        // 商品信息列表
	}
	return resp, nil
}

func GetGoodsDetailById(ctx context.Context, goodsId int64) (*proto.GoodsDetail, error) {
	// 构造缓存键，用于从 Redis 中查询商品详情
	cacheKey := fmt.Sprintf("goods_detail_%d", goodsId)

	// 1. 尝试从 Redis 缓存中获取商品详情
	cachedData, err := redis.GetClient().Get(ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		// 缓存命中
		log.Printf("Cache hit for GoodsId: %d", goodsId)
		var goodsDetail proto.GoodsDetail
		if err := json.Unmarshal([]byte(cachedData), &goodsDetail); err != nil {
			log.Printf("Failed to unmarshal cached data: %v", err)
			return nil, errno.ErrQueryFailed
		}
		return &goodsDetail, nil
	} else if err != nil {
		log.Printf("Failed to get data from cache: %v", err)
	} else {
		log.Printf("Cache miss for GoodsId: %d", goodsId)
	}

	// 缓存未命中，从数据库中查询商品详情
	// 1. 根据商品 ID 从 MySQL 数据库中查询商品详情
	goodsDetail, err := mysql.GetGoodsDetailById(ctx, goodsId)
	if err != nil {
		log.Printf("Failed to query goods detail: %v", err)
		return nil, errno.ErrQueryFailed
	}

	// 2. 检查查询结果是否为空
	if goodsDetail == nil {
		log.Printf("Goods detail not found for GoodsId: %d", goodsId)
		return nil, errno.ErrGoodsDetailNull
	}

	// 3. 检查商品详情数据是否有效（关键字段是否为空或为零）
	if goodsDetail.GoodsId == 0 || goodsDetail.Title == "" || goodsDetail.Price == 0 {
		log.Printf("Invalid goods detail data: %+v", goodsDetail)
		return nil, errno.ErrGoodsDetailNull
	}

	// 4. 构造返回的响应数据
	resp := &proto.GoodsDetail{
		GoodsId:    goodsDetail.GoodsId,
		CategoryId: goodsDetail.CategoryId,
		Status:     int32(goodsDetail.Status),
		Title:      goodsDetail.Title,
		Code:       goodsDetail.Code,      // 商品编码
		BrandName:  goodsDetail.BrandName, // 商品品牌名称
		Brief:      goodsDetail.Brief,
	}

	// 5. 格式化市场价格和价格字段
	if goodsDetail.MarketPrice > 0 {
		// 将市场价格除以 100 并格式化为两位小数的字符串
		resp.MarketPrice = fmt.Sprintf("%.2f", float64(goodsDetail.MarketPrice)/100)
	} else {
		resp.MarketPrice = "0.00"
		log.Printf("MarketPrice is zero or invalid for GoodsId: %d", goodsId)
	}

	if goodsDetail.Price > 0 {
		// 将价格除以 100 并格式化为两位小数的字符串
		resp.Price = fmt.Sprintf("%.2f", float64(goodsDetail.Price)/100)
	} else {
		resp.Price = "0.00"
		log.Printf("Price is zero or invalid for GoodsId: %d", goodsId)
	}

	// 6. 将查询结果序列化为 JSON 数据
	cachedBytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Failed to marshal data: %v", err)
		return nil, errno.ErrQueryFailed
	}

	// 7. 将查询结果存入 Redis 缓存，设置随机过期时间
	// 基础过期时间为 10 分钟，随机增加 0-5 分钟
	baseTTL := 10 * time.Minute
	randomTTL := time.Duration(rand.Intn(5*60)) * time.Second
	totalTTL := baseTTL + randomTTL
	_, err = redis.GetClient().Set(ctx, cacheKey, cachedBytes, totalTTL).Result()
	if err != nil {
		log.Printf("Failed to set data in cache: %v", err)
	}

	// 返回商品详情响应
	log.Printf("Returning goods detail response: %+v", resp)
	return resp, nil
}

// UpdateGoodsDetail 更新商品详情，并删除缓存
func UpdateGoodsDetail(ctx context.Context, goodsId int64, newPrice int64) (*proto.Response, error) {
	// 1. 更新数据库
	err := mysql.UpdateGoodsDetail(ctx, goodsId, newPrice)
	if err != nil {
		log.Printf("Failed to update goods detail: %v", err)
		return nil, errno.ErrUpdateFailed
	}

	// 2. 删除缓存
	cacheKey := fmt.Sprintf("goods_detail_%d", goodsId)
	_, err = redis.GetClient().Del(ctx, cacheKey).Result()
	if err != nil {
		log.Printf("Failed to delete cache: %v", err)
		return nil, errno.ErrCacheDeleteFailed
	}

	log.Printf("Cache deleted for GoodsId: %d", goodsId)
	return &proto.Response{}, nil
}
