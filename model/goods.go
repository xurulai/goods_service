package model

// ORM
// struct -> table
// Goods 商品模型
type Goods struct {
	BaseModel // 继承基础模型，包含通用字段

	GoodsId     int64  `gorm:"notNull;uniqueIndex"` // 商品ID，唯一标识一个商品
	CategoryId  int64  `gorm:"notNull"`             // 商品所属分类ID
	BrandName   string `gorm:"notNull"`             // 品牌名称
	Code        string `gorm:"notNull;uniqueIndex"` // 商品编码，唯一标识一个商品
	Status      int8   `gorm:"notNull"`             // 商品状态（例如：上架、下架、审核中等）
	Title       string `gorm:"notNull"`             // 商品标题
	MarketPrice int64  `gorm:"notNull"`             // 市场价
	Price       int64  `gorm:"notNull"`             // 实际销售价格
	Brief       string `gorm:"type:text"`           // 商品简介
}

// TableName 定义表名
func (Goods) TableName() string {
	return "xx_goods_query" // 返回表名，这里可以根据实际需求修改表名
}
