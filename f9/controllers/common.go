package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"haina.com/share/logging"
)

const (
	CONST_SECURITY_ID = "sid" //证券id
)

//规范sid参数
func Param_Norm_Sid(c *gin.Context) (string, error) {
	id := c.Query(CONST_SECURITY_ID)
	_, err := strconv.Atoi(id)
	if err != nil {
		logging.Error("Invalid param(sid) |%v", err)
	}
	return id, err
}
