package company

import (
	"haina.com/market/finance/models/finchina"
)

/**
  获取股东人数信息
*/
func GetShareholderList(enddate string, sCode string, limit int) (finchina.RetInfoJson, error) {
	return finchina.NewShareHolder().GetListByExps(enddate, sCode, limit)
}

/**
  获取十大流通股东信息
*/
func GetTop10List(enddate string, sCode string, limit int) (finchina.RetTopInfoJson, error) {
	return finchina.NewTop10().GetTop10List(enddate, sCode, limit)
}
