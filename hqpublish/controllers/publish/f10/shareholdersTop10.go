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

type ShareholderslTop10 struct {
}

func NewShareholderslTop10() *ShareholderslTop10 {
	return &ShareholderslTop10{}
}

// 历史股本变动
func (this *ShareholderslTop10) GetShareholdersTop10(c *gin.Context) {

	var _param struct {
		Scode   int    `json:"sid" binding:"required"`
		Count   int32  `json:"count"`
		HType   int32  `json:"htype" binding:"required"`
		EndDate string `json:"enddate"`
	}

	if err := c.BindJSON(&_param); err != nil {
		logging.Debug("Bind Json | %v", err)
		lib.WriteString(c, 40004, nil)
		return
	}
	scode := strconv.Itoa(_param.Scode)

	// 查询redis
	//if _param.HType == 1 {
	red_data, _ := RedisCache.Get(fmt.Sprintf(REDIS_F10_SHAREHOLDERSTOP10, _param.HType, _param.Scode))
	if len(red_data) > 0 { // 如果redis有数据取redis数据
		var fdate f10.Date
		e := json.Unmarshal([]byte(red_data), &fdate)
		if e != nil {
			logging.Error("Json Unmarshal Error | %v", e)
		}

		lib.WriteString(c, 200, fdate)
		return
	}
	//}

	date, err := f10.GetHN_F10_ShareholdersTop10(scode, _param.Count, _param.HType, _param.EndDate)
	if err != nil {
		logging.Error("%v", err)
		lib.WriteString(c, 40002, nil)
		return
	}

	// 存储redis
	byte, err := json.Marshal(date)
	errr := RedisCache.Set(fmt.Sprintf(REDIS_F10_SHAREHOLDERSTOP10, _param.HType, _param.Scode), byte)
	if errr != nil {
		logging.Error("Redis Set Shareholderstop10 Error | %v", errr)
	}

	// 设置过期时间
	RedisCache.Do("EXPIRE", REDIS_F10_SHAREHOLDERSTOP10, TTL.F10HomePage)

	lib.WriteString(c, 200, date)
}
