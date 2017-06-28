package redistore

import (
	"errors"
	"fmt"
	"strconv"

	. "haina.com/share/models"
	"haina.com/share/store/redis"
)

type GlobalSid struct {
	Model `db:"-"`
}

func NewGlobalSid(key string) *GlobalSid {
	return &GlobalSid{
		Model: Model{
			CacheKey: key,
		},
	}
}

func (this *GlobalSid) GetGlobalSidFromRedis() (*[]int32, error) {
	sids, err := redis.LRange(this.CacheKey, 0, -1)
	if err != nil {
		return nil, err
	}
	if len(sids) < 1 {
		return nil, errors.New("sids list is null...")
	}

	var NSids []int32
	for _, sid := range sids {
		nsid, err := strconv.Atoi(sid)
		if err != nil { //此处出错误是因为出现了非数字字符
			errNew := fmt.Sprintf("The sid is not numeric types...%s", sid)
			return nil, errors.New(errNew)
		}
		NSids = append(NSids, int32(nsid))
	}
	return &NSids, nil
}
