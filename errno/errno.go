package errno

import "errors"

var (
	ErrQueryFailed       = errors.New("query db failed")
	ErrGoodsDetailNull   = errors.New("query goodsdetail null")
	ErrUpdateFailed      = errors.New("update goodsdetail failed")
	ErrCacheDeleteFailed = errors.New("delete cache failed")
	ErrGoodsDetailNotFound = errors.New("found goodsdetail failed")
	ErrGetLockFailed = errors.New("get lock failed")
)
