package company

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models"
	"haina.com/market/finance/models/company"
	"haina.com/share/lib"
)

type DividendInfo struct {
}

func NewDividendInfo() *DividendInfo {
	return &DividendInfo{}
}

type divData struct {
	Scode  string      `json:"scode"`  //证券内码
	Tocash float64     `json:"tocash"` //累计分红金额
	Count  int         `json:"count"`  //共分红次数
	List   interface{} `json:"list"`   //分红数据
}

type Data struct {
	Scode   string      `json:"scode"`
	Tocash  float64     `json:"tocash"`
	Count   int         `json:"count"`   //总数
	Success int         `json:"success"` //成功数
	Faild   int         `json:"faild"`   //失败数
	Ing     int         `json:"ing"`     //进行中
	List    interface{} `json:"list"`
}

type roData struct {
}

func (this *DividendInfo) GetDiv(c *gin.Context) {
	scode := c.Query(models.CONTEXT_SCODE)
	sets, e := strconv.Atoi(c.Query("sets"))
	if e != nil {
		lib.WriteString(c, 40004, e.Error())
		return
	}
	fin := new(company.Dividend)
	divs, err := fin.GetDividendList(uint64(sets), strings.Split(scode, ".")[0])
	if err != nil {
		lib.WriteString(c, 40002, err.Error())
		return
	}
	var data divData
	data.Scode = scode
	data.Tocash = fin.ToCash
	data.Count = fin.Divcount
	data.List = divs

	lib.WriteString(c, 200, data)
}

func (this *DividendInfo) GetSEO(c *gin.Context) {
	scode := c.Query(models.CONTEXT_SCODE)
	fin := new(company.SEO)
	seos, err := fin.GetSEOList(strings.Split(scode, ".")[0])
	if err != nil {
		lib.WriteString(c, 40002, err.Error())
		return
	}
	var data Data
	data.Scode = scode
	data.Tocash = fin.SeoRaise
	data.Count = fin.SEOcount
	data.Success = int(fin.SeoSc)
	data.Ing = fin.SEOcount - int(fin.SeoSc)
	data.List = seos

	lib.WriteString(c, 200, data)
}
func (this *DividendInfo) GetRO(c *gin.Context) {
	scode := c.Query(models.CONTEXT_SCODE)
	fin := new(company.RO)
	ros, err := fin.GetROList(strings.Split(scode, ".")[0])
	if err != nil {
		lib.WriteString(c, 40002, err.Error())
	}
	var data Data
	data.Scode = scode
	data.Tocash = fin.RoRaise
	data.Count = fin.ROcount
	data.Success = int(fin.RoSc)
	data.Ing = fin.ROcount - int(fin.RoSc)
	data.List = ros

	lib.WriteString(c, 200, data)
}
