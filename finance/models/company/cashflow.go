// 现金流量表
package company

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Cashflow struct {
	Date  int64  `json:"Date"`
	CGBOA string `json:"CGBOA"`
	CGIIA string `json:"CGIIA"`
	CPBFA string `json:"CPBFA"`
	CPBID string `json:"CPBID"`
	CPBLn string `json:"CPBLn"`
	CPFGS string `json:"CPFGS"`
	CPFSW string `json:"CPFSW"`
	CPFTx string `json:"CPFTx"`
	CRBIv string `json:"CRBIv"`
	CRFDI string `json:"CRFDI"`
	CRFGS string `json:"CRFGS"`
	CRFII string `json:"CRFII"`
	CRMSS string `json:"CRMSS"`
	CUIIA string `json:"CUIIA"`
	CUIIv string `json:"CUIIv"`
	CUIOA string `json:"CUIOA"`
	GDPES string `json:"GDPES"`
	IERCE string `json:"IERCE"`
	NBFBI string `json:"NBFBI"`
	NBFCB string `json:"NBFCB"`
	NCEAI string `json:"NCEAI"`
	NCEIS string `json:"NCEIS"`
	NCFOA string `json:"NCFOA"`
	NCIIA string `json:"NCIIA"`
	NCPFA string `json:"NCPFA"`
	NCRDU string `json:"NCRDU"`
	NCRFU string `json:"NCRFU"`
	NDCBI string `json:"NDCBI"`
	NIcLn string `json:"NIcLn"`
	NIICE string `json:"NIICE"`
	NLend string `json:"NLend"`
	NLnAv string `json:"NLnAv"`
	PcsPE string `json:"PcsPE"`
	PmFPy string `json:"PmFPy"`
	PmISA string `json:"PmISA"`
	PmoFA string `json:"PmoFA"`
}

func NewCashflow() *Cashflow {
	return &Cashflow{}
}

func (this *Cashflow) GetJson(c *gin.Context) (*ResponseInfo, error) {
	node1 := Cashflow{
		Date: time.Now().Unix(),
	}
	list := []Cashflow{}
	list = append(list, node1)
	res := &ResponseInfo{}
	res.List = list
	return res, nil
}
