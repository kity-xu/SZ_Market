package fcmysql

import (
	_ "github.com/go-sql-driver/mysql"
	"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"
)

// 数据对象名称：TQ_FIN_PROINDICDATA    中文名称：衍生财务指标（产品表）

type TQ_FIN_PROINDICDATA struct {
	Model     `db:"-"`
	UPPS      dbr.NullFloat64 `db:"UPPS"`      // 每股未分配利润
	CURRENTRT dbr.NullFloat64 `db:"CURRENTRT"` // 流动比率
	QUICKRT   dbr.NullFloat64 `db:"QUICKRT"`   // 速动比率
}

func NewTQ_FIN_PROINDICDATA() *TQ_FIN_PROINDICDATA {
	return &TQ_FIN_PROINDICDATA{
		Model: Model{
			TableName: TABLE_TQ_FIN_PROINDICDATA,
			Db:        MyCat,
		},
	}
}

// 查询公司业绩报表
func (this *TQ_FIN_PROINDICDATA) GetSingleInfo(comc string) (TQ_FIN_PROINDICDATA, error) {
	var tss TQ_FIN_PROINDICDATA
	err := this.Db.Select("CURRENTRT,QUICKRT,UPPS").
		From(this.TableName).
		Where("COMPCODE='" + comc + "'").
		Where("ISVALID=1 and REPORTTYPE=3").
		OrderBy("ENDDATE DESC").
		Limit(1).
		LoadStruct(&tss)
	return tss, err
}
