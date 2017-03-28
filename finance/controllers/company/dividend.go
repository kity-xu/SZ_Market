package company

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models/company"
	"haina.com/share/lib"
)

type List interface{}

type Data struct {
	ComCode string      `json:"scode"`
	Total   float64     `json:"tocash"`
	List    interface{} `json:"list"`
}

type DividendInfo struct {
}

func NewDividendInfo() *DividendInfo {
	return &DividendInfo{}
}

func (this *DividendInfo) GetDiv(c *gin.Context) {
	scode := c.Query("scode")
	sets, e := strconv.Atoi(c.Query("sets"))
	if e != nil {
		lib.WriteString(c, 40004, "invalid sets..")
		return
	}
	fin := new(company.FinDivid)
	divs, err := fin.GetDivListJson(uint64(sets), strings.Split(scode, ".")[0])
	if err != nil {
		lib.WriteString(c, 300, err.Error())
		return
	}
	var data Data
	data.Total = fin.ToCash
	data.List = divs
	data.ComCode = scode
	lib.WriteString(c, 200, data)
}

func (this *DividendInfo) GetSEO(c *gin.Context) {
	scode := c.Query("scode")
	fin := new(company.FinDivid)
	seos, err := fin.GetSEOListJson(strings.Split(scode, ".")[0])
	if err != nil {
		lib.WriteString(c, 300, err.Error())
		return
	}
	var data Data
	data.Total = fin.SeoRaise
	data.List = seos
	data.ComCode = scode
	lib.WriteString(c, 200, data)
}
func (this *DividendInfo) GetRO(c *gin.Context) {
	scode := c.Query("scode")
	fin := new(company.FinDivid)
	ros, err := fin.GetROListJson(strings.Split(scode, ".")[0])
	if err != nil {
		lib.WriteString(c, 300, err.Error())
	}
	var data Data
	data.Total = fin.RoRaise
	data.List = ros
	data.ComCode = scode
	lib.WriteString(c, 200, data)
}
