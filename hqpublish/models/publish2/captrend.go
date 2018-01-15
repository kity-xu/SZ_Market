package publish2

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"haina.com/share/logging"

	. "haina.com/market/hqpublish/models"

	"haina.com/market/hqpublish/controllers/publish/kline"
	"haina.com/market/hqpublish/models/szdb"
)

// 周期性资金流向返回结构
type ResPeriodCapFlow struct {
	Sid   int32            `json:"sid"`
	Type  int32            `json:"type"`
	Total int32            `json:"total"`
	Begin int32            `json:"begin"`
	Num   int32            `json:"num"`
	CList []*PeriodCapFlow `json:"clist"`
}

// 周期性资金流向(分时、日、周、月、年)
type PeriodCapFlow struct {
	TradeDate  int32   `json:"tradeDate"`  //交易时间
	NetFlowin  float64 `json:"netFlowin"`  //资金净流入
	HugeFlowin float64 `json:"hugeFlowin"` //超大单资金净流入
	BigFlowin  float64 `json:"bigFlowin"`  //大单资金净流入
}

func NewPeriodCapFlow() *PeriodCapFlow {
	return &PeriodCapFlow{}
}

func (p *PeriodCapFlow) GetPeriodCapFlowList(sid, periodID, timeIndex, num, direct int32) (*ResPeriodCapFlow, error) {
	var ret *ResPeriodCapFlow
	var list []*PeriodCapFlow
	switch periodID {
	case 10:
		list = GetCapFlowListMin(sid) //分时
		if len(list) == 0 {
			logging.Error("资金趋势: The src list is null| %d", periodID)
			return nil, DATA_ISNULL
		}
		ret = &ResPeriodCapFlow{
			Sid:   sid,
			Type:  periodID,
			Total: int32(len(list)),
			Begin: list[0].TradeDate,
			Num:   int32(len(list)),
			CList: list,
		}
		return ret, nil
	case 1, 2, 3:
		list = GetCapFlowListODays(sid, periodID) //日、周、月
		if len(list) == 0 {
			logging.Error("资金趋势: The src list is null| %d", periodID)
			return nil, DATA_ISNULL
		}
	default:
		logging.Error("资金趋势: Request param periodID error| %d", periodID)
		return nil, ERROR_REQ_PARAM
	}

	total := int32(len(list))
	if num > total {
		num = total
	}

	if num == 0 { //num==0, 获取全部
		ret = &ResPeriodCapFlow{
			Sid:   sid,
			Type:  periodID,
			Total: total,
			Begin: list[0].TradeDate,
			Num:   total,
			CList: list,
		}
		return ret, nil
	} else { //根据num, 获取部分
		if timeIndex == 0 { //起始日期最新
			var table []*PeriodCapFlow

			lindex := total
			for i := lindex - num; i < lindex; i++ {
				table = append(table, (list)[i])
			}
			ret := &ResPeriodCapFlow{
				Sid:   sid,
				Type:  periodID,
				Total: total,
				Begin: table[0].TradeDate,
				Num:   int32(len(table)),
				CList: table,
			}
			if direct == 1 { //向后
				var sig []*PeriodCapFlow
				sig = append(sig, table[len(table)-1])

				ret = &ResPeriodCapFlow{
					Sid:   sid,
					Type:  periodID,
					Total: total,
					Begin: table[len(table)-1].TradeDate,
					Num:   1,
					CList: sig,
				}
			}
			return ret, nil

		} else { //TimeIndex作为起始日期
			var frontedSwap, palinalSwap []*PeriodCapFlow
			var databuf []*PeriodCapFlow

			for _, v := range list {
				if v.TradeDate <= timeIndex {
					frontedSwap = append(frontedSwap, v)
				}
				if v.TradeDate >= timeIndex {
					palinalSwap = append(palinalSwap, v)
				}
			}

			if direct == 0 { //向前 frontedSwap
				size := len(frontedSwap)
				if size < int(num) {
					databuf = frontedSwap
				} else {
					for i := size - int(num); i < size; i++ {
						databuf = append(databuf, frontedSwap[i])
					}
				}
			} else if direct == 1 { //向后 palinalSwap
				if len(palinalSwap) == 0 { //不加此判断 最新日期向后取，会越界panic
					var table []*PeriodCapFlow

					lindex := total
					for i := lindex - num; i < lindex; i++ {
						table = append(table, list[i])
					}

					var sig []*PeriodCapFlow
					sig = append(sig, table[len(table)-1])

					ret = &ResPeriodCapFlow{
						Sid:   sid,
						Type:  periodID,
						Total: total,
						Begin: table[len(table)-1].TradeDate,
						Num:   1,
						CList: sig,
					}
					return ret, nil
				}

				if int(num) > len(palinalSwap) {
					for i := 0; i < len(palinalSwap); i++ {
						databuf = append(databuf, palinalSwap[i])
					}
				} else {
					for i := 0; i < int(num); i++ {
						databuf = append(databuf, palinalSwap[i])
					}
				}
			} else {
				return nil, ERROR_REQ_PARAM
			}
			ret = &ResPeriodCapFlow{
				Sid:   sid,
				Type:  periodID,
				Total: total,
				Begin: databuf[0].TradeDate,
				Num:   int32(len(databuf)),
				CList: databuf,
			}
			return ret, nil
		}
	}
}

