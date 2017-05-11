package company

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models"
	"haina.com/market/finance/models/company"
	"haina.com/share/lib"
)

type Company struct {
}

func NewCompany() *Company {
	return &Company{}
}

type Share struct {
	Scode string      `json:"scode"`
	List  interface{} `json:"list"`
}

func (this *Company) GetInfo(c *gin.Context) {
	scode := c.Query(models.CONTEXT_SCODE)
	scodePrefix, market, err := ParseSCode(scode)
	if err != nil {
		lib.WriteString(c, 40004, err.Error())
		return
	}

	cominfo, err := new(company.CompInfo).GetCompInfo(scodePrefix, market)
	if err != nil {
		lib.WriteString(c, 40002, err.Error())
		return
	}
	var data Share
	data.Scode = scode
	data.List = cominfo

	lib.WriteString(c, 200, data)
}

// 获取高管信息
func (this *Company) GetManagreInfo(c *gin.Context) {
	scode := c.Query(models.CONTEXT_SCODE)
	scodePrefix, market, err := ParseSCode(scode)
	if err != nil {
		lib.WriteString(c, 40004, err.Error())
		return
	}

	list, err := new(company.HnManager).GetManagerList(scodePrefix, market)
	if err != nil {
		lib.WriteString(c, 40002, err.Error())
		return
	}
	var data Share
	data.Scode = scode
	data.List = list

	lib.WriteString(c, 200, data)
}
