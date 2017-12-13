package f10

import (
	"encoding/json"
	"fmt"

	"haina.im/share/lib"

	"github.com/gin-gonic/gin"
	. "haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish/f10"

	"haina.com/share/logging"
)

type HN_F10_Mobile struct {
}

func NewHN_F10_Mobile() *HN_F10_Mobile {
	return &HN_F10_Mobile{}
}

type F10 struct {
	Scode  string      `json:"scode"`
	Name   *string     `json:"name"`
	Mobile interface{} `json:"f10"`
}

// F10 首页
func (this *HN_F10_Mobile) GetF10_Mobile(c *gin.Context) {

	var _param struct {
		Scode string `json:"scode" binding:"required"`
	}

	if err := c.BindJSON(&_param); err != nil {
		logging.Debug("Bind Json | %v", err)
		lib.WriteString(c, 40004, nil)
		return
	}

	// 查询redis

	red_data, _ := RedisCache.Get(fmt.Sprintf(REDIS_F10_HOMEPAGE, _param.Scode))
	if len(red_data) > 0 { // 如果redis有数据取redis数据
		var fdate F10
		e := json.Unmarshal([]byte(red_data), &fdate)
		if e != nil {
			logging.Error("Json Unmarshal Error | %v", e)
		}

		lib.WriteString(c, 200, fdate)
		return
	}

	f10, name, err := f10.F10Mobile(_param.Scode)
	if err != nil {
		logging.Error("%v", err)
		lib.WriteString(c, 40002, nil)
		return
	}

	result := &F10{
		Scode:  _param.Scode,
		Name:   name,
		Mobile: f10,
	}

	// 存储redis
	byte, err := json.Marshal(result)
	errr := RedisCache.Set(fmt.Sprintf(REDIS_F10_HOMEPAGE, _param.Scode), byte)
	if errr != nil {
		logging.Error("Redis Set PlayBack Error | %v", errr)
	}

	// 设置过期时间
	RedisCache.Do("EXPIRE", REDIS_F10_HOMEPAGE, TTL.F10HomePage)

	lib.WriteString(c, 200, result)
}
