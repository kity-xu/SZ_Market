// TQ_EXPT_SKSTATN    中文名称：个股一致预期_一致预期表(新)
package io_finchina

import (
	. "haina.com/share/models"

	"haina.com/market/hqpublish/models/finchina"
	"haina.com/share/gocraft/dbr"
)

type TQ_EXPT_SKSTATN struct {
	Model      `db:"-"`
	TENDDATE   dbr.NullString  // 预测截止日期（一期）
	TEPS       dbr.NullFloat64 // 一致预期每股收益(一期）
	TPE        dbr.NullFloat64 // 一致预期市盈率(PE)(一期)
	NENDDATE   dbr.NullString  // 预测截止日期（二期）
	NEPS       dbr.NullFloat64 // 一致预期每股收益（二期）
	NPE        dbr.NullFloat64 // 一致预期市盈率(PE)（二期)
	YANENDDATE dbr.NullString  // 预测截止日期（三期）
	YANEPS     dbr.NullFloat64 // 一致预期每股收益（三期）
	YANPE      dbr.NullFloat64 // 一致预期市盈率(PE)（三期)
}

func NewTQ_EXPT_SKSTATN() *TQ_EXPT_SKSTATN {
	return &TQ_EXPT_SKSTATN{
		Model: Model{
			TableName: finchina.TABLE_TQ_EXPT_SKSTATN,
			Db:        MyCat,
		},
	}
}

func (this *TQ_EXPT_SKSTATN) getSingleBySymbolAndExchange(symbol string, exchange string) error {

	builder := this.Db.Select("TENDDATE", "TEPS", "TPE", "NENDDATE", "NEPS", "NPE", "YANENDDATE", "YANEPS", "YANPE").From(this.TableName)

	err := builder.
		Where("SYMBOL=?", symbol).
		Where("EXCHANGE=?", exchange).
		Where("ISVALID=?", 1).
		OrderBy("PUBLISHDATE DESC").
		Limit(1).
		LoadStruct(this)
	if err != nil { //&& err != dbr.ErrNotFound
		return err
	}

	return nil
}

//------------------------------------------------------------------------------

func (this *TQ_EXPT_SKSTATN) GetSingle(symbol string, market string) error {
	exchange, err := this.getExchange(market)
	if err != nil {
		return err
	}

	return this.getSingleBySymbolAndExchange(symbol, exchange)
}

//
func (this *TQ_EXPT_SKSTATN) getExchange(market string) (string, error) {
	/*
		CNSESH		上交所
		CNSESZ		深交所
		STAS00		股份转让市场
	*/
	exchange := ""
	switch market {
	case "100":
		exchange = "CNSESH" //  上交所
	case "200":
		exchange = "CNSESZ" //  深交所
	default:
		return "", finchina.ErrMarket
	}
	return exchange, nil
}
