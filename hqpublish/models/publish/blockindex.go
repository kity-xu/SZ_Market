package publish

import (
	"fmt"

	. "haina.com/share/models"

	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"haina.com/share/logging"

	"github.com/golang/protobuf/proto"
	. "haina.com/market/hqpublish/models"
	"haina.com/share/kityxu/utils"
)

type BlockIndex struct {
	Model `db:"-"`
}

func NewBlockIndex(redis_key string) *BlockIndex {
	return &BlockIndex{
		Model: Model{
			CacheKey: redis_key,
		},
	}
}

//protocol.PayloadBlockindex ------GetPayloadBlockindex

func (this *BlockIndex) GetPayloadBlockindex(req *protocol.RequestBlockindex) (*protocol.PayloadBlockindex, error) {
	table, err := this.GetEleStockListXRXD(req)
	if err != nil {
		logging.Error("%v", err.Error())
		return nil, err
	}
	//logging.Debug("len(table):%v", (*table)[0].List[0])

	var mostkinfos *protocol.KInfoTable
	var max int
	for _, ele := range *table { //找出所有成份股中交易日最全的，以此为板块交易日基准
		logging.Debug("len:%v", len(ele.List))
		if len(ele.List) > max {
			mostkinfos = ele
			max = len(ele.List)
		}
	}
	logging.Debug("----len:%v", len(mostkinfos.List))

	var blockele map[int32][]*protocol.KInfo

	for _, v := range mostkinfos.List {
		for i, ele := range *table {
			if len(ele.List) == 0 {
				logging.Info("没有这只成份股")
				break
			}
			if len(ele.List) < i+1 {
				logging.Info("该股票已停牌%v", ele.List[0].NSID)
				break
			}
			if ele.List[i].NTime < v.NTime {
				continue
			} else if ele.List[i].NTime == v.NTime {
				blockele[v.NTime] = append(blockele[v.NTime], ele.List[i])
			} else {
				logging.Info("该交易日没有%v这只股票", ele.List[0].NSID)
				break
			}
		}
	}

	return nil, nil
}

// 获取该板块ID下的所有除权后的成分股集合
func (this *BlockIndex) GetEleStockListXRXD(req *protocol.RequestBlockindex) (*[]*protocol.KInfoTable, error) {
	bkey := fmt.Sprintf(this.CacheKey, 1100, req.SetID)

	block := &protocol.ElementList{} //板块

	data, err := RedisStore.GetBytes(bkey)
	if err != nil {
		logging.Error("%v", err.Error())
		return nil, err
	}

	if err = proto.Unmarshal(data, block); err != nil {
		logging.Error("%v", err.Error())
		return nil, err
	}

	rd := &protocol.RequestXRXD{
		Method: 1,
	}
	index := getLrangeIndex()

	ktable := make([]*protocol.KInfoTable, 0, 1024)
	for _, lst := range block.List { // lst 成分股
		lkey := fmt.Sprintf(REDISKEY_SECURITY_HDAY, lst.NSid)

		klist, err := RedisStore.LRange(lkey, 0, index)
		if err != nil {
			logging.Debug("%v", err.Error())
			return nil, err
		}

		table := &protocol.KInfoTable{} // 成分股历史K线
		for _, info := range klist {
			kinfo := &protocol.KInfo{}
			if err = proto.Unmarshal([]byte(info), kinfo); err != nil {
				logging.Debug("%v", err.Error())
				return nil, err
			}
			if kinfo.NTime < Baseday {
				break
			}
			table.List = append(table.List, kinfo)
		}

		//table
		logging.Info("--除权前--len:%v-----", len(table.List))

		//对K线进行除权操作
		fgs, err := NewXRXD().FactorGroupTotal(rd, lst.Facs, table.List)
		if err != nil {
			logging.Debug("%v", err.Error())
			return nil, err
		}

		rdKTable := &protocol.KInfoTable{
			List: make([]*protocol.KInfo, 0, 1024),
		}
		for _, v := range fgs {
			rdKTable.List = append(rdKTable.List, v.Ls[:]...)
		}
		ktable = append(ktable, rdKTable)
		logging.Info("--除权后--len:%v-----", len(rdKTable.List))
	}

	return &ktable, nil
}

func getKinfoTableAfterEDER(facs []*Factor, table []*protocol.KInfo) {
	xd := NewXRXD()
	xd.ReverseKList(table) //kinfos 翻转  小-->大

}

func getLrangeIndex() int {
	today := utils.Today()

	return (today/10000 - Baseday/10000 + 1) * Workday
}
