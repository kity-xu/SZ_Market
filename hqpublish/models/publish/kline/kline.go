//K线
package kline

import (
	"ProtocolBuffer/format/kline"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"time"

	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"haina.com/market/hqpublish/models/publish"
	"haina.com/market/hqpublish/models/publish/security"
	sm "haina.com/share/models"

	. "haina.com/market/hqpublish/models"
	"haina.com/share/lib"

	"github.com/golang/protobuf/proto"
	"haina.com/share/logging"
)

type KLine struct {
	sm.Model `db:"-"`
}

func NewKLine(rediskey string) *KLine {
	return &KLine{
		Model: sm.Model{
			CacheKey: rediskey,
		},
	}
}

// 获取某一NSID的某类历史K线（全部）
func (this *KLine) GetHisKLineAll(req *protocol.RequestHisK) (*[]*protocol.KInfo, error) {
	key := fmt.Sprintf(this.CacheKey, req.SID)
	list, err := this.getHisKlineFromeRedisCache(key, req)
	if list == nil {
		list, err = this.getHisKlineFromeFileStore(key, req)
		if err != nil {
			var table []*protocol.KInfo
			return &table, err
		}
	}
	if list == nil {
		return nil, publish.ERROR_KLINE_DATA_NULL
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
		return nil, publish.ERROR_REDIS_LIST_NULL
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
	var kind, filename string

	switch protocol.HAINA_KLINE_TYPE(req.Type) {
	case protocol.HAINA_KLINE_TYPE_KDAY:
		kind = FStore.Day
	case protocol.HAINA_KLINE_TYPE_KWEEK:
		kind = FStore.Week
	case protocol.HAINA_KLINE_TYPE_KMONTH:
		kind = FStore.Month
	case protocol.HAINA_KLINE_TYPE_KYEAR:
		kind = FStore.Year
	default:
		return nil, publish.INVALID_REQUEST_PARA
	}

	market := req.SID / 1000000
	if market == 100 {
		filename = fmt.Sprintf("%s/sh/%s/%d.dat", FStore.Path, kind, req.SID)
	} else if market == 200 {
		filename = fmt.Sprintf("%s/sz/%s/%d.dat", FStore.Path, kind, req.SID)
	} else {
		return nil, publish.INVALID_FILE_PATH
	}

	if !lib.IsFileExist(filename) {
		return nil, publish.INVALID_FILE_PATH
	}

	///do something
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, publish.FILE_HMINDATA_NULL
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
	if err == nil {
		if len(bs) == 0 {
			return nil, publish.ERROR_REDIS_LIST_NULL
		}
		if err = proto.Unmarshal(bs, ktable); err != nil {
			return nil, err
		}
	} else {
		ktable = nil
	}
	return ktable, nil
}

func getHQpostExecutedTime() (string, error) {

	ss, err := RedisStore.Get(REDISKEY_HQPOST_EXECUTED_TIME)
	if err != nil {
		return "", err
	}
	if ss == "" {
		return "", fmt.Errorf("redis get HQpostExecutedTime is null..")
	}
	return ss, nil
}

func IsHQpostRunOver() (bool, error) {
	ss, err := getHQpostExecutedTime()
	if err != nil {
		logging.Error("%v", err.Error())
		return false, err
	}

	dd, err := strconv.ParseInt(ss, 10, 64)
	if err != nil {
		logging.Error("%v", err.Error())
		return false, err
	}

	timestamp := time.Now().Unix()

	tm := time.Unix(timestamp, 0)
	format := tm.Format("200601021504")
	monment, err := strconv.ParseInt(format, 10, 64)
	logging.Info("monment:%v-----dd:%v", monment, dd)
	if err != nil {
		logging.Error("%v", err.Error())
		return false, err
	}

	if dd%10000 < monment%10000 { //hqpost更新完毕
		return true, nil
	}
	return false, nil
}

func IsExistInRedis(key string) bool {
	bs, err := RedisStore.GetBytes(key)
	if err == nil && len(bs) > 0 {
		return true
	}
	return false
}

// 判断某支股票是否停牌
func IsDelist(sid int32) bool {
	key := "hq:st:static:%d"
	key_sc := fmt.Sprintf(key, sid)

	static := &security.StockStatic{}
	bs, err := RedisStore.GetBytes(key_sc)
	if err != nil {
		logging.Error("%v", err.Error())
		return true //按停牌算
	}

	binary.Read(bytes.NewBuffer(bs), binary.LittleEndian, static)
	if static.SzStatus[0] == 'S' { //停牌
		return true
	}
	return false
}

func IsIndex(sid int32) (bool, error) {
	key := fmt.Sprintf("hq:st:name:%d", sid)
	bs, err := RedisStore.GetBytes(key)
	if err != nil {
		return false, err
	}

	var kinfo = &security.TagSecurityName{}

	binary.Read(bytes.NewBuffer(bs), binary.LittleEndian, kinfo)
	if kinfo.SzSType[1] == 'I' {
		return true, nil
	}
	return false, nil
}
