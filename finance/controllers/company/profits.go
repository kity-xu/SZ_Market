// 利润表
package company

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models/company"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type Profits struct {
}

func NewProfits() *Profits {
	return &Profits{}
}

func (this *Profits) getJson(c *gin.Context) (*company.ResponseInfo, error) {
	return company.NewProfits().GetJson(c)
}

func (this *Profits) GET(c *gin.Context) {
	data, err := this.getJson(c)
	if err != nil {
		logging.Debug("%v", err)
		lib.WriteString(c, 40002, nil)
		return
	}

	lib.WriteString(c, 200, data)
}
