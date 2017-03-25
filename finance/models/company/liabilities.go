// 资产负债表
package company

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Liabilities struct {
	Date int64 `json:"Date"`

	//资产
	AcRe string `json:"AcRe"`
	CuDe string `json:"CuDe"`
	DeMg string `json:"DeMg"`
	DFAs string `json:"DFAs"`
	DfSv string `json:"DfSv"`
	DiRe string `json:"DiRe"`
	DTAs string `json:"DTAs"`
	FAFS string `json:"FAFS"`
	FAHT string `json:"FAHT"`
	FiAs string `json:"FiAs"`
	GWil string `json:"GWil"`
	HTMI string `json:"HTMI"`
	IbDe string `json:"IbDe"`
	InAs string `json:"InAs"`
	InRe string `json:"InRe"`
	LdLt string `json:"LdLt"`
	LTAR string `json:"LTAR"`
	LTEI string `json:"LTEI"`
	LTPE string `json:"LTPE"`
	Metl string `json:"Metl"`
	MnFd string `json:"MnFd"`
	NoRe string `json:"NoRe"`
	OtRe string `json:"OtRe"`
	PrEx string `json:"PrEx"`
	Prpy string `json:"Prpy"`
	REFI string `json:"REFI"`
	ToAs string `json:"ToAs"`

	//负债
	AcEx  string `json:"AcEx"`
	AcPy  string `json:"AcPy"`
	AdRE  string `json:"AdRE"`
	BdPy  string `json:"BdPy"`
	CmPy  string `json:"CmPy"`
	DETLb string `json:"DETLb"`
	DfIn  string `json:"DfIn"`
	DFLb  string `json:"DFLb"`
	DpCl  string `json:"DpCl"`
	DpFB  string `json:"DpFB"`
	DvPy  string `json:"DvPy"`
	FASR  string `json:"FASR"`
	InPy  string `json:"InPy"`
	LnFB  string `json:"LnFB"`
	LnFC  string `json:"LnFC"`
	LTBw  string `json:"LTBw"`
	LTPy  string `json:"LTPy"`
	NCL1  string `json:"NCL1"`
	NtPy  string `json:"NtPy"`
	PCSc  string `json:"PCSc"`
	PlLn  string `json:"PlLn"`
	PrSk  string `json:"PrSk"`
	SaPy  string `json:"SaPy"`
	SBPy  string `json:"SBPy"`
	STLn  string `json:"STLn"`
	TaLb  string `json:"TaLb"`
	TFLb  string `json:"TFLb"`
	TxPy  string `json:"TxPy"`

	//所有者权益
	BPCOEAI string `json:"BPCOEAI"`
	BPCOESI string `json:"BPCOESI"`
	BPCSET  string `json:"BPCSET"`
	CDFCS   string `json:"CDFCS"`
	CpSp    string `json:"CpSp"`
	GRPr    string `json:"GRPr"`
	LEAI    string `json:"LEAI"`
	LESI    string `json:"LESI"`
	MiIt    string `json:"MiIt"`
	OEAI    string `json:"OEAI"`
	OEIn    string `json:"OEIn"`
	OESET   string `json:"OESET"`
	OtCI    string `json:"OtCI"`
	PCSe    string `json:"PCSe"`
	PICa    string `json:"PICa"`
	PrSc    string `json:"PrSc"`
	SpRs    string `json:"SpRs"`
	TLSE    string `json:"TLSE"`
	TrSc    string `json:"TrSc"`
	UdPr    string `json:"UdPr"`
}

func NewLiabilities() *Liabilities {
	return &Liabilities{}
}

func (this *Liabilities) GetJson(c *gin.Context) (*ResponseInfo, error) {
	node1 := Liabilities{
		Date: time.Now().Unix(),
	}
	list := []Liabilities{}
	list = append(list, node1)
	res := &ResponseInfo{}
	res.List = list

	return res, nil
}
