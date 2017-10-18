package redistore

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"haina.com/share/logging"

	"haina.com/share/lib"

	"haina.com/market/hqpost/models"
	. "haina.com/share/models"

	"ProtocolBuffer/projects/hqpost/go/protocol"

	"haina.com/share/store/redis"
)

type MBlockIndex struct {
	Model `db:"-"`
}

func NewMBlockIndex(key string) *MBlockIndex {
	return &MBlockIndex{
		Model: Model{
			CacheKey: key,
		},
	}
}

func (this *MBlockIndex) getBlockID() ([]string, error) {
	keys, err := redis.Keys(this.CacheKey)
	if err != nil {
		return nil, err
	}

	var bids []string
	for _, key := range keys {
		ss := strings.Split(key, ":")
		if len(ss) < 4 {
			return nil, err
		}
		bids = append(bids, ss[3])
	}
	return bids, nil
}

func (this *MBlockIndex) UpdateBlockIndexLine(path string) error {
	bids, err := this.getBlockID()
	if err != nil {
		return err
	}

	for _, bid := range bids {
		table, err := this.getBlockIndexKLineTable(bid, path)
		if err != nil || len(*table) < 1 {
			logging.Error("%v", err.Error())
			continue
		}
		day := this.mergeMinlineTODay(table, bid)
		if err = this.appendBlockIndexFile(day, path+bid); err != nil {
			return err
		}
	}

	return nil
}

func (this *MBlockIndex) getBlockIndexKLineTable(bid, path string) (*[]*protocol.KInfo, error) {
	bkey := fmt.Sprintf(models.REDISKEY_BLOCKINDEX_MIN, bid)
	//fmt.Println("-----bkey:", bkey)
	dayline, err := redis.LRange(bkey, 0, -1)
	if err != nil || len(dayline) < 1 {
		return nil, fmt.Errorf("get block:%v index minline failed...", bid)
	}
	if len(dayline) < 1 {
		logging.Debug("block:%v minline is null", bid)
		return nil, nil
	}

	table := make([]*protocol.KInfo, 0, 256)
	for _, min := range dayline {
		mk := &protocol.KInfo{}
		if err = binary.Read(bytes.NewBuffer([]byte(min)), binary.LittleEndian, mk); err != nil && err != io.EOF {
			return nil, err
		}
		table = append(table, mk)
	}
	return &table, nil
}

// mergeMinlineTODay 合并分钟线
func (this *MBlockIndex) mergeMinlineTODay(table *[]*protocol.KInfo, bid string) *protocol.KInfo {
	var (
		i          int
		AvgPxTotal uint32
		tmp        protocol.KInfo //pb类型
	)

	Bid, _ := strconv.Atoi(bid)

	for _, v := range *table {
		if tmp.NHighPx < v.NHighPx || tmp.NHighPx == 0 { //最高价
			tmp.NHighPx = v.NHighPx
		}
		if tmp.NLowPx > v.NLowPx || tmp.NLowPx == 0 { //最低价
			tmp.NLowPx = v.NLowPx
		}
		tmp.LlVolume += v.LlVolume //成交量
		tmp.LlValue += v.LlValue   //成交额
		AvgPxTotal += v.NAvgPx

		i++
	}
	tmp.NSID = int32(Bid)
	tmp.NTime = (*table)[0].NTime/10000 + 20000000 //时间
	tmp.NOpenPx = (*table)[0].NOpenPx              //开盘价
	tmp.NLastPx = (*table)[len(*table)-1].NLastPx  //最新价
	tmp.NAvgPx = AvgPxTotal / uint32(i+1)          //平均价

	return &tmp
}

func (this *MBlockIndex) appendBlockIndexFile(day *protocol.KInfo, path string) error {
	var precpx int32 //昨收价

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	if !lib.IsFileExist(path) {
		precpx = 1000 * 10000

	} else { //存在
		kinfo := &protocol.KInfo{}
		size := binary.Size(kinfo)
		if _, err = file.Seek(-int64(size), 2); err != nil {
			return err
		}
		bs := make([]byte, size)
		if _, err = file.Read(bs); err != nil && err != io.EOF {
			return err
		}

		if err = binary.Read(bytes.NewBuffer(bs), binary.LittleEndian, kinfo); err != nil && err != io.EOF {
			return err
		}
		precpx = kinfo.NLastPx
	}
	day.NPreCPx = precpx

	buff := new(bytes.Buffer)
	if err = binary.Write(buff, binary.LittleEndian, day); err != nil {
		return err
	}

	if _, err = file.Write(buff.Bytes()); err != nil {
		return err
	}
	return nil
}
