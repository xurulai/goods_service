package model

import (
	"time"
)

//BaseModel 是一个通用的基础模型，用于为所有数据表提供统一的字段结构和行为。

type BaseModel struct {
	ID       uint `gorm:"primaryKey"` // 主键ID
	CreateAt time.Time                // 创建时间
	UpdateAt time.Time                // 更新时间
	CreateBy string                   // 创建者
	UpdateBy string                   // 更新者
	Version  int16                    // 乐观锁版本号
	isDel    int8 `gorm:"index"`      // 软删除标志
}