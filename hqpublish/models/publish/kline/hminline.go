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
	if err != nil {
		if err != publish.ERROR_REDIS_LIST_NULL {
			return nil, err
		}
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
	var table = kline.KInfoTable{}
	bs, err := GetCache(key)
	if err != nil { //错误或没找到
		if err = getHMinlineFromeRedisStore(key, &table); err != nil {
			return nil, err
		}

		if len(table.List) < 1 {
			return nil, publish.READ_REDIS_STORE_NULL
		}

		if err = this.setPaylodToRedisCache(key, req.Type, &table); err != nil {
			logging.Error("%v", err.Error())
		}
	} else {
		if err = proto.Unmarshal(bs, &table); err != nil {
			return nil, err
		}
	}
	return &table, nil
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
	var dir, filename string

	market := req.SID / 1000000
	if market == 100 {
		dir = fmt.Sprintf("%s/sh/%d/", FStore.Path, req.SID)
	} else if market == 200 {
		dir = fmt.Sprintf("%s/sz/%d/", FStore.Path, req.SID)
	} else {
		return nil, publish.INVALID_FILE_PATH
	}

	switch protocol.HAINA_KLINE_TYPE(req.Type) {
	case protocol.HAINA_KLINE_TYPE_KMIN1:
		filename = dir + FStore.Min
		break
	case protocol.HAINA_KLINE_TYPE_KMIN5:
		filename = dir + FStore.Min5
		break
	case protocol.HAINA_KLINE_TYPE_KMIN15:
		filename = dir + FStore.Min15
		break
	case protocol.HAINA_KLINE_TYPE_KMIN30:
		filename = dir + FStore.Min30
		break
	case protocol.HAINA_KLINE_TYPE_KMIN60:
		filename = dir + FStore.Min60
		break
	default:
		return nil, publish.INVALID_REQUEST_PARA
	}

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
	//var cache *kline.HMinLineDay //指针
	//	var cacheTable = &kline.HMinTable{}

	//var tmpTime int32

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
		//		if tmpTime == 0 {
		//			cache = &kline.HMinLineDay{}
		//			cache.List = append(cache.List, &kinfo)
		//			tmpTime = kinfo.NTime / 10000
		//			continue
		//		}

		//		if tmpTime == kinfo.NTime/10000 {
		//			cache.List = append(cache.List, &kinfo)
		//			if i == lengh-size { //防止最后一条数据丢失、
		//				cache.Date = 20*1000000 + tmpTime
		//				cacheTable.List = append(cacheTable.List, cache)
		//			}
		//		} else {
		//			cache.Date = 20*1000000 + tmpTime
		//			cacheTable.List = append(cacheTable.List, cache) //得到了一天的分钟线，cache加入cacheTable

		//			cache = &kline.HMinLineDay{}            //重新创建HMinLineDay结构，cache指向它
		//			cache.List = append(cache.List, &kinfo) //本次kinfo加入新cache
		//			tmpTime = kinfo.NTime / 10000           //更新时间
		//		}
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
