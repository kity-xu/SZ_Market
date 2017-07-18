//历史K线数据（年、月、周、日(PB)
package redistore

import (
	"fmt"
	"time"

	. "haina.com/share/models"

	"ProtocolBuffer/projects/hqpost/go/protocol"

	"haina.com/market/hqpost/models/filestore"

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
func (this *HKLine) LPushHisKLine(sid int32, line *protocol.KInfo) error {
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

func (this *HKLine) UpdateWeekKLineToRedis(sid int32, today *protocol.KInfo) error {
	key := fmt.Sprintf(this.CacheKey, sid)

	n, e := redis.Llen(key) //新股上市
	if n == 0 && e == nil {
		data, err := proto.Marshal(today)
		if err != nil {
			logging.Error("Marshal error...", err.Error())
			return err
		}

		if err = redis.Lpush(key, data); err != nil {
			logging.Error("Lpush error...", err.Error())
			return err
		}
		return nil
	}

	ss, err := redis.LRange(key, 0, 0)
	if err != nil {
		return err
	}
	if len(ss) == 0 {
		return models.ERROR_REDIS_LIST_NULL
	}

	var kinfo protocol.KInfo
	if err := proto.Unmarshal([]byte(ss[0]), &kinfo); err != nil {
		return err
	}

	b1, _ := filestore.DateAdd(kinfo.NTime) //找到该日期所在周日的那天
	b2, _ := filestore.DateAdd(today.NTime)

	if b1.Equal(b2) { //同属一周
		result := compareKInfo(&kinfo, today)

		data, err := proto.Marshal(result)
		if err != nil {
			logging.Error("Marshal error...", err.Error())
			return err
		}

		if err = redis.LSet(key, 0, data); err != nil {
			logging.Error("Lset error...", err.Error())
			return err
		}
		return nil
	} else { //不属于同周
		result := *today
		result.NPreCPx = kinfo.NLastPx //昨收价

		data, err := proto.Marshal(&result)
		if err != nil {
			logging.Error("Marshal error...", err.Error())
			return err
		}

		if err = redis.Lpush(key, data); err != nil {
			logging.Error("Lpush error...", err.Error())
			return err
		}
		return nil
	}
}

func (this *HKLine) UpdateMonthKLineToRedis(sid int32, today *protocol.KInfo) error {
	key := fmt.Sprintf(this.CacheKey, sid)

	n, e := redis.Llen(key) //新股上市
	if n == 0 && e == nil {
		data, err := proto.Marshal(today)
		if err != nil {
			logging.Error("Marshal error...", err.Error())
			return err
		}

		if err = redis.Lpush(key, data); err != nil {
			logging.Error("Lpush error...", err.Error())
			return err
		}
		return nil
	}

	ss, err := redis.LRange(key, 0, 0)
	if err != nil {
		return err
	}
	if len(ss) == 0 {
		return models.ERROR_REDIS_LIST_NULL
	}

	var kinfo protocol.KInfo
	if err := proto.Unmarshal([]byte(ss[0]), &kinfo); err != nil {
		return err
	}

	if kinfo.NTime/100 == today.NTime/100 { //同属一月
		result := compareKInfo(&kinfo, today)

		data, err := proto.Marshal(result)
		if err != nil {
			logging.Error("Marshal error...", err.Error())
			return err
		}

		if err = redis.LSet(key, 0, data); err != nil {
			logging.Error("Lset error...", err.Error())
			return err
		}
		return nil
	} else { //不属于同月
		result := *today
		result.NPreCPx = kinfo.NLastPx //昨收价

		data, err := proto.Marshal(&result)
		if err != nil {
			logging.Error("Marshal error...", err.Error())
			return err
		}

		if err = redis.Lpush(key, data); err != nil {
			logging.Error("Lpush error...", err.Error())
			return err
		}
		return nil
	}
}

func (this *HKLine) UpdateYearKLineToRedis(sid int32, today *protocol.KInfo) error {
	key := fmt.Sprintf(this.CacheKey, sid)

	n, e := redis.Llen(key) //新股上市
	if n == 0 && e == nil {
		data, err := proto.Marshal(today)
		if err != nil {
			logging.Error("Marshal error...", err.Error())
			return err
		}

		if err = redis.Lpush(key, data); err != nil {
			logging.Error("Lpush error...", err.Error())
			return err
		}
		return nil
	}
	proto.Unmarshal(data)
	proto.Marshal()

	ss, err := redis.LRange(key, 0, 0)
	if err != nil {
		return err
	}
	if len(ss) == 0 {
		return models.ERROR_REDIS_LIST_NULL
	}

	var kinfo protocol.KInfo
	if err := proto.Unmarshal([]byte(ss[0]), &kinfo); err != nil {
		return err
	}

	if kinfo.NTime/10000 == today.NTime/10000 { //同属一年
		result := compareKInfo(&kinfo, today)

		data, err := proto.Marshal(result)
		if err != nil {
			logging.Error("Marshal error...", err.Error())
			return err
		}

		if err = redis.LSet(key, 0, data); err != nil {
			logging.Error("Lset error...", err.Error())
			return err
		}
		return nil
	} else { //不属于同年
		result := *today
		result.NPreCPx = kinfo.NLastPx //昨收价

		data, err := proto.Marshal(&result)
		if err != nil {
			logging.Error("Marshal error...", err.Error())
			return err
		}

		if err = redis.Lpush(key, data); err != nil {
			logging.Error("Lpush error...", err.Error())
			return err
		}
		return nil
	}
}

//Append today kline
func (this *HKLine) AppendTodayLine(sid int32, latest *protocol.KInfo) error {
	key := fmt.Sprintf(this.CacheKey, sid)

	data, err := proto.Marshal(latest)
	if err != nil {
		logging.Error("%v", err.Error())
		return err
	}
	err = redis.Lpush(key, data)
	return err
}

//比较历史最新和当天
func compareKInfo(tmp *protocol.KInfo, today *protocol.KInfo) *protocol.KInfo {
	var swap protocol.KInfo

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

func (this *HKLine) HQpostExecutedTime() {
	timestamp := time.Now().Unix()

	tm := time.Unix(timestamp, 0)
	ss := tm.Format("200601021504")
	if err := redis.Set(this.CacheKey, []byte(ss)); err != nil {
		logging.Error("%v", err.Error())
	}
}
