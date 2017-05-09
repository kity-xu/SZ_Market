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

func (this *DividendInfo) GetDiv(c *gin.Context) {
	scode := c.Query(models.CONTEXT_SCODE)
	market := strings.Split(scode, ".")
	if len(market) < 2 {
		return
	}

	var count uint64
	sets := c.Query(models.CONTEXT_COUNT)
	if sets != "" {
		n, err := strconv.Atoi(sets)
		count = uint64(n)
		if err != nil {
			lib.WriteString(c, 40004, err.Error())
			return
		}
	} else {
		count = models.CONTEXT_RETNUM
	}

	fin := new(company.Dividend)
	div, err := fin.GetDividendList(count, market[0], market[1])
	if err != nil {
		lib.WriteString(c, 40002, err.Error())
		return
	}
	div.Scode = scode
	lib.WriteString(c, 200, div)
}

func (this *DividendInfo) GetSEO(c *gin.Context) {
	scode := c.Query(models.CONTEXT_SCODE)
	market := strings.Split(scode, ".")
	if len(market) < 2 {
		return
	}

	fin := new(company.SEO)
	seo, err := fin.GetSEOList(market[0], market[1])
	if err != nil {
		lib.WriteString(c, 40002, err.Error())
		return
	}
	seo.Scode = scode

	lib.WriteString(c, 200, seo)
}
func (this *DividendInfo) GetRO(c *gin.Context) {
	scode := c.Query(models.CONTEXT_SCODE)
	market := strings.Split(scode, ".")
	if len(market) < 2 {
		return
	}
	fin := new(company.RO)
	ro, err := fin.GetROList(market[0], market[1])
	if err != nil {
		lib.WriteString(c, 40002, err.Error())
	}
	ro.Scode = scode

	lib.WriteString(c, 200, ro)
}
