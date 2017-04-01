// F10 财务分析接口共用
package finchina

import (
	. "haina.com/share/models"

	"haina.com/share/gocraft/dbr"
)

// ---------------------------------------------------------------------
type SymbolToCompcode struct {
	Model    `db:"-"`
	SYMBOL   string         //股票代码
	COMPCODE dbr.NullString //公司代码(公司内码) 通过 SYMBOL 得到
	CUR      dbr.NullString //货币单位
}

func NewSymbolToCompcode() *SymbolToCompcode {
	return &SymbolToCompcode{
		Model: Model{
			TableName: TABLE_TQ_OA_STCODE,
			Db:        MyCat,
		},
	}
}

func (this *SymbolToCompcode) getCompcode(symbol string) error {
	err := this.Db.Select("SYMBOL", "COMPCODE", "CUR").From(this.TableName).Where("SYMBOL = ?", symbol).Limit(1).LoadStruct(this)
	return err
}
