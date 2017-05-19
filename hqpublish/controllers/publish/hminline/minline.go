package hminline

import (
	"github.com/gin-gonic/gin"
)

type KMinline struct {
}

func NewKMinline() *KMinline {
	return &KMinline{}
}

func (this *KMinline) POST(c *gin.Context) {
	TestMinLine()
}
