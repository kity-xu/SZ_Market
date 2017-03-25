package company
import (
	"haina.com/market/finance/models/finchina"
)

/**
  获取股本结构信息
*/
func GetStructure(sCode string) (finchina.RetTrucInfoJson, error) {
	return finchina.NewSharestruchg().GetSingleByExps(sCode)
}

/**
  获取股本变动信息
*/
func GetChangesStrInfo(enddate string, sCode string, limit int) (finchina.RetShaInfoJson, error) {
	return finchina.NewChangesEquity().GetChangesStrJson(enddate, sCode, limit)
}
