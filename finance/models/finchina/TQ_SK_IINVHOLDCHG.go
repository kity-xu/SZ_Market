package finchina

import (
	//	"haina.com/share/logging"

	//"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"
)

/**
  机构持股明细表
*/

type TQ_SK_IINVHOLDCHG struct {
	Model `db:"-" `
}

func NewTQ_SK_IINVHOLDCHG() *TQ_SK_IINVHOLDCHG {
	return &TQ_SK_IINVHOLDCHG{
		Model: Model{
			TableName: TABLE_TQ_SK_IINVHOLDCHG,
			Db:        MyCat,
		},
	}
}

// 获取机构持股数
func (this *TQ_SK_IINVHOLDCHG) GetInstitutionStockNum(compCode, reportDate string) float64 {
	exps := map[string]interface{}{
		"GPCOMPCODE=?": compCode,
		"REPORTDATE=?": reportDate,
		"ISVALID=?":    1,
	}

	var num float64

	builder := this.Db.Select("SUM(HOLDQTY)").From(this.TableName) //变动起始日
	err := this.SelectWhere(builder, exps).LoadStruct(&num)
	if err != nil {
		return float64(0)
	}
	return num
}
