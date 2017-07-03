// 分价成交
package publish

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"sort"
	"strconv"

	"haina.com/share/logging"
	. "haina.com/share/models"

	. "haina.com/market/hqpublish/models"
)

type TradePriceRecordM struct {
	Model `db:"-"`
}

func NewTradePriceRecordM() *TradePriceRecordM {
	return &TradePriceRecordM{
		Model: Model{
			CacheKey: REDISKEY_TRADE_PRICE,
		},
	}
}

// 分价列表
func (this *TradePriceRecordM) GetTradePriceRecordL(req *protocol.RequestTradePriceR) (*protocol.PayloadTradePriceR, error) {

	var tpr protocol.PayloadTradePriceR

	// 获取所有分价列表
	var err error
	key := fmt.Sprintf(this.CacheKey, req.SID)
	rul, err := RedisStore.Hgetall(key)

	for i, _ := range rul {
		tpi := &protocol.TradePriceI{}
		var tpif protocol.TradePriceRecord
		bufer := bytes.NewBuffer([]byte(rul[i]))
		if err := binary.Read(bufer, binary.LittleEndian, &tpif); err != nil && err != io.EOF {
			logging.Error("%v", err)
			return nil, err
		}
		res, err := strconv.Atoi(i)
		if err != nil {
			logging.Info("类型转换 error%v", err)
		}
		tpi.SNInfo = &tpif
		tpi.Pricew = uint32(res)
		tpr.TPRList = append(tpr.TPRList, tpi)
	}
	tpr.SID = req.SID

	getSECStruct(&tpr.TPRList)

	return &tpr, err
}

// 根据价格排序
type tradePrice []*protocol.TradePriceI

func (this tradePrice) Len() int {
	return len(this)
}

func (this tradePrice) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this tradePrice) Less(i, j int) bool {
	return this[j].Pricew < this[i].Pricew
}

//升序
func getASCStruct(a *[]*protocol.TradePriceI) {
	sort.Sort(sort.Reverse(tradePrice(*a)))
}

//降序
func getSECStruct(a *[]*protocol.TradePriceI) {
	sort.Sort(tradePrice(*a))
}
