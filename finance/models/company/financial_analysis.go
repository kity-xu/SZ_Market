package company

//------------------------------------------------------------------------------
// F10 财务分析接口共用
type RespFinAnaJson struct {
	SCode  string      `json:"scode"`
	MU     string      `json:"MU"`
	AS     string      `json:"AS"`
	Length int         `json:"length"`
	List   interface{} `json:"list"`
}

type Responser interface {
	GetJson(scode string, report_type int, per_page int, page int) (*RespFinAnaJson, error)
}

type Session struct {
	Responser
	*RespFinAnaJson
}

func (this *RespFinAnaJson) NewSession(res Responser) *Session {
	return &Session{
		RespFinAnaJson: this,
		Responser:      res,
	}
}

//------------------------------------------------------------------------------
