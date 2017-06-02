package financemysql

import (
	"github.com/gocraft/dbr"
)

// 数据对象名称：TQ_SK_BASICINFO    中文名称：股票基本信息表

type TQ_SK_BASICINFO struct {
	SYMBOL     dbr.NullString `db:"SYMBOL"`     // 证券代码
	SETYPE     dbr.NullString `db:"SETYPE"`     // 证券类别
	EXCHANGE   dbr.NullString `db:"EXCHANGE"`   // 市场代码
	LISTSTATUS dbr.NullString `db:"LISTSTATUS"` // 上市状态
	LISTDATE   dbr.NullString `db:"LISTDATE"`   // 上市日期
	DELISTDATE dbr.NullString `db:"DELISTDATE"` // 退市日期
}

// 查询证券信息
func (this *TQ_SK_BASICINFO) GetBasicinfoList(sess *dbr.Session, symb string) (TQ_SK_BASICINFO, error) {
	var tsb TQ_SK_BASICINFO
	err := sess.Select("SYMBOL,SETYPE,EXCHANGE,LISTSTATUS,LISTDATE,DELISTDATE").From("TQ_SK_BASICINFO").
		Where("SYMBOL='" + symb + "' and  ISVALID=1").
		Limit(1).
		LoadStruct(&tsb)
	return tsb, err
}
