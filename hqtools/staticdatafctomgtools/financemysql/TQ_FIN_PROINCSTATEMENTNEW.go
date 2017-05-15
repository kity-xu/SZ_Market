package financemysql

import (
	"github.com/gocraft/dbr"
)

// 数据对象名称：TQ_FIN_PROINCSTATEMENTNEW    中文名称：一般企业利润表(新准则产品表)

type TQ_FIN_PROINCSTATEMENTNEW struct {
	BASICEPS dbr.NullFloat64 `db:"BASICEPS"` // 基本每股收益
	//TOTASSET    dbr.NullFloat64 `db:"TOTASSET"`    // 总资产   // ？所查表中没有此字段
	TOTPROFIT   dbr.NullFloat64 `db:"TOTPROFIT"`   // 利润总额
	NETPROFIT   dbr.NullFloat64 `db:"NETPROFIT"`   // 净利润
	PUBLISHDATE dbr.NullInt64   `db:"PUBLISHDATE"` // 信息发布日期
}

// 查询公司业绩报表
func (this *TQ_FIN_PROINCSTATEMENTNEW) GetSingleInfo(sess *dbr.Session, comc string) (TQ_FIN_PROINCSTATEMENTNEW, error) {
	var tsp TQ_FIN_PROINCSTATEMENTNEW

	err := sess.Select("*").From("TQ_FIN_PROINCSTATEMENTNEW").
		Where("COMPCODE=" + comc + " and  ISVALID=1").
		OrderBy("ENDDATE DESC ").
		Limit(1).
		LoadStruct(&tsp)
	return tsp, err
}
