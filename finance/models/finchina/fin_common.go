// F10 财务分析接口共用
package finchina

import (
	. "haina.com/share/models"

	"haina.com/share/gocraft/dbr"
)

type RequestParam struct {
	SCodeOrigin string // 保留原参数scode
	SCode       string // 截取股票代码数字部分后的scode
	Type        int    // 类型(1:一季报 2:中报 3:三季报 4:年报)
	PerPage     int    // 每页条数,默认4
	Page        int    // 第几页的页码,默认1
}

type Responser interface {
	GetJson(*RequestParam) (*ResponseFinAnaJson, error)
}

type ResponseFinAnaJson struct {
	SCode  string      `json:"scode"`
	MU     string      `json:"MU"`
	AS     string      `json:"AS"`
	Length int         `json:"length"`
	List   interface{} `json:"list"`
}

type Session struct {
	Responser
	*ResponseFinAnaJson
}

func (this *ResponseFinAnaJson) NewSession(res Responser) *Session {
	return &Session{
		ResponseFinAnaJson: this,
		Responser:          res,
	}
}

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
