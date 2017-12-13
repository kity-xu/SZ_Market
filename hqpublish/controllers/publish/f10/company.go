package f10

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	. "haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish/f10"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type Company struct {
}

func NewCompany() *Company {
	return &Company{}
}

type Share struct {
	Scode     string      `json:"sid"`
	ComDetail interface{} `json:"comDetail"`
	Leader    interface{} `json:"leader"`
}

// 获取公司详细信息
func (this *Company) GetF10_ComInfo(c *gin.Context) {

	var _param struct {
		Scode string `json:"scode" binding:"required"`
	}

	if err := c.BindJSON(&_param); err != nil {
		logging.Debug("Bind Json | %v", err)
		lib.WriteString(c, 40004, nil)
		return
	}

	// 查询redis

	red_data, _ := RedisCache.Get(fmt.Sprintf(REDIS_F10_COMINFO, _param.Scode))
	if len(red_data) > 0 { // 如果redis有数据取redis数据
		var fdate f10.Compinfo
		e := json.Unmarshal([]byte(red_data), &fdate)
		if e != nil {
			logging.Error("Json Unmarshal Error | %v", e)
		}

		lib.WriteString(c, 200, fdate)
		return
	}

	// 公司信息
	date, err := f10.GetF10Company(_param.Scode)
	if err != nil {
		logging.Error("%v", err)
	}
	// 存储redis
	byte, err := json.Marshal(date)
	errr := RedisCache.Set(fmt.Sprintf(REDIS_F10_COMINFO, _param.Scode), byte)
	if errr != nil {
		logging.Error("Redis Set PlayBack Error | %v", errr)
	}

	// 设置过期时间
	RedisCache.Do("EXPIRE", REDIS_F10_COMINFO, TTL.F10HomePage)

	lib.WriteString(c, 200, date)
}
