package fcmysql

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	. "haina.com/share/models"
)

const (
	TABLE_TQ_SK_XDRY = "TQ_SK_XDRY" // 证券内码表
)

type XDRYFactor struct {
	BEGINDATE dbr.NullInt64   //起始日期VARCHAR(8) 	           本次除权因子的有效起始日期（即实际上的除权除息日）
	ENDDATE   dbr.NullInt64   //截止日期VARCHAR(8)              本次除权因子的有效截止日期（当尚无下一次除权因素的具体日期前，为19000101）
	XDY       dbr.NullFloat64 //当次除权因子NUMERIC(32,19)	   本次除权日，因分红送股转增等因素，依照除权前后价值不变动的原则计算的当次除权的折价因子
	LTDXDY    dbr.NullFloat64 //逆推累积除权因子NUMERIC(29,16)   以当前最新一天交易价格为标准，计算每次时间区间的累积除权因子，既每天实际交易价格与逆推复权价格之间的比值关系
	THELTDXDY dbr.NullFloat64 //顺推累计除权因子NUMERIC(29,16)   以上市第一天为标准，计算每次时间区间的累积除权因子，既顺推复权价格与每天实际交易价格与之间的比值关系
}

type TQ_SK_XDRY struct {
	Model `db:"-" `
}

func NewTQ_SK_XDRY() *TQ_SK_XDRY {
	return &TQ_SK_XDRY{
		Model: Model{
			TableName: TABLE_TQ_SK_XDRY,
			Db:        MyCat,
		},
	}
}

func (this *TQ_SK_XDRY) GetFactorsBySecode(secode string) ([]*XDRYFactor, error) {
	rows := make([]*XDRYFactor, 0, 100)

	cond := this.Db.Select("*").
		From(this.TableName).
		Where("SECODE='" + secode + "'").
		Where("ISVALID=1").
		Where("DATASOURCE=2").
		OrderBy("BEGINDATE")

	_, err := cond.LoadStructs(&rows)
	if err != nil {
		logging.Error("%v", err)
		return nil, err
	}

	return rows, nil
}
