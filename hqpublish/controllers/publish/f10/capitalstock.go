package f10

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	. "haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish/f10"
	"haina.com/share/lib"
	"haina.com/share/logging"
	"strconv"
)

type CapitalStock struct {
}

func NewCapitalStock() *CapitalStock {
	return &CapitalStock{}
}

// 历史股本变动
func (this *CapitalStock) GetF10_CapitalStock(c *gin.Context) {

	var _param struct {
		Scode int `json:"sid" binding:"required"`
		Count int `json:"count"`
	}

	if err := c.BindJSON(&_param); err != nil {
		logging.Debug("Bind Json | %v", err)
		lib.WriteString(c, 40004, nil)
		return
	}
	scode := strconv.Itoa(_param.Scode)
	// 查询redis
	red_data, _ := RedisCache.Get(fmt.Sprintf(REDIS_F10_CAPITALSTOCK, _param.Scode))
	if len(red_data) > 0 { // 如果redis有数据取redis数据
		var fdate f10.CSDate
		e := json.Unmarshal([]byte(red_data), &fdate)
		if e != nil {
			logging.Error("Json Unmarshal Error | %v", e)
		}

		lib.WriteString(c, 200, fdate)
		return
	}

	limit := _param.Count
	if limit != 10 {
		limit = 10
	}
	csdate, err := f10.GetF10CapitalStock(scode, limit)
	if err != nil {
		logging.Error("%v", err)
		lib.WriteString(c, 40002, nil)
		return
	}

	// 存储redis
	byte, err := json.Marshal(csdate)
	errr := RedisCache.Set(fmt.Sprintf(REDIS_F10_CAPITALSTOCK, _param.Scode), byte)
	if errr != nil {
		logging.Error("Redis Set F10CapitalStock Error | %v", errr)
	}

	// 设置过期时间
	key := fmt.Sprintf(REDIS_F10_CAPITALSTOCK, _param.Scode)
	RedisCache.Do("EXPIRE", key, TTL.F10HomePage)

	lib.WriteString(c, 200, csdate)
}
