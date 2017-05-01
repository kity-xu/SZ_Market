package publish

import (
	"fmt"

	. "haina.com/share/models"

	"ProtocolBuffer/format/redis/pbdef/kline"

	//"github.com/golang/protobuf/proto"
	//redigo "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

type MinKLine struct {
	Model `db:"-"`
	reply kline.ReplyMinK
}

func NewMinKLine() *MinKLine {
	return &MinKLine{
		Model: Model{
			CacheKey: REDISKEY_SECURITY_MIN,
		},
		reply: kline.ReplyMinK{
			Data: &kline.MinK{},
		},
	}
}

// 获取分钟K线
func (this *MinKLine) GetMinKLine(request *kline.RequestMinK) (*kline.ReplyMinK, error) {
	key := fmt.Sprintf(this.CacheKey, request.SID)

	ls, err := redis.LRange(key, 0, -1)
	if err != nil {
		return nil, err
	}
	if len(ls) == 0 {
		return nil, REDIS_ERROR_LIST_NULL
	}

	this.reply.Code = 200
	this.reply.Data.SID = request.SID
	this.reply.Data.BeginTime = request.BeginTime
	this.reply.Data.Num = int32(len(ls))

	for _, v := range ls {
		a := &kline.KInfo{
			NSID:     1,
			NTime:    1,
			NPreCPx:  1,
			NOpenPx:  1,
			NHighPx:  1,
			NLowPx:   1,
			NLastPx:  1,
			LlVolume: 1,
			LlValue:  1,
			NAvgPx:   1,
		}
		this.reply.Data.List = append(this.reply.Data.List, a)
	}

	logging.Info("%v", this.reply)

	return &this.reply, nil
}
