package fcmysql

import (
	_ "github.com/go-sql-driver/mysql"
	"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"
)

// 数据对象名称：TQ_SK_PERFORMANCE    中文名称：上市公司业绩快报表

type TQ_SK_PERFORMANCE struct {
	Model       `db:"-"`
	BASICEPS    dbr.NullFloat64 `db:"TOTASSET"`    // 基本每股收益
	WEIGHTEDROE dbr.NullFloat64 `db:"WEIGHTEDROE"` // 净资产收益率(加权)
	TOTSHAREQUI dbr.NullFloat64 `db:"TOTSHAREQUI"` // 股东权益合计
	NAPS        dbr.NullFloat64 `db:"NAPS"`        // 每股净资产
}

func NewTQ_SK_PERFORMANCE() *TQ_SK_PERFORMANCE {
	return &TQ_SK_PERFORMANCE{
		Model: Model{
			TableName: TABLE_TQ_SK_PERFORMANCE,
			Db:        MyCat,
		},
	}
}

// 查询公司业绩报表
func (this *TQ_SK_PERFORMANCE) GetSingleInfo(comc string) (TQ_SK_PERFORMANCE, error) {
	var tsp TQ_SK_PERFORMANCE

	err := this.Db.Select("BASICEPS,WEIGHTEDROE,TOTSHAREQUI,NAPS").
		From(this.TableName).
		Where("COMPCODE=" + comc + " and  ISVALID=1").
		OrderBy("PUBLISHDATE  DESC ").
		Limit(1).
		LoadStruct(&tsp)
	return tsp, err
}
