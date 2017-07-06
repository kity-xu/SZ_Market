package fcmysql

import (
	_ "github.com/go-sql-driver/mysql"
	. "haina.com/share/models"
)

// 数据对象名称：TQ_SK_INTERVALQT    中文名称：股票固定区间行情表

type TQ_SK_INTERVALQT struct {
	Model `db:"-"`
	VOL5D int64 `db:"VOL5D"` // 五日交易量
}

func NewTQ_SK_INTERVALQT() *TQ_SK_INTERVALQT {
	return &TQ_SK_INTERVALQT{
		Model: Model{
			TableName: TABLE_TQ_SK_INTERVALQT,
			Db:        MyCat,
		},
	}
}

// 查询证券信息
func (this *TQ_SK_INTERVALQT) GetSingleInfo(sec string) (TQ_SK_INTERVALQT, error) {
	var tss TQ_SK_INTERVALQT
	err := this.Db.Select("VOL5D").From(this.TableName).
		Where("SECODE='" + sec + "' and  ISVALID=1").
		OrderBy("TRADEDATE DESC").
		Limit(1).
		LoadStruct(&tss)
	return tss, err
}
