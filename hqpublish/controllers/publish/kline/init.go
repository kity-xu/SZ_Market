package kline

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"
	"strconv"
	"sync"
	"time"

	"haina.com/market/hqpublish/models/publish"
	"haina.com/share/logging"
)

//沪深交易市场当前交易日
var (
	Trade_100 int32 = 0
	Trade_200 int32 = 0
)

var once sync.Once

func init() {
	once.Do(Timer)
}

func Timer() {
	go func() {
		t1 := time.NewTimer(time.Minute * 5)
		for {
			//f()
			now := time.Now()
			ntime, _ := strconv.Atoi(now.Format("20060102"))

			if 800 < ntime%10000 && ntime%10000 < 1000 { // 8:00 < ntime < 9:00
				<-t1.C
				t1.Reset(time.Minute * 5)
				Trade_100 = 0
				Trade_200 = 0
			}
			// 计算下一个零点
			//now := time.Now()
			//next := now.Add(time.Minute * 30)
			//next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), 0, 0, 0, next.Location())
			//t := time.NewTimer(next.Sub(now))
			//<-t.C
		}
	}()
}

func InitMarketTradeDate() {
	initMarketTradeDate()
}

func initMarketTradeDate() {
	if Trade_100 != 0 && Trade_200 != 0 {
		return
	}

	mlist, err := getMarketStatus()
	if err != nil {
		logging.Error("%v", err)
		return
	}
	for _, v := range mlist.MSList {
		if v.NTradeDate != 0 {
			Trade_200, Trade_100 = v.NTradeDate, v.NTradeDate
		}
	}
	logging.Info("---------------Trade_100:%v", Trade_100)
	logging.Info("---------------Trade_100:%v", Trade_200)
}

//市场状态获取当前交易日
func getMarketStatus() (*protocol.PayloadMarketStatus, error) {
	var req = protocol.RequestMarketStatus{
		Num:          2,
		MarketIDList: []int32{100000000, 200000000},
	}

	res, err := publish.NewMarketStatus().GetPayloadObj(&req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
