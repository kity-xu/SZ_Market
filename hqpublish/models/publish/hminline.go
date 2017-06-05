//分时数据（历史分钟线）
package publish

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"

	"haina.com/share/lib"

	. "haina.com/share/models"

	"ProtocolBuffer/projects/hqpublish/go/protocol"

	//"haina.com/share/logging"

	"github.com/golang/protobuf/proto"
	. "haina.com/market/hqpublish/models"
)

type HisMinLine struct {
	Model `db:"-"`
}

func NewHisMinLine(redis_key string) *HisMinLine {
	return &HisMinLine{
		Model: Model{
			CacheKey: redis_key,
		},
	}
}

func (this *HisMinLine) PayloadHisMinLine(req *protocol.RequestHisMinK) (*protocol.PayloadHisMinK, error) {
	table, err := this.getHisLineMin_01(req)

	if err != nil || len(table.List) < 1 {
		return nil, err
	}

	var tmpTime, count, day, sid int32
	var kinfos []*protocol.KInfo

	for i := len(table.List) - 1; i > 0; i-- {
		if tmpTime == 0 {
			tmpTime = table.List[i].NTime / 10000
			kinfos = append(kinfos, table.List[i])
			sid = table.List[i].NSID
			count++
			continue
		}

		if tmpTime != table.List[i].NTime/10000 {
			day++
			if req.DayNum == day {
				break
			}
		}

		count++
		kinfos = append(kinfos, table.List[i])
		tmpTime = table.List[i].NTime / 10000
	}

	GetASCStruct(&kinfos) //升序排序

	hismin := &protocol.PayloadHisMinK{
		SID:    sid,
		DayNum: day,
		Num:    count,
		KList:  kinfos,
	}
	return hismin, nil
}

// 获取某一NSID的某类历史分钟线（全部）
func (this *HisMinLine) getHisLineMin_01(req *protocol.RequestHisMinK) (*protocol.KInfoTable, error) {
	key := fmt.Sprintf(this.CacheKey, req.SID)
	var (
		blist []byte
		err   error
	)

	blist, err = GetCache(key)
	if err != nil {
		blist, err = this.getHisMin01_FromeFileStore(req)
		if err != nil {
			return nil, err
		}

		ks, err := bytesUnMarshal(blist)
		if err != nil {
			return nil, err
		}

		table := &protocol.KInfoTable{
			List: *ks,
		}

		bs, err := proto.Marshal(table)
		if err != nil {
			return nil, err
		}

		SetCache(key, TTL.Min, bs)

		return table, nil
	} else {
		table := &protocol.KInfoTable{}
		err = proto.Unmarshal(blist, table)
		if err != nil {
			return nil, err
		}
		return table, nil
	}
}

func bytesUnMarshal(data []byte) (*[]*protocol.KInfo, error) {
	var table []*protocol.KInfo
	var line protocol.KInfo

	size := binary.Size(&line)

	for i := 0; i < len(data); i += size {
		v := data[i : i+size]
		kinfo := protocol.KInfo{}
		buffer := bytes.NewBuffer(v)
		if err := binary.Read(buffer, binary.LittleEndian, &kinfo); err != nil && err != io.EOF {
			return nil, err
		}
		table = append(table, &kinfo)
	}
	return &table, nil
}

func (this *HisMinLine) getHisMin01_FromeFileStore(req *protocol.RequestHisMinK) ([]byte, error) {
	var dir, filename string

	market := req.SID / 1000000
	if market == 100 {
		dir = fmt.Sprintf("%s/sh/%d/", FStore.Path, req.SID)
	} else if market == 200 {
		dir = fmt.Sprintf("%s/sz/%d/", FStore.Path, req.SID)
	} else {
		return nil, INVALID_FILE_PATH
	}

	filename = dir + FStore.Min

	if !lib.IsFileExist(filename) {
		return nil, INVALID_FILE_PATH
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return data, nil
}
