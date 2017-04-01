package company

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models/company"
	"haina.com/share/lib"
)

type Share struct {
	Scode string `json:"scode"`
	List  interface{}
}

type Administrator struct {
}

func NewAdministrator() *Administrator {
	return &Administrator{}
}

// 获取高管信息
func (this *Administrator) GetAdminInfo(c *gin.Context) {
	secode := c.Query("scode")

	list, err := new(company.AdministratorJson).GetAdminListJson(secode)
	if err != nil {
		lib.WriteString(c, 300, err.Error())
		return
	}
	var data Share
	data.Scode = secode
	data.List = list

	lib.WriteString(c, 200, data)
}

// 高管持股变动信息
func (this *Administrator) GetAdminEquityChange(c *gin.Context) {
	secode := c.Query("scode")
	list, err := new(company.AdministratorJson).GetAdminEquityChangeJson(secode)
	if err != nil {
		lib.WriteString(c, 300, err.Error())
		return
	}
	var data Share
	data.Scode = secode
	data.List = list

	lib.WriteString(c, 200, data)
}
