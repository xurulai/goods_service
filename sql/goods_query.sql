CREATE TABLE `xx_goods_query` (
                         `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
                         `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                         `create_by` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '创建者',
                         `update_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
                         `update_by` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '更新者',
                         `version` SMALLINT(5) UNSIGNED NOT NULL DEFAULT '0' COMMENT '乐观锁版本号',
                         `is_del` tinyint(4) UNSIGNED NOT NULL DEFAULT '0' COMMENT '是否删除：0正常1删除',
                         `user_id` BIGINT(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '用户id',
                         `goods_id` BIGINT(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '商品id',
                         `category_id` BIGINT(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '类目id',
                         `brand_name` VARCHAR(255) NOT NULL COMMENT '品牌名',
                         `code` VARCHAR(64) NOT NULL COMMENT '码',
                         `status` tinyint(4) UNSIGNED NOT NULL DEFAULT '0' COMMENT '是否上架：0上架1下架',
                         `title` VARCHAR(255) NOT NULL COMMENT '名称',
                         `market_price` BIGINT(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '市场价/划线价（分）',
                         `price` BIGINT(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '售价（分）',
                         `brief` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '简介',
                         UNIQUE (goods_id),
                         INDEX (category_id),
                         INDEX (is_del)
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT = '商品查询表';
-- 真正企业里面一般都是按照分类>SPU>SKU维度建表，这里只是方便课上教学。