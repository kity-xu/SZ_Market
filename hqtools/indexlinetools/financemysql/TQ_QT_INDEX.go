package financemysql

import (
	"github.com/gocraft/dbr"
)

type TQ_QT_INDEX struct {
	TRADEDATE dbr.NullFloat64 `db:"TRADEDATE"` // 交易日期
	SECODE    dbr.NullFloat64 `db:"SECODE"`    // 证券内码
	EXCHANGE  dbr.NullFloat64 `db:"EXCHANGE"`  // 交易市场
	LCLOSE    dbr.NullFloat64 `db:"LCLOSE"`    // 前收盘价
	TOPEN     dbr.NullFloat64 `db:"TOPEN"`     // 开盘价
	TCLOSE    dbr.NullFloat64 `db:"TCLOSE"`    // 收盘价
	THIGH     dbr.NullFloat64 `db:"THIGH"`     // 最高价
	TLOW      dbr.NullFloat64 `db:"TLOW"`      // 最低价
	VOL       dbr.NullFloat64 `db:"VOL"`       // 成交量
	AMOUNT    dbr.NullFloat64 `db:"AMOUNT"`    // 成交金额
	DEALS     dbr.NullFloat64 `db:"DEALS"`     // 成交笔数
	TOTMKTCAP dbr.NullFloat64 `db:"TOTMKTCAP"` // 总市值
}

func (this *TQ_QT_INDEX) GetIndexInfoList(sess *dbr.Session, secode string) ([]TQ_QT_INDEX, error) {
	var index []TQ_QT_INDEX
	_, err := sess.Select("*").From("TQ_QT_INDEX").
		Where("SECODE =" + secode).
		OrderBy("TRADEDATE").LoadStructs(&index)
	return index, err
}
