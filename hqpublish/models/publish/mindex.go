package publish

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	ctrl "haina.com/market/hqpublish/controllers"
	. "haina.com/market/hqpublish/models"
	. "haina.com/share/models"

	pro "ProtocolBuffer/projects/hqpublish/go/protocol"

	hsgrr "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
	"haina.com/share/store/redis"

	"github.com/golang/protobuf/proto"
)

var (
	_ = redis.Init
	_ = GetCache
	_ = ctrl.MakeRespDataByBytes
	_ = errors.New
	_ = fmt.Println
	_ = hsgrr.Bytes
	_ = logging.Info
	_ = bytes.NewBuffer
	_ = binary.Read
	_ = io.ReadFull
)

type MIndex struct {
	Model `db:"-"`
}

func NewMIndex() *MIndex {
	return &MIndex{
		Model: Model{
			CacheKey: REDISKEY_L2CACHE_INDEX_MOBILE,
		},
	}
}

func (this MIndex) GetMIndexObj() (*pro.PayloadMIndex, error) {
	b, err := RedisCache.GetBytes(this.CacheKey)
	if err == nil {
		ob := pro.PayloadMIndex{}
		err := proto.Unmarshal(b, &ob)
		if err == nil {
			return &ob, nil
		}
		logging.Warning("%v", err)
	}

	index := pro.PayloadMIndex{
		InfoList:     make([]*pro.Infobar, 0, 3),
		HotBlockList: make([]*pro.TagBlockSortInfo, 0, 3),
		IncrList:     make([]*pro.TagStockSortInfo, 0, 5),
		DeclList:     make([]*pro.TagStockSortInfo, 0, 5),
		InflowList:   make([]*pro.TagStockSortInfo, 0, 5),
		OutflowList:  make([]*pro.TagStockSortInfo, 0, 5),
	}

	{ //InfoList 信息栏
		var req pro.RequestSnapshot
		req.SID = 100000001 // 上证指数
		datash, err := NewIndexSnapshot().GetIndexSnapshotObj(&req)
		if err != nil {
			logging.Error("%v", err)
			return nil, err
		}
		req.SID = 200399001 // 深圳成指
		datasz, err := NewIndexSnapshot().GetIndexSnapshotObj(&req)
		if err != nil {
			logging.Error("%v", err)
			return nil, err
		}
		req.SID = 200399006 // 创业板
		datacy, err := NewIndexSnapshot().GetIndexSnapshotObj(&req)
		if err != nil {
			logging.Error("%v", err)
			return nil, err
		}
		index.InfoList = append(index.InfoList, DataTreating(1, datash))
		index.InfoList = append(index.InfoList, DataTreating(2, datasz))
		index.InfoList = append(index.InfoList, DataTreating(3, datacy))
	}

	{ //HotBlockList // 板块排序 热点
		req := pro.RequestBlock{
			Classify: 1100,
			FieldID:  -4006,
			Begin:    0,
			Num:      3,
		}
		block, err := NewBlock(REDISKEY_BLOCK).GetBlockReplyByRequest(&req)
		if err != nil {
			logging.Error("%v", err)
			return nil, err
		}
		index.HotBlockList = append(index.HotBlockList, block.List[:]...)
	}

	{ //IncrList   // 排序 涨幅
		req := pro.RequestSort{
			SetID:   2,
			FieldID: -2008,
			Begin:   0,
			Num:     5,
		}
		sort, err := NewSort(REDISKEY_SORT_KDAY_H).GetPayloadSort(&req)
		if err != nil {
			logging.Error("%v", err)
			return nil, err
		}
		index.IncrList = append(index.IncrList, sort.List[:]...)
	}

	{ //DeclList   // 排序 跌幅
		req := pro.RequestSort{
			SetID:   2,
			FieldID: 2008,
			Begin:   0,
			Num:     5,
		}
		sort, err := NewSort(REDISKEY_SORT_KDAY_H).GetPayloadSort(&req)
		if err != nil {
			logging.Error("%v", err)
			return nil, err
		}
		index.DeclList = append(index.DeclList, sort.List[:]...)
	}

	{ //InflowList // 排序 流入
		req := pro.RequestSort{
			SetID:   2,
			FieldID: -2025,
			Begin:   0,
			Num:     5,
		}
		sort, err := NewSort(REDISKEY_SORT_KDAY_H).GetPayloadSort(&req)
		if err != nil {
			logging.Error("%v", err)
			return nil, err
		}
		index.InflowList = append(index.InflowList, sort.List[:]...)
	}

	{ //OutflowList// 排序 流出
		req := pro.RequestSort{
			SetID:   2,
			FieldID: 2025,
			Begin:   0,
			Num:     5,
		}
		sort, err := NewSort(REDISKEY_SORT_KDAY_H).GetPayloadSort(&req)
		if err != nil {
			logging.Error("%v", err)
			return nil, err
		}
		index.OutflowList = append(index.OutflowList, sort.List[:]...)
	}

	bs, err := proto.Marshal(&index)
	if err == nil {
		RedisCache.Setex(this.CacheKey, TTL.MIndex, []byte(bs))
	} else {
		logging.Warning("%v", err)
	}

	return &index, nil
}

// 处理证券快照
func DataTreating(ind int, pst *pro.IndexSnapshot) *pro.Infobar {
	var sname = ""
	if ind == 1 {
		sname = "上证指数"
	}
	if ind == 2 {
		sname = "深圳成指"
	}
	if ind == 3 {
		sname = "创业板指"
	}
	return &pro.Infobar{
		NSID:       pst.NSID,
		SzSName:    sname,
		NLastPx:    pst.NLastPx,
		LlVolume:   pst.LlVolume,
		NPxChg:     pst.NPxChg,
		PxChgRatio: pst.NPxChgRatio,
	}
}
