//周K线
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

type WeekLine struct {
	Model `db:"-"`
}

func NewWeekLine() *WeekLine {
	return &WeekLine{
		Model: Model{
			CacheKey: REDISKEY_SECURITY_HWEEK,
		},
	}
}

// 获取周线
func (this *WeekLine) GetWeekLine(sid int32, betime int32) (*[]*kline.KInfo, int, error) {
	var wlines kline.KInfoTable
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

	if err := proto.Unmarshal(data, &wlines); err != nil {
		logging.Error("%v", err.Error())
		return nil, 0, err
	}

	for _, v := range wlines.List {
		if v.NTime >= betime {
			array = append(array, v)
		}
	}
	if len(array) < 1 {
		return &array, len(wlines.List), errors.New("To get the data is empty, may start time is wrong...")
	}

	return &array, len(wlines.List), nil
}
