package fcmysql

import (
	"github.com/gocraft/dbr"
)

type TQ_OA_STCODE struct {
	SECODE     dbr.NullString `db:"SECODE"`     // 证券内码
	SENAME     dbr.NullString `db:"SENAME"`     // 证券全称
	SYMBOL     dbr.NullString `db:"SYMBOL"`     // 证券代码
	COMPCODE   dbr.NullString `db:"COMPCODE"`   // 公司内码
	EXCHANGE   dbr.NullString `db:"EXCHANGE"`   // 交易市场
	SETYPE     dbr.NullString `db:"SETYPE"`     // 证券类别
	LISTSTATUS dbr.NullString `db:"LISTSTATUS"` // 上市状态
}

// 查询所有沪深所有股票代码
func (this *TQ_OA_STCODE) GetComCodeList(sess *dbr.Session, str string) ([]TQ_OA_STCODE, error) {
	var code []TQ_OA_STCODE
	_, err := sess.Select("*").From("TQ_OA_STCODE").
		Where("ISVALID =1").
		Where("SETYPE=101").
		Where("SECODE in (" + str + ")").
		OrderBy("SECODE").LoadStructs(&code)
	return code, err
}
