//K线
package publish

import (
	"ProtocolBuffer/format/kline"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"

	"ProtocolBuffer/projects/hqpublish/go/protocol"

	sm "haina.com/share/models"

	. "haina.com/market/hqpublish/models"
	"haina.com/share/lib"

	"github.com/golang/protobuf/proto"

	redigo "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
)

var _ = fmt.Println
var _ = redigo.Bytes
var _ = logging.Info

type KLine struct {
	sm.Model `db:"-"`
}

func NewKLine(rediskey string) *KLine {
	return &KLine{
		sm.Model: sm.Model{
			CacheKey: rediskey,
		},
	}
}

// 获取某一NSID的某类历史K线（全部）
func (this *KLine) GetHisKLineAll(req *protocol.RequestHisK) (*[]*protocol.KInfo, error) {
	key := fmt.Sprintf(this.CacheKey, req.SID)

	list, err := this.getHisKlineFromeRedisCache(key, req)
	if err != nil {
		if err != ERROR_REDIS_LIST_NULL {
			return nil, err
		}
		list, err = this.getHisKlineFromeFileStore(key, req)
		if err != nil {
			return nil, err
		}
	}
	if list == nil {
		return nil, ERROR_KLINE_DATA_NULL
	}

	var table []*protocol.KInfo
	for _, v := range list.List {
		pro := &protocol.KInfo{
			NSID:     v.NSID,
			NTime:    v.NTime,
			NPreCPx:  v.NPreCPx,
			NOpenPx:  v.NOpenPx,
			NHighPx:  v.NHighPx,
			NLowPx:   v.NLowPx,
			NLastPx:  v.NLastPx,
			LlVolume: v.LlVolume,
			LlValue:  v.LlValue,
			NAvgPx:   v.NAvgPx,
		}
		table = append(table, pro)
	}

	return &table, nil
}

//从数据redis全部取出
func getHisKlineFromeRedisStore(key string) (*kline.KInfoTable, error) {
	ss, err := RedisStore.LRange(key, 0, -1)
	if err != nil {
		return nil, err
	}

	if len(ss) == 0 {
		return nil, ERROR_REDIS_LIST_NULL
	}

	var table = &kline.KInfoTable{}
	for _, v := range ss {
		kinfo := &kline.KInfo{}
		if err := proto.Unmarshal([]byte(v), kinfo); err != nil {
			return nil, err
		}
		table.List = append(table.List, kinfo)
	}
	return table, nil
}

//从文件获取数据
func (this *KLine) getHisKlineFromeFileStore(key string, req *protocol.RequestHisK) (*kline.KInfoTable, error) {
	var dir, filename string

	market := req.SID / 1000000
	if market == 100 {
		dir = fmt.Sprintf("%s/sh/%d/", FStore.Path, req.SID)
	} else if market == 200 {
		dir = fmt.Sprintf("%s/sz/%d/", FStore.Path, req.SID)
	} else {
		return nil, INVALID_FILE_PATH
	}

	switch protocol.HAINA_KLINE_TYPE(req.Type) {
	case protocol.HAINA_KLINE_TYPE_KDAY:
		filename = dir + FStore.Day
		if !lib.IsFileExist(filename) {
			filename = dir + FStore.Index
		}
		break
	case protocol.HAINA_KLINE_TYPE_KWEEK:
		filename = dir + FStore.Week
		break
	case protocol.HAINA_KLINE_TYPE_KMONTH:
		filename = dir + FStore.Month
		break
	case protocol.HAINA_KLINE_TYPE_KYEAR:
		filename = dir + FStore.Year
		break
	default:
		return nil, INVALID_REQUEST_PARA
	}

	if !lib.IsFileExist(filename) {
		return nil, INVALID_FILE_PATH
	}

	///do something
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var line kline.KInfo
	size := binary.Size(&line)

	var table = &kline.KInfoTable{}

	for i := 0; i < len(data); i += size {
		v := data[i : i+size]
		kinfo := kline.KInfo{}
		buffer := bytes.NewBuffer(v)
		if err = binary.Read(buffer, binary.LittleEndian, &kinfo); err != nil && err != io.EOF {
			return nil, err
		}
		table.List = append(table.List, &kinfo)
	}

	if err = this.setPaylodToRedisCache(key, req.Type, table); err != nil { //读文件后的数据入缓冲
		logging.Error("%v", err.Error())
	}
	return table, nil
}

func (this *KLine) setPaylodToRedisCache(key string, stype int32, table *kline.KInfoTable) error {
	var ttl int
	switch protocol.HAINA_KLINE_TYPE(stype) {
	case protocol.HAINA_KLINE_TYPE_KDAY:
		ttl = TTL.Day
		break
	case protocol.HAINA_KLINE_TYPE_KWEEK:
		ttl = TTL.Week
		break
	case protocol.HAINA_KLINE_TYPE_KMONTH:
		ttl = TTL.Month
		break
	case protocol.HAINA_KLINE_TYPE_KYEAR:
		ttl = TTL.Year
		break
	default:
		break
	}

	data, err := proto.Marshal(table)
	if err != nil {
		return err
	}
	if err = SetCache(key, ttl, data); err != nil {
		return err
	}
	return nil
}

//从缓冲redis获取数据
func (this *KLine) getHisKlineFromeRedisCache(key string, req *protocol.RequestHisK) (*kline.KInfoTable, error) {
	var ktable = &kline.KInfoTable{}
	bs, err := GetCache(key)
	if err != nil { //错误或没找到
		ktable, err = getHisKlineFromeRedisStore(key)
		if err != nil {
			return nil, err
		}

		if err = this.setPaylodToRedisCache(key, req.Type, ktable); err != nil {
			logging.Error("%v", err.Error())
		}
	} else {
		if err = proto.Unmarshal(bs, ktable); err != nil {
			return nil, err
		}
	}
	return ktable, nil
}
