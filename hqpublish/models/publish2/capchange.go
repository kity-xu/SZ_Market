//资金统计
package publish2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	. "haina.com/market/hqpublish/models"
	"haina.com/share/logging"
)

// redis 结构
type MarketsStatistics struct {
	MarketID     int32 ///<市场ID
	RiseLimitNum int32 ///<涨停家数
	Rise8_Num    int32 ///<涨幅8-涨停家数
	Rise6_8Num   int32 ///<涨幅6-8家数
	Rise4_6Num   int32 ///<涨幅4-6家数
	Rise2_4Num   int32 ///<涨幅2-4家数
	Rise0_2Num   int32 ///<涨幅0-2家数
	Fall2_0Num   int32 ///<下跌-2-0家数
	Fall4_2Num   int32 ///<下跌-4--2家数
	Fall6_4Num   int32 ///<下跌-6--4家数
	Fall8_6Num   int32 ///<下跌-8--6家数
	Fall_8Num    int32 ///<下跌-8-跌停家数
	FallLimitNum int32 ///<跌停家数
}

type StatisticsJson struct {
	MarketID     int32 `json:"marketId"`
	FallLimitNum int32 `json:"fallStop"`
	Fall_8Num    int32 `json:"fall8"`
	Fall8_6Num   int32 `json:"fall6"`
	Fall6_4Num   int32 `json:"fall4"`
	Fall4_2Num   int32 `json:"fall2"`
	Fall2_0Num   int32 `json:"fall0"`
	Rise0_2Num   int32 `json:"rise2"`
	Rise2_4Num   int32 `json:"rise4"`
	Rise4_6Num   int32 `json:"rise6"`
	Rise6_8Num   int32 `json:"rise8"`
	Rise8_Num    int32 `json:"riseStop"`
	RiseLimitNum int32 `json:"riseTop"`
}

func NewMarketsStatistics() *MarketsStatistics {
	return &MarketsStatistics{}
}

func (MarketsStatistics) GetMarketsStatistics() (*StatisticsJson, error) {
	key := fmt.Sprintf("hq:market:statistics:%d", 2) // 2 A股
	da, err := RedisStore.Get(key)
	if len(da) == 0 || err != nil {
		logging.Error("%s", err.Error())
		return nil, err
	}

	logging.Info("len :%d", len(da))

	data := &MarketsStatistics{}

	logging.Info("len data:%d", binary.Size(data))

	if err = binary.Read(bytes.NewBuffer([]byte(da)), binary.LittleEndian, data); err != nil && err != io.EOF {
		logging.Error("%s", err.Error())
		return nil, err
	}

	result := &StatisticsJson{
		MarketID:     data.MarketID,
		FallLimitNum: data.FallLimitNum,
		Fall_8Num:    data.Fall_8Num,
		Fall8_6Num:   data.Fall8_6Num,
		Fall6_4Num:   data.Fall6_4Num,
		Fall4_2Num:   data.Fall4_2Num,
		Fall2_0Num:   data.Fall2_0Num,
		Rise0_2Num:   data.Rise0_2Num,
		Rise2_4Num:   data.Rise2_4Num,
		Rise4_6Num:   data.Rise4_6Num,
		Rise6_8Num:   data.Rise6_8Num,
		Rise8_Num:    data.Rise8_Num,
		RiseLimitNum: data.RiseLimitNum,
	}
	return result, nil
}
