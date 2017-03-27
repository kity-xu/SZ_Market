// F10 财务分析接口共用
package company

/*
type RequestParam struct {
	SCodeOrigin string // 保留原参数scode
	SCode       string // 截取股票代码数字部分后的scode
	Type        int    // 类型(1:一季报 2:中报 3:三季报 4:年报)
	PerPage     int    // 每页条数,默认4
	Page        int    // 第几页的页码,默认1
}

type ResponseFinAnaJson struct {
	SCode  string      `json:"scode"`
	MU     string      `json:"MU"`
	AS     string      `json:"AS"`
	Length int         `json:"length"`
	List   interface{} `json:"list"`
}

type Responser interface {
	GetJson(*RequestParam) (*ResponseFinAnaJson, error)
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
*/
