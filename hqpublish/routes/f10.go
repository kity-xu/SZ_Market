package routes

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/hqpublish/controllers/publish/f10"
)

func RegF10(rg *gin.RouterGroup) {
	// F10首页
	rg.POST("/f10/home", f10.NewHN_F10_Mobile().GetF10_Mobile)

	rg.POST("/f10/comInfo", f10.NewCompany().GetF10_ComInfo)
}
