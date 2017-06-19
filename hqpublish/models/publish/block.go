//板块
package publish

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"

	. "haina.com/share/models"

	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"haina.com/share/logging"

	"github.com/golang/protobuf/proto"
	. "haina.com/market/hqpublish/models"
)

//TagStockSortInfo

//板块排序结构体
type TagBlockSortInfo struct {
	NBlockID        int32                   //板块ID
	SzBlockName     [12]byte                //板块名称
	NAveChgRate     int32                   //平均涨跌幅*10000
	LlVolume        int64                   ///< 板块总成交量
	LlValue         int64                   ///< 板块总成交额(*10000)
	SzSName         [SECURITY_NAME_LEN]byte //领涨股名称
	NNum            int32                   //成分股票个数
	NChgRatio       int32                   //上涨比例
	NLongNum        int32                   //上涨家数
	NShortNum       int32                   //下跌家数
	LlValueOfInFlow int64                   ///<资金流入额(*10000)
}

var (
	REDIS_KEY_CACHE_BLOCK = "hq:init:bk:%d"
)

type Block struct {
	Model `db:"-"`
}

func NewBlock(redis_key string) *Block {
	return &Block{
		Model: Model{
			CacheKey: redis_key,
		},
	}
}

func (this *Block) GetBlockReplyByRequest(req *protocol.RequestBlock) (*protocol.PayloadBlock, error) {
	//按所传参数做数据解析
	if req.Begin < 0 || req.Num < 0 {
		return nil, INVALID_REQUEST_PARA
	}

	var blocks []*protocol.TagBlockSortInfo

	ckey := fmt.Sprintf(REDIS_KEY_CACHE_BLOCK, req.Classify)
	data, err := RedisCache.GetBytes(ckey)
	if err != nil {
		logging.Debug("cache redis is nil...%v", err.Error())
		if err = this.GetBlockFromeRediaData(req, &blocks); err != nil {
			logging.Error("%v", err.Error())
			return nil, err
		}
	} else {
		blist := &protocol.BlockList{}
		if err = proto.Unmarshal(data, blist); err != nil {
			logging.Error("-----------%v", err.Error())
			return nil, err
		}

		dkey := fmt.Sprintf(REDISKEY_SORT_KDAY_H, protocol.HAINA_PUBLISH_SORT_BLOCKID_BK_S, absInt32(req.FieldID))
		dblock, err := RedisStore.GetBytes(dkey)
		if err != nil {
			logging.Error("---***%v", err.Error())
			return nil, err
		}

		bk := &TagBlockSortInfo{}
		size := binary.Size(bk)

		for i := 0; i < len(dblock); i += size {
			block := &TagBlockSortInfo{}
			if err = binary.Read(bytes.NewBuffer(dblock[i:i+size]), binary.LittleEndian, block); err != nil && err != io.EOF {
				logging.Error("---%v", err.Error())
				return nil, err
			}

			for _, v := range blist.List {
				if block.NBlockID == v.SetID {
					pbk := &protocol.TagBlockSortInfo{
						NBlockID:        block.NBlockID,
						SzBlockName:     byte12ToString(block.SzBlockName),
						NAveChgRate:     block.NAveChgRate,
						LlVolume:        block.LlVolume,
						LlValue:         block.LlValue,
						SzSName:         byte40ToString(block.SzSName),
						NNum:            block.NNum,
						NChgRatio:       block.NChgRatio,
						NLongNum:        block.NLongNum,
						NShortNum:       block.NShortNum,
						LlValueOfInFlow: block.LlValueOfInFlow,
					}
					blocks = append(blocks, pbk)
					break
				}
			}
		}
	}

	if len(blocks)-1 < int(req.Begin) {
		return nil, INVALID_REQUEST_PARA
	}

	if req.FieldID < 0 {
		swapBlock(&blocks)
	}

	var board []*protocol.TagBlockSortInfo

	if req.Num == 0 {
		board = blocks[0:]
	} else {
		if req.Num+req.Begin < int32(len(blocks)) {
			board = blocks[req.Begin : req.Num+req.Begin]
		} else {
			board = blocks[req.Begin:]
		}
	}

	payload := &protocol.PayloadBlock{
		Classify: req.Classify,
		FieldID:  req.FieldID,
		Total:    int32(len(blocks)),
		Begin:    req.Begin,
		Num:      int32(len(board)),
		List:     board,
	}
	return payload, nil
}

