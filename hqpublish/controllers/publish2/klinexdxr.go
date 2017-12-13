package publish2

import "github.com/gin-gonic/gin"

type KlineXDXR struct{}

func NewKlineXDXR() *KlineXDXR {
	return &KlineXDXR{}
}

func (*KlineXDXR) POST(c *gin.Context) {

}
