// TQ_EXPT_SKSTATRT    中文名称：个股一致预期_评级汇总表
package io_finchina

import (
	. "haina.com/share/models"

	"haina.com/market/hqpublish/models/finchina"
	"haina.com/share/gocraft/dbr"
)

type TQ_EXPT_SKSTATRT struct {
	Model      `db:"-"`
	BUY        dbr.NullInt64 // 买入
	ADDBUYNUM  dbr.NullInt64 // 增持
	NEUTRALNUM dbr.NullInt64 // 中性
	BREDUCENUM dbr.NullInt64 // 减持
	SELL       dbr.NullInt64 // 卖出
}

func NewTQ_EXPT_SKSTATRT() *TQ_EXPT_SKSTATRT {
	return &TQ_EXPT_SKSTATRT{
		Model: Model{
			TableName: finchina.TABLE_TQ_EXPT_SKSTATRT,
			Db:        MyCat,
		},
	}
}

func (this *TQ_EXPT_SKSTATRT) getSingleBySymbolAndExchange(symbol string, exchange string, period_type int) error {

	builder := this.Db.Select("BUY", "ADDBUYNUM", "NEUTRALNUM", "BREDUCENUM", "SELL").From(this.TableName)

	err := builder.
		Where("SYMBOL=?", symbol).
		Where("EXCHANGE=?", exchange).
		Where("PERIODTYPE=?", period_type).
		OrderBy("PUBLISHDATE DESC").
		Limit(1).
		LoadStruct(this)
	if err != nil && err != dbr.ErrNotFound {
		return err
	}

	return nil
}

//------------------------------------------------------------------------------

func (this *TQ_EXPT_SKSTATRT) GetSingle(symbol string, market string, period_type int) error {
	sc := NewTQ_OA_STCODE()
	exchange, err := sc.getExchange(market)
	if err != nil {
		return err
	}

	return this.getSingleBySymbolAndExchange(symbol, exchange, period_type)
}
