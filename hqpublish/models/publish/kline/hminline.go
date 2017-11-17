package kline

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"

	"haina.com/share/lib"

	. "haina.com/share/models"

	"ProtocolBuffer/format/kline"
	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"haina.com/market/hqpublish/models/publish"
	"haina.com/share/logging"

	"github.com/golang/protobuf/proto"
	. "haina.com/market/hqpublish/models"
)

type HMinKLine struct {
	Model `db:"-"`
}

func NewHMinKLine(redis_key string) *HMinKLine {
	return &HMinKLine{
		Model: Model{
			CacheKey: redis_key,
		},
	}
}

// 获取某一NSID的某类历史分钟线（全部）
func (this *HMinKLine) GetHMinKLineAll(req *protocol.RequestHisK) (*[]*protocol.KInfo, error) {
	key := fmt.Sprintf(this.CacheKey, req.SID)
	var kinfolist *kline.KInfoTable

	kinfolist, err := this.getHMinlineFromeRedisCache(key, req)
	if kinfolist == nil {
		kinfolist, err = this.getHMinlineFromeFileStore(key, req)
		if err != nil {
			return nil, err
		}
	}

	var ktale []*protocol.KInfo
	if kinfolist != nil {
		for _, v := range kinfolist.List {
			kinfo := &protocol.KInfo{
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
			ktale = append(ktale, kinfo)
		}
	} else {
		return nil, publish.ERROR_KLINE_DATA_NULL
	}
	return &ktale, nil
}

func (this *HMinKLine) getHMinlineFromeRedisCache(key string, req *protocol.RequestHisK) (*kline.KInfoTable, error) {
	var table = &kline.KInfoTable{}
	bs, err := GetCache(key)
	if err == nil {
		if len(bs) == 0 {
			return nil, publish.ERROR_REDIS_LIST_NULL
		}
		if err = proto.Unmarshal(bs, table); err != nil {
			return nil, err
		}
	} else {
		table = nil
	}
	return table, nil
}

func getHMinlineFromeRedisStore(key string, table *kline.KInfoTable) error {
	ss, err := RedisStore.LRange(key, 0, -1)
	if err != nil {
		return err
	}

	if len(ss) == 0 {
		return publish.ERROR_REDIS_LIST_NULL
	}

	for _, v := range ss {
		kinfo := &kline.KInfo{}
		if err := proto.Unmarshal([]byte(v), kinfo); err != nil {
			return err
		}
		table.List = append(table.List, kinfo)
	}
	return nil
}

func (this *HMinKLine) getHMinlineFromeFileStore(key string, req *protocol.RequestHisK) (*kline.KInfoTable, error) {
	var kind, filename string
	logging.Info("-------************filename----------")

	switch protocol.HAINA_KLINE_TYPE(req.Type) {
	case protocol.HAINA_KLINE_TYPE_KMIN1:
		kind = FStore.Min
	case protocol.HAINA_KLINE_TYPE_KMIN5:
		kind = FStore.Min5
	case protocol.HAINA_KLINE_TYPE_KMIN15:
		kind = FStore.Min15
	case protocol.HAINA_KLINE_TYPE_KMIN30:
		kind = FStore.Min30
	case protocol.HAINA_KLINE_TYPE_KMIN60:
		kind = FStore.Min60
	default:
		return nil, publish.INVALID_REQUEST_PARA
	}
	market := req.SID / 1000000
	if market == 100 {
		filename = fmt.Sprintf("%s/sh/%s/%d.dat", FStore.Path, kind, req.SID)
	} else if market == 200 {
		filename = fmt.Sprintf("%s/sz/%s/%d.dat", FStore.Path, kind, req.SID)
	} else {
		logging.Info("---filename----------")
		return nil, publish.INVALID_FILE_PATH
	}
	logging.Info("---filename:%s", filename)

	if !lib.IsFileExist(filename) {
		return nil, publish.INVALID_FILE_PATH
	}
	///do something
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var line kline.KInfo
	size := binary.Size(&line)

	var table = &kline.KInfoTable{}
	lengh := len(data)
	if lengh == 0 {
		return nil, publish.FILE_HMINDATA_NULL
	}

	for i := 0; i < lengh; i += size {
		v := data[i : i+size]
		kinfo := kline.KInfo{}
		buffer := bytes.NewBuffer(v)
		if err = binary.Read(buffer, binary.LittleEndian, &kinfo); err != nil && err != io.EOF {
			return nil, err
		}
		table.List = append(table.List, &kinfo)
	}
	if err = this.setPaylodToRedisCache(key, req.Type, table); err != nil { //这里文件存储结构转换为redis存储结构
		logging.Error("%v", err.Error()) //此处不能因为该错误而返回，但又不能忽略错误，故打印
	}
	return table, nil
}

func (this *HMinKLine) setPaylodToRedisCache(key string, stype int32, table *kline.KInfoTable) error {
	var ttl int
	switch protocol.HAINA_KLINE_TYPE(stype) {
	case protocol.HAINA_KLINE_TYPE_KMIN1:
		ttl = TTL.Min1
		break
	case protocol.HAINA_KLINE_TYPE_KMIN5:
		ttl = TTL.Min5
		break
	case protocol.HAINA_KLINE_TYPE_KMIN15:
		ttl = TTL.Min15
		break
	case protocol.HAINA_KLINE_TYPE_KMIN30:
		ttl = TTL.Min30
		break
	case protocol.HAINA_KLINE_TYPE_KMIN60:
		ttl = TTL.Min60
		break
	default:
		return publish.INVALID_REQUEST_PARA
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
