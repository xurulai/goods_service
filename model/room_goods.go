package model

// RoomGoods 直播间商品模型
type RoomGoods struct {
	BaseModel // 继承基础模型，包含通用字段

	RoomId    int64 `gorm:"notNull"`    // 直播间ID，关联直播间表
	GoodsId   int64 `gorm:"notNull"`    // 商品ID，关联商品表
	Weight    int64 `gorm:"notNull"`    // 权重，用于排序或推荐逻辑
	IsCurrent int8  `gorm:"is_current"` // 当前是否为直播间讲解的商品（1 表示是，0 表示不是）
}

// TableName 定义表名
func (RoomGoods) TableName() string {
	return "xx_room_goods" // 返回表名，这里可以根据实际需求修改表名
}