func (this *Block) GetBlockFromeRediaData(req *protocol.RequestBlock, blocks *[]*protocol.TagBlockSortInfo) error {
	key := fmt.Sprintf(this.CacheKey, req.Classify)
	data, err := RedisStore.GetBytes(key)
	if err != nil {
		return err
	}
	var simp = &protocol.BlockList{}
	if err = proto.Unmarshal(data, simp); err != nil {
		return err
	}

	skey := fmt.Sprintf(REDISKEY_SORT_KDAY_H, protocol.HAINA_PUBLISH_SORT_BLOCKID_BK_S, absInt32(req.FieldID))

	db, err := RedisStore.GetBytes(skey)
	if err != nil {
		logging.Error("---***%v", err.Error())
		return err
	}

	bk := &TagBlockSortInfo{}
	size := binary.Size(bk)

	blist := &protocol.BlockList{} //cache

	for i := 0; i < len(db); i += size {
		block := &TagBlockSortInfo{}
		if err = binary.Read(bytes.NewBuffer(db[i:i+size]), binary.LittleEndian, block); err != nil && err != io.EOF {
			logging.Error("---%v", err.Error())
			return err
		}

		for _, v := range simp.List {
			if block.NBlockID == v.SetID {
				pbk := &protocol.TagBlockSortInfo{
					NBlockID:        block.NBlockID,
					SzBlockName:     byte12ToString(block.SzBlockName),
					NAveChgRate:     block.NAveChgRate,
					LlVolume:        block.LlVolume,
					LlValue:         block.LlValue,
					SzSName:         byte40ToString(block.SzSName),
					NNum:            block.NNum,
					NChgRatio:       block.NChgRatio,
					NLongNum:        block.NLongNum,
					NShortNum:       block.NShortNum,
					LlValueOfInFlow: block.LlValueOfInFlow,
				}
				*blocks = append(*blocks, pbk)

				cacheB := &protocol.Block{
					SetID:   block.NBlockID,
					SetName: v.SetName,
				}
				blist.List = append(blist.List, cacheB)
				break
			}
		}
	}

	//入Cache
	dCache, err := proto.Marshal(blist)
	if err != nil {
		return err
	}
	ckey := fmt.Sprintf(REDIS_KEY_CACHE_BLOCK, req.Classify)

	if err = RedisCache.Set(ckey, dCache); err != nil {
		return err
	}

	if err = lifeTime(ckey, TTL.Block); err != nil {
		return err
	}

	if len(*blocks) < 0 {
		return errors.New("There is no plate data...")
	}

	//logging.Debug("len:%v", len(*blocks))
	return nil
}

func lifeTime(key string, tm int) error {
	now := time.Now()
	nowhm := now.Hour()*100 + now.Minute()

	// 计算TTL: 当前时间到下一个9:25之间的秒数
	stop := now
	if nowhm >= tm {
		stop = stop.AddDate(0, 0, 1)
	}

	local, _ := time.LoadLocation("Local")
	stopstr := fmt.Sprintf("%04d-%02d-%02d %02d:%02d", stop.Year(), int(stop.Month()), stop.Day(), 9, 25)
	ttlstop, _ := time.ParseInLocation("2006-01-02 15:04", stopstr, local)
	ttls := ttlstop.Unix()

	if _, err := RedisCache.Do("EXPIREAT", key, ttls); err != nil { // 缓存Redis TTL设置 下一个9:25自动删除
		return err
	}
	return nil
}

//首尾交换未知
func swapBlock(table *[]*protocol.TagBlockSortInfo) {
	lengh := len(*table)

	for i := 0; i < lengh; i++ {
		(*table)[i], (*table)[lengh-i-1] = (*table)[lengh-i-1], (*table)[i]
		if i == lengh-i-2 || i == lengh-i-3 {
			break
		}
	}
}