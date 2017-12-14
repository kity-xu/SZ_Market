// 融资融券表
package finchina

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	"haina.com/share/models"
)

type TQ_SK_FINMRGNTRADE struct {
	models.Model `db:"-" `
	TRADEDATE    dbr.NullInt64
	FINBALANCE   dbr.NullFloat64
	MRGNRESQTY   dbr.NullFloat64
	FINMRGHBAL   dbr.NullFloat64
}

func NewTQ_SK_FINMRGNTRADE() *TQ_SK_FINMRGNTRADE {
	return &TQ_SK_FINMRGNTRADE{
		Model: models.Model{
			CacheKey:  "redis_key",
			TableName: TABLE_TQ_SK_FINMRGNTRADE,
			Db:        models.MyCat,
		},
	}
}

func (this *TQ_SK_FINMRGNTRADE) GetSMTFromFC(count int, which string) ([]TQ_SK_FINMRGNTRADE, error) {
	var SMT []TQ_SK_FINMRGNTRADE

	exps := map[string]interface{}{
		"EXCHANGE=?": which,
		"ISVALID=?":  1,
	}
	builder := this.Db.Select("TRADEDATE,FINBALANCE,MRGNRESQTY,FINMRGHBAL").From(this.TableName)
	_, err := this.SelectWhere(builder, exps).OrderBy("TRADEDATE desc").Limit(uint64(count)).LoadStructs(&SMT)
	if err != nil {
		logging.Error("%s", err.Error())
		return SMT, err
	}
	return SMT, err
}
