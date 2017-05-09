// 利润表
package company

import (
	"time"

	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models"
	"haina.com/market/finance/models/company"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type Profits struct {
}

func NewProfits() *Profits {
	return &Profits{}
}

func (this *Profits) GET(c *gin.Context) {
	scode := c.Query(models.CONTEXT_SCODE)
	stype := c.Query(models.CONTEXT_TYPE)
	spage := c.Query(models.CONTEXT_PAGE)
	sperp := c.Query(models.CONTEXT_PERPAGE)

	req := CheckAndNewRequestParam(scode, stype, sperp, spage)
	if req == nil {
		lib.WriteString(c, 40004, nil)
		return
	}

	data, err := this.getJson(req)
	if err != nil {
		logging.Debug("%v", err)
		lib.WriteString(c, 40002, nil)
		return
	}
	data.SCode = scode

	lib.WriteString(c, 200, data)
}

func (this *Profits) getJson(req *RequestParam) (*company.RespFinAnaJson, error) {
	sli := make([]company.ProfitsJson, 0, req.PerPage)
	ls, err := company.NewProfits().GetList(req.SCode, req.Market, req.Type, req.PerPage, req.Page)
	if err != nil {
		return nil, err
	}
	for _, v := range ls {
		one := company.ProfitsJson{
			AILs: v.ASSEIMPALOSS.Float64,
			AREp: v.REINEXPE.Float64,
			BPAC: v.PARENETP.Float64,
			CoEp: v.POUNEXPE.Float64,
			CoRe: v.POUNINCO.Float64,
			CORe: v.BIZCOST.Float64,
			DPES: v.DILUTEDEPS.Float64,
			EPS:  v.BASICEPS.Float64,
			FnEp: v.FINEXPE.Float64,
			IDEp: v.POLIDIVIEXPE.Float64,
			InRe: v.INTEINCO.Float64,
			ItEp: v.INTEEXPE.Float64,
			ITEp: v.INCOTAXEXPE.Float64,
			MgEp: v.MANAEXPE.Float64,
			MIIn: v.MINYSHARRIGH.Float64,
			NOEp: v.NONOEXPE.Float64,
			NORe: v.NONOREVE.Float64,
			NtIn: v.NETPROFIT.Float64,
			OATx: v.BIZTAX.Float64,
			OCOR: v.BIZTOTCOST.Float64,
			OpPr: v.PERPROFIT.Float64,
			OpRe: v.BIZINCO.Float64,
			SaEp: v.SALESEXPE.Float64,
			TOpR: v.BIZTOTINCO.Float64,
			ToPr: v.TOTPROFIT.Float64,
		}

		if v.ENDDATE.Valid {
			tm, err := time.Parse("20060102", v.ENDDATE.String)
			if err != nil {
				return nil, err
			}
			one.Date = tm.Unix()
		}

		sli = append(sli, one)
	}

	// A股
	jsn := &company.RespFinAnaJson{
		MU:     "人民币元",
		AS:     "新会计准则",
		Length: len(sli),
		List:   sli,
	}

	return jsn, nil
}
