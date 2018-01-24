package publish

import (
	"github.com/gin-gonic/gin"
	"haina.com/share/lib"
)

type TestPage struct{}

func NewTest() *TestPage {
	return &TestPage{}
}

func (this *TestPage) Test(c *gin.Context) {
	lib.WriteString(c, 200, "请求成功！")
}
