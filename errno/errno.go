package errno

import "errors"

var (
	ErrQueryFailed = errors.New("query db failed")
	ErrGoodsDetailNull = errors.New("query goodsdetail null")
)
