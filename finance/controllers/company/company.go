package company

import (
	"strings"

	"github.com/gin-gonic/gin"
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
	scode := c.Query("scode")

	cominfo, err := new(company.CompInfo).GetCompInfoFromFC(strings.Split(scode, ".")[0])
	if err != nil {
		lib.WriteString(c, 300, err.Error())
	}
	var data Share
	data.Scode = scode
	data.List = cominfo

	lib.WriteString(c, 200, data)
}

// 获取高管信息
func (this *Company) GetManagreInfo(c *gin.Context) {
	secode := c.Query("scode")

	list, err := new(company.HnManager).GetManagerInfo(strings.Split(secode, ".")[0])
	if err != nil {
		lib.WriteString(c, 300, err.Error())
		return
	}
	var data Share
	data.Scode = secode
	data.List = list

	lib.WriteString(c, 200, data)
}
