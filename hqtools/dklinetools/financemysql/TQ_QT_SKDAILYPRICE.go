package financemysql

import (
	"github.com/gocraft/dbr"
)

type Stock struct {
	TRADEDATE    dbr.NullInt64   `db:"TRADEDATE"`    // 交易日期
	SECODE       dbr.NullInt64   `db:"SECODE"`       // 证券内码
	EXCHANGE     dbr.NullInt64   `db:"EXCHANGE"`     // 交易市场
	LCLOSE       dbr.NullFloat64 `db:"LCLOSE"`       // 前收盘价
	TOPEN        dbr.NullFloat64 `db:"TOPEN"`        // 开盘价
	TCLOSE       dbr.NullFloat64 `db:"TCLOSE"`       // 收盘价
	THIGH        dbr.NullFloat64 `db:"THIGH"`        // 最高价
	TLOW         dbr.NullFloat64 `db:"TLOW"`         // 最低价
	VOL          dbr.NullInt64   `db:"VOL"`          // 成交量
	AMOUNT       dbr.NullFloat64 `db:"AMOUNT"`       // 成交金额
	DEALS        dbr.NullInt64   `db:"DEALS"`        // 成交笔数
	AVGPRICE     dbr.NullFloat64 `db:"AVGPRICE"`     // 当日均价
	AVGVOL       dbr.NullFloat64 `db:"AVGVOL"`       // 平均每笔成交量
	AVGTRAMT     dbr.NullFloat64 `db:"AVGTRAMT"`     // 平均每笔成交金额
	CHANGE       dbr.NullFloat64 `db:"CHANGE"`       // 涨跌
	PCHG         dbr.NullFloat64 `db:"PCHG"`         // 涨跌幅
	AMPLITUDE    dbr.NullFloat64 `db:"AMPLITUDE"`    // 振幅
	NEGOTIABLEMV dbr.NullFloat64 `db:"NEGOTIABLEMV"` // 流通市值
	TOTMKTCAP    dbr.NullFloat64 `db:"TOTMKTCAP"`    // 总市值
	TURNRATE     dbr.NullFloat64 `db:"TURNRATE"`     // 换手率
	ISVALID      dbr.NullInt64   `db:"ISVALID"`      // 是否有效

}

func (this *Stock) GetSKTListFC(sess1 *dbr.Session, secode string) ([]Stock, error) {
	var stock []Stock
	_, err := sess1.Select("*").From("TQ_QT_SKDAILYPRICE").
		Where("SECODE =" + secode).
		Where("TCLOSE > 0").
		Where("ISVALID=1").
		OrderBy("TRADEDATE").LoadStructs(&stock)
	return stock, err
}
func (this *Stock) GetSKTList5FC(sess *dbr.Session, secode string) ([]Stock, error) {
	var stock []Stock
	_, err := sess.Select(" TRADEDATE,VOL").From("TQ_QT_SKDAILYPRICE").
		Where("SECODE =" + secode).
		Where("ISVALID=1").
		OrderBy("TRADEDATE desc").Limit(5).LoadStructs(&stock)
	return stock, err
}
