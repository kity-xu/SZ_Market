package financemysql

import (
	"github.com/gocraft/dbr"
)

// 数据对象名称：TQ_FIN_PROINCSTATEMENTNEW    中文名称：一般企业利润表(新准则产品表)

type TQ_FIN_PROINCSTATEMENTNEW struct {
	BIZINCO   dbr.NullFloat64 `db:"BIZINCO"`   // 本期主营业务收入
	PERPROFIT dbr.NullFloat64 `db:"PERPROFIT"` // 本期主营业务利润
}

// 查询公司业绩报表
func (this *TQ_FIN_PROINCSTATEMENTNEW) GetSingleInfo(sess *dbr.Session, comc string) (TQ_FIN_PROINCSTATEMENTNEW, error) {
	var tss TQ_SK_BUSIINFO
	err := sess.Select("BIZINCO, PERPROFIT").From("TQ_FIN_PROINCSTATEMENTNEW").
		Where("COMPCODE='" + comc + "' and  ISVALID=1 and REPORTTYPE=1").
		OrderBy("ENDDATE DESC").
		Limit(1).
		LoadStruct(&tss)
	return tss, err
}
