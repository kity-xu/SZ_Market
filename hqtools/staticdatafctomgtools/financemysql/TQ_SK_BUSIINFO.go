package financemysql

import (
	"github.com/gocraft/dbr"
)

// 数据对象名称：TQ_SK_BUSIINFO    中文名称：公司经营情况

type TQ_SK_BUSIINFO struct {
	TCOREBIZINCOME dbr.NullFloat64 `db:"TCOREBIZINCOME"` // 本期主营业务收入
	TCOREBIZPROFIT dbr.NullFloat64 `db:"TCOREBIZPROFIT"` // 本期主营业务利润
}

// 查询公司业绩报表
func (this *TQ_SK_BUSIINFO) GetSingleInfo(sess *dbr.Session, comc string) (TQ_SK_BUSIINFO, error) {
	var tss TQ_SK_BUSIINFO
	err := sess.Select("*").From("TQ_SK_BUSIINFO").
		Where("COMPCODE='" + comc + "' and  ISVALID=1").OrderBy("PUBLISHDATE DESC").Limit(1).LoadStruct(&tss)
	return tss, err
}
