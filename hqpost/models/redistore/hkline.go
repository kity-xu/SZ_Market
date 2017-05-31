//历史K线数据（年、月、周、日(PB)
package redistore

import (
	"errors"
	"fmt"

	. "haina.com/share/models"

	"ProtocolBuffer/format/kline"

	"haina.com/market/hqpost/models"

	"github.com/golang/protobuf/proto"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

type HKLine struct {
	Model `db:"-"`
}

func NewHKLine(key string) *HKLine {
	return &HKLine{
		Model: Model{
			CacheKey: key,
		},
	}
}

//Insert
func (this *HKLine) LPushHisKLine(sid int32, line *kline.KInfo) error {
	key := fmt.Sprintf(this.CacheKey, sid)
	data, err := proto.Marshal(line)
	if err != nil {
		logging.Error("%v", err.Error())
		return err
	}
	if err = redis.Lpush(key, data); err != nil {
		logging.Error("%v", err.Error())
		return err
	}
	return nil
}

//Select
func (this *HKLine) LRangeHisKLine(sid int32, num int, table *[]kline.KInfo) error {
	if num < 1 {
		return errors.New("Invalid request parameters num...")
	}
	key := fmt.Sprintf(this.CacheKey, sid)
	ss, err := redis.LRange(key, 0, num-1)
	if err != nil {
		return err
	}
	if len(ss) == 0 {
		return models.ERROR_REDIS_LIST_NULL
	}

	for _, by := range ss {
		kinfo := kline.KInfo{}
		if err := proto.Unmarshal([]byte(by), &kinfo); err != nil {
			return err
		}
		*table = append(*table, kinfo)
	}
	return nil
}

//Update
func (this *HKLine) LSetHisKLine(sid int32, latest *kline.KInfo) error {
	key := fmt.Sprintf(this.CacheKey, sid)

	data, err := proto.Marshal(latest)
	if err != nil {
		logging.Error("%v", err.Error())
		return err
	}
	err = redis.LSet(key, 0, data)
	return err
}

//比较历史最新和当天
func CompareKInfo(tmp *kline.KInfo, today *kline.KInfo) *kline.KInfo {
	var swap kline.KInfo

	swap.NSID = tmp.NSID
	swap.NTime = tmp.NTime
	swap.NPreCPx = tmp.NPreCPx //昨收价
	swap.NOpenPx = tmp.NOpenPx
	if tmp.NHighPx > today.NHighPx {
		swap.NHighPx = tmp.NHighPx
	} else {
		swap.NHighPx = today.NHighPx
	}
	if tmp.NLowPx > today.NLowPx {
		swap.NLowPx = today.NLowPx
	} else {
		swap.NLowPx = tmp.NLowPx
	}
	swap.NLastPx = today.NLastPx
	swap.LlVolume = today.LlVolume + tmp.LlVolume
	swap.LlValue = today.LlValue + tmp.LlValue
	swap.NAvgPx = (today.NAvgPx + tmp.NAvgPx) / 2
	return &swap
}
