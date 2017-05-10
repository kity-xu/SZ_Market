package financemysql

import (
	"github.com/gocraft/dbr"
)

type ComCode struct {
	SECODE     dbr.NullString `db:"SECODE"`     // 证券内码
	SYMBOL     dbr.NullString `db:"SYMBOL"`     // 证券代码
	COMPCODE   dbr.NullString `db:"COMPCODE"`   // 公司内码
	EXCHANGE   dbr.NullString `db:"EXCHANGE"`   // 交易市场
	SETYPE     dbr.NullString `db:"SETYPE"`     // 证券类别
	LISTSTATUS dbr.NullString `db:"LISTSTATUS"` // 上市状态
}

// 查询所有沪深所有股票代码
func (this *ComCode) GetComCodeList(sess *dbr.Session) ([]ComCode, error) {
	var code []ComCode
	_, err := sess.Select("*").From("TQ_OA_STCODE").
		Where("EXCHANGE in ('001002','001003') and SETYPE='101' and  LISTSTATUS=1 and ISVALID =1").
		OrderBy("SECODE").LoadStructs(&code)
	return code, err
}

// 查询指数信息
func (this *ComCode) GetIndexInfoList(sess *dbr.Session) ([]ComCode, error) {
	var code []ComCode
	_, err := sess.Select("*").From("TQ_OA_STCODE").
		Where("SETYPE ='701' AND (SYMBOL LIKE '399%' OR SYMBOL LIKE '000%') AND LISTSTATUS =1 and ISVALID =1").
		OrderBy("SYMBOL").LoadStructs(&code)
	return code, err
}
