// 利润表
package company

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Profits struct {
	Date int64 `json:"Date"`

	AAPC string `json:"AAPC"`
	AILs string `json:"AILs"`
	AREp string `json:"AREp"`
	BAEp string `json:"BAEp"`
	BPAC string `json:"BPAC"`
	CoEp string `json:"CoEp"`
	CoRe string `json:"CoRe"`
	CORe string `json:"CORe"`
	DPES string `json:"DPES"`
	EPS  string `json:"EPS"`
	FnEp string `json:"FnEp"`
	ICEp string `json:"ICEp"`
	IDEp string `json:"IDEp"`
	InRe string `json:"InRe"`
	ItEp string `json:"ItEp"`
	ITEp string `json:"ITEp"`
	MgEp string `json:"MgEp"`
	MIIn string `json:"MIIn"`
	NCoE string `json:"NCoE"`
	NInR string `json:"NInR"`
	NOEp string `json:"NOEp"`
	NORe string `json:"NORe"`
	NtIn string `json:"NtIn"`
	OATx string `json:"OATx"`
	OCOR string `json:"OCOR"`
	OOCs string `json:"OOCs"`
	OpEp string `json:"OpEp"`
	OpPr string `json:"OpPr"`
	OpRe string `json:"OpRe"`
	SaEp string `json:"SaEp"`
	SAPC string `json:"SAPC"`
	TOpR string `json:"TOpR"`
	ToPr string `json:"ToPr"`
}

func NewProfits() *Profits {
	return &Profits{}
}

func (this *Profits) GetJson(c *gin.Context) (*ResponseInfo, error) {
	node1 := Profits{
		Date: time.Now().Unix(),
	}
	list := []Profits{}
	list = append(list, node1)
	res := &ResponseInfo{}
	res.List = list
	return res, nil
}
