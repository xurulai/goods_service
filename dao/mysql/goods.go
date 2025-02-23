package mysql

import (
	"context"
	"goods_srv/errno"
	"goods_srv/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// dao 层用来执行数据库相关的操作

// GetGoodsByRoomId 根据roomID查询直播间绑定的所有商品信息
func GetGoodsByRoomId(ctx context.Context, roomId int64) ([]*model.RoomGoods, error) {
	// 定义一个切片变量 data，用于存储查询结果
	// model.RoomGoods 是一个结构体，表示直播间与商品的绑定关系
	var data []*model.RoomGoods

	// 使用 gorm 的 WithContext 方法，将上下文传递给数据库操作
	// 确保数据库操作可以正确处理超时、取消等操作
	err := db.WithContext(ctx).
		// 指定操作的模型，这里操作的是 model.RoomGoods 表
		Model(&model.RoomGoods{}).
		// 添加查询条件，过滤出 room_id 等于传入的 roomId 的记录
		Where("room_id = ?", roomId).
		// 按照权重字段（weight）排序，确保返回的结果有序
		Order("weight").
		// 执行查询操作，将结果存储到 data 中
		Find(&data).Error

	// 如果查询出错且不是空数据的错误
	if err != nil && err != gorm.ErrEmptySlice {
		// 返回一个自定义的错误，表示查询失败
		return nil, errno.ErrQueryFailed
	}

	// 如果查询成功（无论是否有数据），返回查询结果
	return data, nil
}

// GetGoodsByIdList据id列表批量查询商品详情
func GetGoodsByIdList(ctx context.Context, idList []int64) ([]*model.Goods, error) {
	// 定义一个切片变量 data，用于存储查询结果
	// model.Goods 是一个结构体，表示商品信息
	var data []*model.Goods

	// 使用 gorm 的 WithContext 方法，将上下文传递给数据库操作
	// 确保数据库操作可以正确处理超时、取消等操作
	err := db.WithContext(ctx).
		// 指定操作的模型，这里操作的是 model.Goods 表
		Model(&model.Goods{}).
		// 添加查询条件，过滤出 goods_id 在 idList 中的记录
		Where("goods_id in ?", idList).
		// 使用 Clauses 方法添加自定义的排序逻辑
		// 确保查询结果按照 idList 中的顺序返回
		Clauses(clause.OrderBy{
			Expression: clause.Expr{
				SQL:                "FIELD(goods_id,?)",   // 使用 MySQL 的 FIELD 函数进行排序
				Vars:               []interface{}{idList}, // 传入 idList 作为排序参数
				WithoutParentheses: true,                  // 不需要额外的括号
			},
		}).
		// 执行查询操作，将结果存储到 data 中
		Find(&data).Error

	// 如果查询出错且不是空数据的错误
	if err != nil && err != gorm.ErrEmptySlice {
		// 返回一个自定义的错误，表示查询失败
		return nil, errno.ErrQueryFailed
	}

	// 如果查询成功（无论是否有数据），返回查询结果
	return data, nil
}

// GetGoods ById据id查询商品信息
func GetGoodsDetailById(ctx context.Context, goodsId int64) (*model.Goods, error) {
	// 定义一个切片变量 data，用于存储查询结果
	// model.Goods 是一个结构体，表示商品信息
	var data = &model.Goods{}

	// 使用 gorm 的 WithContext 方法，将上下文传递给数据库操作
	// 确保数据库操作可以正确处理超时、取消等操作
	err := db.WithContext(ctx).
		// 指定操作的模型，这里操作的是 model.Goods 表
		Model(&model.Goods{}).
		// 添加查询条件，过滤出 goods_id 在 idList 中的记录
		Where("goods_id = ?", goodsId).
		// 执行查询操作，将结果存储到 data 中
		First(data).Error

	// 如果查询出错且不是空数据的错误
	if err != nil && err != gorm.ErrEmptySlice {
		// 返回一个自定义的错误，表示查询失败
		return nil, errno.ErrQueryFailed
	}

	// 如果查询成功（无论是否有数据），返回查询结果
	return data, nil
}
