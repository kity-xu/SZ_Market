package financemysql

import (
	"github.com/gocraft/dbr"
)

// 数据对象名称：TQ_FIN_PROINCSTATEMENTNEW    中文名称：一般企业利润表(新准则产品表)

type TQ_FIN_PROINCSTATEMENTNEW struct {
	BASICEPS    dbr.NullFloat64 `db:"BASICEPS"`    // 基本每股收益
	TOTPROFIT   dbr.NullFloat64 `db:"TOTPROFIT"`   // 利润总额
	NETPROFIT   dbr.NullFloat64 `db:"NETPROFIT"`   // 净利润
	PUBLISHDATE dbr.NullInt64   `db:"PUBLISHDATE"` // 信息发布日期
	BIZINCO     dbr.NullFloat64 `db:"BIZINCO"`     // 主营业务收入
	PERPROFIT   dbr.NullFloat64 `db:"PERPROFIT"`   // 主营业务利润
	INVEINCO    dbr.NullFloat64 `db:"INVEINCO"`    // 投资收益
}

func (this *TQ_FIN_PROINCSTATEMENTNEW) GetSingleInfo(sess *dbr.Session, comc string) (TQ_FIN_PROINCSTATEMENTNEW, error) {
	var tsp TQ_FIN_PROINCSTATEMENTNEW

	err := sess.Select("BASICEPS,TOTPROFIT,NETPROFIT,PUBLISHDATE,BIZINCO,PERPROFIT,INVEINCO").From("TQ_FIN_PROINCSTATEMENTNEW").
		Where("COMPCODE=" + comc + " and  ISVALID=1").
		Where("REPORTTYPE=1").
		OrderBy("ENDDATE DESC ").
		Limit(1).
		LoadStruct(&tsp)
	return tsp, err
}
