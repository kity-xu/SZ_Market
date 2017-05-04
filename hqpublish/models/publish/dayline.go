//日K线
package publish

import (
	"errors"
	"fmt"

	. "haina.com/share/models"

	"ProtocolBuffer/format/kline"

	"github.com/golang/protobuf/proto"

	redigo "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

var _ = fmt.Println
var _ = redigo.Bytes
var _ = logging.Info

type DayLine struct {
	Model `db:"-"`
}

func NewDayLine() *DayLine {
	return &DayLine{
		Model: Model{
			CacheKey: REDISKEY_SECURITY_HDAY,
		},
	}
}

// 获取日K线
func (this *DayLine) GetHisDayLine(sid int32, betime int32) (*[]*kline.KInfo, int, error) {
	var dlines kline.KInfoTable
	var array []*kline.KInfo
	key := fmt.Sprintf(this.CacheKey, sid)

	dls, err := redis.Get(key)
	if err != nil {
		return nil, 0, err
	}
	if dls == "" {
		return nil, 0, ERROR_REDIS_LIST_NULL
	}

	data := []byte(dls)

	if err := proto.Unmarshal(data, &dlines); err != nil {
		logging.Error("%v", err.Error())
		return nil, 0, err
	}

	for _, v := range dlines.List {
		if v.NTime >= betime {
			array = append(array, v)
		}
	}
	if len(array) < 1 {
		return &array, len(dlines.List), errors.New("To get the data is empty, may start time is wrong...")
	}

	return &array, len(dlines.List), nil
}
