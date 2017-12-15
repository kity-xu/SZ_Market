package f10

import (
	"encoding/json"
	"fmt"

	"haina.com/share/lib"

	"github.com/gin-gonic/gin"
	. "haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish/f10"

	"haina.com/share/logging"
	"strconv"
)

type HN_F10_Mobile struct {
}

func NewHN_F10_Mobile() *HN_F10_Mobile {
	return &HN_F10_Mobile{}
}

// F10 首页
func (this *HN_F10_Mobile) GetF10_Mobile(c *gin.Context) {

	var _param struct {
		Scode int `json:"sid" binding:"required"`
	}

	if err := c.BindJSON(&_param); err != nil {
		logging.Debug("Bind Json | %v", err)
		lib.WriteString(c, 40004, nil)
		return
	}

	scode := strconv.Itoa(_param.Scode)
	// 查询redis
	red_data, _ := RedisCache.Get(fmt.Sprintf(REDIS_F10_HOMEPAGE, scode))
	if len(red_data) > 0 { // 如果redis有数据取redis数据
		var fdate f10.F10MobileTerminal
		e := json.Unmarshal([]byte(red_data), &fdate)
		if e != nil {
			logging.Error("Json Unmarshal Error | %v", e)
		}

		lib.WriteString(c, 200, fdate)
		return
	}
	f10, err := f10.F10Mobile(scode)
	if err != nil {
		logging.Error("%v", err)
		lib.WriteString(c, 40002, nil)
		return
	}

	// 存储redis
	byte, err := json.Marshal(f10)
	errr := RedisCache.Set(fmt.Sprintf(REDIS_F10_HOMEPAGE, _param.Scode), byte)
	if errr != nil {
		logging.Error("Redis Set HomePage Error | %v", errr)
	}

	// 设置过期时间
	key := fmt.Sprintf(REDIS_F10_HOMEPAGE, _param.Scode)
	RedisCache.Do("EXPIRE", key, TTL.F10HomePage)

	lib.WriteString(c, 200, f10)
}