// 当日分时线资金流向
func GetCapFlowListMin(sid int32) []*PeriodCapFlow {
	key := fmt.Sprintf(REDIS_CAP_FLOW_MIN, sid)
	str, err := RedisStore.LRange(key, 0, -1)
	if err != nil {
		logging.Error("分时资金趋势 Redis get error |%v", err.Error())
		return nil
	}

	var list []*PeriodCapFlow
	for _, data := range str {
		ele := &TagTradeScaleStat{}
		if err = binary.Read(bytes.NewBuffer([]byte(data)), binary.LittleEndian, ele); err != nil && err != io.EOF {
			logging.Error("分时资金趋势 binary.read error |%v", err.Error())
			return nil
		}
		p := &PeriodCapFlow{
			TradeDate:  ele.NTime,
			NetFlowin:  float64(ele.LlHugeBuyValue+ele.LlBigBuyValue-ele.LlHugeSellValue-ele.LlBigSellValue) / 10000,
			HugeFlowin: float64(ele.LlHugeBuyValue-ele.LlHugeSellValue) / 10000,
			BigFlowin:  float64(ele.LlBigBuyValue-ele.LlBigSellValue) / 10000,
		}
		list = append(list, p)
	}
	return list
}

// 历史日、周、月、年资金流向
func GetCapFlowListODays(sid, periodID int32) []*PeriodCapFlow {
	return getFlowListFromCache(sid, periodID)
}

// get list from cache
func getFlowListFromCache(sid, periodID int32) []*PeriodCapFlow {
	key := fmt.Sprintf(REDIS_CACHE_CAPITAL_FLOW, sid, periodID)
	var list []*PeriodCapFlow
	data, err := RedisCache.Get(key)
	if len(data) == 0 || err != nil {
		list, err = getFlowListFromSZDB(sid, periodID)
		if err != nil {
			logging.Error("资金趋势：redisCache not fund & SZ db not fund |%v", err)
			return nil
		}
		SortCapFlow(&list)                    // 升序排序
		formartCapTrend(&list, sid, periodID) //加最后一根
		setFlowListToCache(key, &list)        // set cache
	} else {
		if err = json.Unmarshal([]byte(data), &list); err != nil {
			logging.Error("资金趋势：Unmarshal redisCache error |%v", err)
			return nil
		}
	}
	return list
}

// get list from SZDB
func getFlowListFromSZDB(sid int32, periodID int32) ([]*PeriodCapFlow, error) {
	list := make([]*PeriodCapFlow, 0, 512)

	if periodID == 1 {
		pfl, err := szdb.NewSZ_HQ_SECURITYFUNDFLOW().GetHisSecurityFlowFull(sid)
		if len(pfl) == 0 || err != nil {
			return nil, err
		}
		for _, v := range pfl {
			p := &PeriodCapFlow{
				TradeDate:  v.TRADEDATE,
				NetFlowin:  v.HUGEBUYVALUE.Float64 + v.BIGBUYVALUE.Float64 - v.HUGESELLVALUE.Float64 - v.BIGSELLVALUE.Float64,
				HugeFlowin: v.HUGEBUYVALUE.Float64 - v.HUGESELLVALUE.Float64,
				BigFlowin:  v.BIGBUYVALUE.Float64 - v.BIGSELLVALUE.Float64,
			}
			list = append(list, p)
		}

	} else {
		pfl, err := szdb.NewSZ_HQ_SECURITYFUNDFLOW_PERIOD().GetSecurityFundFlowPeriod(sid, periodID)
		if len(pfl) == 0 || err != nil {
			return nil, err
		}

		for _, v := range pfl {
			p := &PeriodCapFlow{
				TradeDate:  int32(v.LASTDATE.Int64), //v.ENTRYDATE
				NetFlowin:  v.HUGEBUYVALUE.Float64 + v.BIGBUYVALUE.Float64 - v.HUGESELLVALUE.Float64 - v.BIGSELLVALUE.Float64,
				HugeFlowin: v.HUGEBUYVALUE.Float64 - v.HUGESELLVALUE.Float64,
				BigFlowin:  v.BIGBUYVALUE.Float64 - v.BIGSELLVALUE.Float64,
			}
			list = append(list, p)
		}

	}
	return list, nil
}

