package company

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models"
	"haina.com/market/finance/models/company"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type DividendInfo struct {
}

func NewDividendInfo() *DividendInfo {
	return &DividendInfo{}
}

func (this *DividendInfo) GetDiv(c *gin.Context) {
	var count uint64
	scode := c.Query(models.CONTEXT_SCODE)
	sets := c.Query(models.CONTEXT_COUNT)

	scodePrefix, market, err := ParseSCode(scode)
	if err != nil {
		logging.Error("%v", err)
		lib.WriteString(c, 40004, nil)
		return
	}

	if sets != "" {
		n, err := strconv.Atoi(sets)
		if err != nil {
			logging.Error("%v", err)
			lib.WriteString(c, 40004, nil)
			return
		}
		count = uint64(n)
	} else {
		count = models.CONTEXT_RETNUM
	}

	div, err := new(company.Dividend).GetDividendList(count, scodePrefix, market)
	if err != nil {
		logging.Error("%v", err)
		lib.WriteString(c, 40002, nil)
		return
	}
	div.Scode = scode
	lib.WriteString(c, 200, div)
}

func (this *DividendInfo) GetSEO(c *gin.Context) {
	scode := c.Query(models.CONTEXT_SCODE)
	scodePrefix, market, err := ParseSCode(scode)
	if err != nil {
		logging.Error("%v", err)
		lib.WriteString(c, 40004, nil)
		return
	}

	seo, err := new(company.SEO).GetSEOList(scodePrefix, market)
	if err != nil {
		logging.Error("%v", err)
		lib.WriteString(c, 40002, nil)
		return
	}

	seo.Scode = scode
	lib.WriteString(c, 200, seo)
}
func (this *DividendInfo) GetRO(c *gin.Context) {
	scode := c.Query(models.CONTEXT_SCODE)
	scodePrefix, market, err := ParseSCode(scode)
	if err != nil {
		logging.Error("%v", err)
		lib.WriteString(c, 40004, nil)
		return
	}

	ro, err := new(company.RO).GetROList(scodePrefix, market)
	if err != nil {
		logging.Error("%v", err)
		lib.WriteString(c, 40002, nil)
	}
	ro.Scode = scode

	lib.WriteString(c, 200, ro)
}
