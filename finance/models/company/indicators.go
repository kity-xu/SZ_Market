// 关键指标
package company

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Indicators struct {
	Date int64 `json:"Date"`

	// 每股指标
	DEPS      string `json:"DEPS"`
	EPS       string `json:"EPS"`
	EPSE      string `json:"EPSE"`
	EPSTTM    string `json:"EPSTTM"`
	PSAF      string `json:"PSAF"`
	PSCR      string `json:"PSCR"`
	PSECF     string `json:"PSECF"`
	PSNA      string `json:"PSNA"`
	PSNBCF    string `json:"PSNBCF"`
	PSNBCFTTM string `json:"PSNBCFTTM"`
	PSNCF     string `json:"PSNCF"`
	PSNCFTTM  string `json:"PSNCFTTM"`
	PSOE      string `json:"PSOE"`
	PSR       string `json:"PSR"`
	PSRE      string `json:"PSRE"`
	PSRTTM    string `json:"PSRTTM"`
	PSSCF     string `json:"PSSCF"`
	PSSR      string `json:"PSSR"`
	PSTR      string `json:"PSTR"`
	PSUP      string `json:"PSUP"`

	// 盈利能力
	CPR       string `json:"CPR"`
	DSR       string `json:"DSR"`
	DSRTTM    string `json:"DSRTTM"`
	EBITDR    string `json:"EBITDR"`
	EBITDRTTM string `json:"EBITDRTTM"`
	FEDBR     string `json:"FEDBR"`
	FEDBRTTM  string `json:"FEDBRTTM"`
	JROA      string `json:"JROA"`
	JROATTM   string `json:"JROATTM"`
	LAIDBR    string `json:"LAIDBR"`
	LAIDBRTTM string `json:"LAIDBRTTM"`
	MEDBR     string `json:"MEDBR"`
	MEDBRTTM  string `json:"MEDBRTTM"`
	NPADNRGL  string `json:"NPADNRGL"`
	NPAPC     string `json:"NPAPC"`
	NPOR      string `json:"NPDBR"`
	NPORTTM   string `json:"NPDBRTTM"`
	NSR       string `json:"NSR"`
	NSRTTM    string `json:"NSRTTM"`
	OCDBR     string `json:"OCDBR"`
	OCDBRTTM  string `json:"OCDBRTTM"`
	OPR       string `json:"OPR"`
	ROEA      string `json:"ROEA"`
	ROED      string `json:"ROED"`
	ROEDD     string `json:"ROEDD"`
	ROEDW     string `json:"ROEDW"`
	ROETTM    string `json:"ROETTM"`
	ROEW      string `json:"ROEW"`
	SCR       string `json:"SCR"`
	SEDBR     string `json:"SEDBR"`
	SEDBRTTM  string `json:"SEDBRTTM"`
	SGPM      string `json:"SGPM"`
	SGPMTTM   string `json:"SGPMTTM"`

	// 偿债能力
	BPCSDIBD   string `json:"BPCSDIBD"`
	BPCSDTL    string `json:"BPCSDTL"`
	CFLR       string `json:"CFLR"`
	CR         string `json:"CR"`
	ER         string `json:"ER"`
	LTLRWC     string `json:"LTLRWC"`
	NBAGCFDCL  string `json:"NBAGCFDCL"`
	NBAGCFDIBD string `json:"NBAGCFDIBD"`
	NBAGCFDND  string `json:"NBAGCFDND"`
	NBAGCFDTL  string `json:"NBAGCFDTL"`
	QR         string `json:"QR"`
	SQR        string `json:"SQR"`
	TIE        string `json:"TIE"`
	TNWDIBD    string `json:"TNWDIBD"`
	TNWDND     string `json:"TNWDND"`
	TNWDR      string `json:"TNWDR"`

	// 成长能力
	APCSNPYG      string `json:"APCSNPYG"`
	APCSNPYGD     string `json:"APCSNPYGD"`
	BAGCFNYOYG    string `json:"BAGCFNYOYG"`
	BAGCFPSNYOYG  string `json:"BAGCFPSNYOYG"`
	BEPSYG        string `json:"BEPSYG"`
	BIYG          string `json:"BIYG"`
	BPYG          string `json:"BPYG"`
	BSPCERBYGR    string `json:"BSPCERBYGR"`
	DEPSYG        string `json:"DEPSYG"`
	NAPSRBYGR     string `json:"NAPSRBYGR"`
	NAYG          string `json:"NAYG"`
	NPYG          string `json:"NPYG"`
	OP5YSPBPCNPGA string `json:"OP5YSPBPCNPGA"`
	REDYG         string `json:"REDYG"`
	SGR           string `json:"SGR"`
	TARBYGR       string `json:"TARBYGR"`
	TAYG          string `json:"TAYG"`
	TPYG          string `json:"TPYG"`

	// 营运能力
	APTD string `json:"APTD"`
	APTR string `json:"APTR"`
	ARTD string `json:"ARTD"`
	ARTR string `json:"ARTR"`
	CATR string `json:"CATR"`
	FATR string `json:"FATR"`
	ITD  string `json:"ITD"`
	ITR  string `json:"ITR"`
	OC   string `json:"OC"`
	SETR string `json:"SETR"`
	TATR string `json:"TATR"`

	// 现金状况
	CACENI       string `json:"CACENI"`
	CSDAA        string `json:"CSDAA"`
	FCFl         string `json:"FCFl"`
	NBAGCF       string `json:"NBAGCF"`
	NBAGCFDNE    string `json:"NBAGCFDNE"`
	NBAGCFDNETTM string `json:"NBAGCFDNETTM"`
	NBAGCFDR     string `json:"NBAGCFDR"`
	NBAGCFDRTTM  string `json:"NBAGCFDRTTM"`
	NPCL         string `json:"NPCL"`
	SGPCRS       string `json:"SGPCRS"`
	SGPCRSDR     string `json:"SGPCRSDR"`
	SGPCRSDRTTM  string `json:"SGPCRSDRTTM"`
	TACRR        string `json:"TACRR"`

	// 分红能力
	CCEB string `json:"CCEB"`
	CDPM string `json:"CDPM"`
	DPM  string `json:"DPM"`
	DPR  string `json:"DPR"`
	DPS  string `json:"DPS"`
	RER  string `json:"RER"`

	// 资本结构
	ALR     string `json:"ALR"`
	BPCSDIC string `json:"BPCSDIC"`
	BPDTA   string `json:"BPDTA"`
	CADTA   string `json:"CADTA"`
	CLDTL   string `json:"CLDTL"`
	EM      string `json:"EM"`
	FAR     string `json:"FAR"`
	IAR     string `json:"IAR"`
	IBDDIC  string `json:"IBDDIC"`
	LTASR   string `json:"LTASR"`
	LTBDTA  string `json:"LTBDTA"`
	LTLDSET string `json:"LTLDSET"`
	NCADTA  string `json:"NCADTA"`
	NCLDTL  string `json:"NCLDTL"`
	SHER    string `json:"SHER"`
	WC      string `json:"WC"`

	// 收益质量
	IDP         string `json:"IDP"`
	IIJVCDTP    string `json:"IIJVCDTP"`
	IIJVCDTPTTM string `json:"IIJVCDTPTTM"`
	ONIDTP      string `json:"NIDTP"`
	ONIDTPTTM   string `json:"NIDTPTTM"`
	NIVCDTP     string `json:"NIVCDTP"`
	NIVCDTPTTM  string `json:"NIVCDTPTTM"`
	NNOIDTP     string `json:"NNOIDTP"`
	NNOIDTPTTM  string `json:"NNOIDTPTTM"`
	NPADNRGALNP string `json:"NPADNRGALNP"`

	// 杜邦分析
	BPCNPSNP string `json:"BPCNPSNP"`
	EMDA     string `json:"EMDA"`
	NIDTP    string `json:"NIDTP"`
	NPDBR    string `json:"NPDBR"`
}

type ResponseInfo struct {
	SCode string      `json:"scode"`
	MU    string      `json:"MU"`
	AS    string      `json:"AS"`
	List  interface{} `json:"List"`
}

func NewIndicators() *Indicators {
	return &Indicators{}
}
func (this *Indicators) GetJson(c *gin.Context) (*ResponseInfo, error) {
	node1 := Indicators{
		Date: time.Now().Unix(),
	}
	list := []Indicators{}
	list = append(list, node1)
	res := &ResponseInfo{}
	res.List = list

	return res, nil
}