// set list to cache
func setFlowListToCache(key string, list *[]*PeriodCapFlow) error {
	bys, err := json.Marshal(list)
	if err != nil {
		logging.Error("资金趋势：Marshal redisCache error |%v", err)
		return err
	}

	if err = RedisCache.Setex(key, TTL.Day, bys); err != nil {
		logging.Error("资金趋势：Set redisCache error |%v", err)
		return err
	}
	return nil
}

//------------------------------------------------------------------------------------------------------//

// 资金流向升序排序
func SortCapFlow(list *[]*PeriodCapFlow) {
	if len(*list) == 0 {
		return
	}
	sort.Sort(sort.Reverse(clist(*list)))
}

type clist []*PeriodCapFlow

func (this clist) Len() int {
	return len(this)
}

func (this clist) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this clist) Less(i, j int) bool {
	return this[j].TradeDate < this[i].TradeDate
}

//----------------------------------------------------------maybeAdd------------------------//
// 处理最新一天的资金趋势
func formartCapTrend(funds *[]*PeriodCapFlow, sid int32, ntype int32) *[]*PeriodCapFlow {
	kline.InitMarketTradeDate()
	last := (*funds)[len(*funds)-1]
	if last.TradeDate == kline.Trade_100 { // 库里更新取库里，不做其他操作
		return nil
	}
	captoday := &PeriodCapFlow{
		TradeDate: kline.Trade_100,
	}

	cap, _ := CapFlowToday(sid)
	switch ntype {
	case 1:
		captoday.NetFlowin = float64(cap.LlValueOfInFlow) / 10000
		captoday.HugeFlowin = float64(cap.LlHugeBuyValue) / 10000
		captoday.BigFlowin = float64(cap.LlBigBuyValue) / 10000
		*funds = append(*funds, captoday)
	case 2:
		if last.TradeDate < kline.Trade_100 {
			b1, _ := DateAdd(last.TradeDate) //找到该日期所在周日的那天
			b2, _ := DateAdd(captoday.TradeDate)

			if !b1.Equal(b2) { //不同属一周（周一）新建
				captoday.NetFlowin = float64(cap.LlValueOfInFlow) / 10000
				captoday.HugeFlowin = float64(cap.LlHugeBuyValue) / 10000
				captoday.BigFlowin = float64(cap.LlBigBuyValue) / 10000
				*funds = append(*funds, captoday)
			} else { //同属一周，更新最后一根
				last.TradeDate = captoday.TradeDate
				last.NetFlowin += float64(cap.LlValueOfInFlow) / 10000
				last.HugeFlowin += float64(cap.LlHugeBuyValue) / 10000
				last.BigFlowin += float64(cap.LlBigBuyValue) / 10000
				(*funds)[len(*funds)-1] = last
			}
		}
	case 3:
		if last.TradeDate < kline.Trade_100 {
			if last.TradeDate/100 != captoday.TradeDate/100 { //不同月
				captoday.NetFlowin = float64(cap.LlValueOfInFlow) / 10000
				captoday.HugeFlowin = float64(cap.LlHugeBuyValue) / 10000
				captoday.BigFlowin = float64(cap.LlBigBuyValue) / 10000
				*funds = append(*funds, captoday)
			} else {
				last.TradeDate = captoday.TradeDate
				last.NetFlowin += float64(cap.LlValueOfInFlow) / 10000
				last.HugeFlowin += float64(cap.LlHugeBuyValue) / 10000
				last.BigFlowin += float64(cap.LlBigBuyValue) / 10000
				(*funds)[len(*funds)-1] = last
			}
		}
	}
	return nil
}

// 当日资金流向
func CapFlowToday(sid int32) (*TagTradeScaleStat, error) {
	keyd := fmt.Sprintf("hq:trade:day:%d", sid)
	data, err := RedisStore.GetBytes(keyd)
	if err != nil {
		logging.Error("%v", err.Error())
		return nil, err
	}
	var ele = TagTradeScaleStat{}
	if err = binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &ele); err != nil && err != io.EOF {
		logging.Error("%v", err.Error())
		return nil, err
	}
	return &ele, nil
}
