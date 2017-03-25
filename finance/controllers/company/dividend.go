package company

import (
	//"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models/finchina"
	"haina.com/share/lib"
)

type List interface{}

type Data struct {
	Lengh   int         `json:"lengh"`
	ComCode string      `json:"comcode"`
	List    interface{} `json:"list"`
}

type DividendInfo struct {
}

func NewDividendInfo() *DividendInfo {
	return &DividendInfo{}
}

func (this *DividendInfo) GetDiv(c *gin.Context) {
	scode := strings.Split(c.Query("scode"), ".")[0]
	sets, e := strconv.Atoi(c.Query("sets"))
	if e != nil {
		lib.WriteString(c, 400, "invalid sets..")
		return
	}
	divs, err := new(finchina.Dividend).GetDivList(uint64(sets), scode)
	if err != nil {
		lib.WriteString(c, 300, err.Error())
		return
	}
	data := getDivListjson(divs)
	data.ComCode = scode

	lib.WriteString(c, 200, data)
}
func (this *DividendInfo) GetSEO(c *gin.Context) {
	scode := strings.Split(c.Query("scode"), ".")[0]
	seos, err := new(finchina.Dividend).GetSEOList(scode)
	if err != nil {
		lib.WriteString(c, 300, err.Error())
	}
	data := getSEOListjson(seos)
	data.ComCode = scode
	lib.WriteString(c, 200, data)
}
func (this *DividendInfo) GetRO(c *gin.Context) {
	scode := strings.Split(c.Query("scode"), ".")[0]
	ros, err := new(finchina.Dividend).GetROList(scode)
	if err != nil {
		lib.WriteString(c, 300, err.Error())
	}
	data := getROListjson(ros)
	data.ComCode = scode
	lib.WriteString(c, 200, data)
}

func getDivListjson(divs []finchina.Div) Data {
	var i int = 0
	var data Data
	jsn := make([]finchina.DivJson, 0)

	for _, v := range divs {
		var js finchina.DivJson
		js.Bene = v.GRAOBJ.String
		js.Bonus = v.PROBONUSRT.Float64
		js.Date = v.DIVIYEAR.String
		js.Dividend = v.PRETAXCASHMAXDVCNY.Float64
		js.DivDate = v.CASHDVARRBEGDATE.String
		js.INDate = v.PUBLISHDATE.String
		js.DNDate = v.SHHDMEETRESPUBDATE.String
		js.ExDate = v.XDRDATE.String
		js.RegDate = v.EQURECORDDATE.String
		js.Tran = v.TRANADDRT.Float64

		if strings.EqualFold(v.LISTDATE.String, "19000101") {
			js.LisDate = "--"
		} else {
			js.LisDate = v.LISTDATE.String
		}

		jsn = append(jsn, js)
		i++
	}
	data.Lengh = i
	data.List = jsn
	return data
}

func getSEOListjson(seos []finchina.SEO) Data {
	var i int = 0
	var data Data
	jsn := make([]*finchina.SEOJson, 0)

	for _, v := range seos {
		js := &finchina.SEOJson{
			IECD:    v.CSRCAPPRAGREEDATE.String,
			LisDate: v.UPDATEDATE.String,
			//PNDate:  v.ADDISSPUBDATE.String,
			PPrice: v.ENQUMAXPRICE.Float64,
			Price:  v.ISSPRICE.Float64,
			PVal:   v.PLANTOTRAISEAMT.Float64,
			PVol:   v.PLANISSMAXQTY.Float64,
			Range:  v.ISSUEOBJECT.String,
			SEOD:   v.SHARECORDDATE.String,
			SRCD:   v.CSRCAPPDPUBDATE.String,
			Val:    v.NEWTOTRAISEAMT.Float64,
			Vol:    v.ACTISSQTY.Float64,
			Type:   v.ISSUEMODEMEMO.String,
		}
		jsn = append(jsn, js)
		i++
	}
	data.Lengh = i
	data.List = jsn
	return data
}

func getROListjson(ros []finchina.RO) Data {
	var i int = 0
	var data Data
	jsn := make([]finchina.ROJson, 0)

	for _, v := range ros {
		var js finchina.ROJson
		js.AGMD = v.UPDATEDATE.String
		js.Code = v.ALLOTCODE.String
		js.DNDate = v.LISTPUBDATE.String
		js.ERDate = v.EXRIGHTDATE.String
		js.LisDate = v.LISTDATE.String
		js.PNDate = v.PUBLISHDATE.String
		js.PProp = v.ACTTOTALLOTRT.Float64
		js.Price = v.ALLOTPRICE.Float64
		js.Prop = v.ALLOTRT.Float64
		js.PVol = v.PLANISSMAXQTY.Float64
		js.RegDate = v.EQURECORDDATE.String
		js.ROPD = v.PAYBEGDATE.String + "~" + v.PAYENDDATE.String
		js.Short = v.ALLOTSNAME.String
		js.Vol = v.ACTISSQTY.Float64

		jsn = append(jsn, js)
		i++
	}
	data.Lengh = i
	data.List = jsn

	return data
}
