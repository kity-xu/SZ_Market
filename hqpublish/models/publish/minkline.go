package publish

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	. "haina.com/share/models"

	"ProtocolBuffer/format/redis/pbdef/kline"

	redigo "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

var _ = fmt.Println
var _ = redigo.Bytes
var _ = logging.Info

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
		return nil, ERROR_REDIS_LIST_NULL
	}

	if request.BeginTime > 150100 {
		return nil, ERROR_KLINE_BEGIN_TIME
	}

	for _, v := range ls {
		k := &kline.KInfo{}
		buffer := bytes.NewBuffer([]byte(v))
		if err := binary.Read(buffer, binary.LittleEndian, k); err != nil && err != io.EOF {
			return nil, err
		}
		if k.NTime >= request.BeginTime {
			//fmt.Printf("%#v\n", k)
			this.reply.Data.List = append(this.reply.Data.List, k)
		}
	}

	/*
		if this.reply.Data.GetList() == nil {
			return nil, ERROR_KLINE_DATA_NULL
		} // */

	this.reply.Code = 200
	this.reply.Data.SID = request.SID
	this.reply.Data.BeginTime = request.BeginTime
	this.reply.Data.Num = int32(len(this.reply.Data.List))

	return &this.reply, nil
}
