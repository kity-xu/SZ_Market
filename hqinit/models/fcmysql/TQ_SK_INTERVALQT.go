package fcmysql

import (
	"github.com/gocraft/dbr"
)

// 数据对象名称：TQ_SK_INTERVALQT    中文名称：股票固定区间行情表

type TQ_SK_INTERVALQT struct {
	VOL5D int64 `db:"VOL5D"` // 五日交易量
}

// 查询证券信息
func (this *TQ_SK_INTERVALQT) GetSingleInfo(sess *dbr.Session, sec string) (TQ_SK_INTERVALQT, error) {
	var tss TQ_SK_INTERVALQT
	err := sess.Select("VOL5D").From("TQ_SK_INTERVALQT").
		Where("SECODE='" + sec + "' and  ISVALID=1").
		OrderBy("ORDER BY TRADEDATE DESC").
		Limit(1).
		LoadStruct(&tss)
	return tss, err
}
