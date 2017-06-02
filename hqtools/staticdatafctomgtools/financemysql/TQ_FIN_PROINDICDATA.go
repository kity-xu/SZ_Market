package financemysql

import (
	"github.com/gocraft/dbr"
)

// 数据对象名称：TQ_FIN_PROINDICDATA    中文名称：衍生财务指标（产品表）

type TQ_FIN_PROINDICDATA struct {
	CURRENTRT dbr.NullFloat64 `db:"CURRENTRT"` //  流动比率
	QUICKRT   dbr.NullFloat64 `db:"QUICKRT"`   // 速动比率
}

// 查询公司业绩报表
func (this *TQ_FIN_PROINDICDATA) GetSingleInfo(sess *dbr.Session, comc string) (TQ_FIN_PROINDICDATA, error) {
	var tss TQ_FIN_PROINDICDATA
	err := sess.Select("CURRENTRT,QUICKRT").From("TQ_FIN_PROINDICDATA").
		Where("COMPCODE='" + comc + "' and  ISVALID=1").
		OrderBy("ENDDATE DESC").
		Limit(1).
		LoadStruct(&tss)
	return tss, err
}
