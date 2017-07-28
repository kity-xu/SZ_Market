package servers

import (
	"ProtocolBuffer/projects/hqinit/go/protocol"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/proto"
	"haina.com/market/hqinit/models/fcmysql"
	"haina.com/share/logging"
	. "haina.com/share/models"
)

type StockBlockRedis struct {
	Model `db:"-"`
}

func NewStockBlockRedis() *StockBlockRedis {
	return &StockBlockRedis{
		Model: Model{},
	}
}

var (
	ALL                     string //全部
	DISTRICT                string //地区
	CONCEPT                 string //概念
	INDUSTRY                string //行业
	REDISKEY_BLOCK_CLASSIFY = "hq:init:bk:%d"
	REDISKEY_BLOCK_BOARD    = "hq:init:bk:%d:%d"
	REDISKEY_STOCK_BLOCK    = "hq:init:bk:stock:%d"
)

func (this *StockBlockRedis) Block() {
	start := time.Now()

	logging.Info("begin ...")
	c, errr := redis.Dial("tcp", "127.0.0.1:61380")
	//c, errr := redis.Dial("tcp", "47.94.16.69:61380")
	c.Send("AUTH", "8dc40c2c4598ae5a")
	if errr != nil {
		logging.Info("redis conn error %v", errr)
	}

	boar1j, err := fcmysql.NewTQ_COMP_BOARDMAP().GetBoardmapRedis()
	if err != nil {
		logging.Debug("mysql 1j", err)
	}

	DISTRICT = strconv.Itoa(int(protocol.REDIS_BLOCK_CLASSIFY_District))
	CONCEPT = strconv.Itoa(int(protocol.REDIS_BLOCK_CLASSIFY_Concept))
	INDUSTRY = strconv.Itoa(int(protocol.REDIS_BLOCK_CLASSIFY_Industry))
	ALL = strconv.Itoa(int(protocol.REDIS_BLOCK_CLASSIFY_All))

	disMap := make(map[int32][]*protocol.Element)
	conMap := make(map[int32][]*protocol.Element)
	indusMap := make(map[int32][]*protocol.Element)

	for _, v := range boar1j { //KeyCode 去掉"CN"
		switch v.BOARDCODE.String {
		case DISTRICT: //地区
			ele := &protocol.Element{
				NSid:    stringToInt32((v.SECODE.String)),
				Keyname: v.KEYNAME.String,
			}
			index := stringToInt32("81" + ((v.KEYCODE.String)[2:]))

			disMap[index] = append(disMap[index], ele)
		case CONCEPT: //概念
			ele := &protocol.Element{
				NSid:    stringToInt32(v.SECODE.String),
				Keyname: v.KEYNAME.String,
			}
			index := stringToInt32("81" + v.KEYCODE.String)
			conMap[index] = append(conMap[index], ele)
		case INDUSTRY: //行业
			ele := &protocol.Element{
				NSid:    stringToInt32(v.SECODE.String),
				Keyname: v.KEYNAME.String,
			}
			index := stringToInt32("81" + v.KEYCODE.String)
			indusMap[index] = append(indusMap[index], ele)
		default:
		}
	}
	//-------------------------------------------------------------------------------//

	var boards1 = &protocol.BlockList{}
	for bid, element := range disMap { //key,value: 某个地区下的成份股

		var secstr string
		for _, v := range element {
			secstr += "'" + int32Tostring(v.NSid) + "',"
		}

		secstr = strings.TrimRight(secstr, ",")
		//logging.Debug("secstr--%v", secstr)

		//查数据库
		stock, err := fcmysql.NewFcSecuNameTab().GetComCodeList(secstr)
		//logging.Debug("--len stock-%v", len(stock))
		if err != nil {
			logging.Error("%v", err.Error())
			return
		}
		var elms = &protocol.ElementList{}
		var sid int32
		for _, v := range stock {
			if strings.EqualFold(v.EXCHANGE.String, "001002") { //上海
				sid = 100*1000000 + stringToInt32(v.SYMBOL.String)
			} else if strings.EqualFold(v.EXCHANGE.String, "001003") { //深圳
				sid = 200*1000000 + stringToInt32(v.SYMBOL.String)
			} else {
				sid = stringToInt32(v.SYMBOL.String)
			}

			elm := &protocol.Element{
				NSid:    sid,
				Keyname: v.SENAME.String,
			}
			//logging.Debug("-----1102:%v", elm)

			elms.List = append(elms.List, elm)
		}

		//以板块分类的成份股
		data, err := proto.Marshal(elms)
		if err != nil {
			logging.Error("%v", err.Error())
			return
		}
		key := fmt.Sprintf(REDISKEY_BLOCK_BOARD, protocol.REDIS_BLOCK_CLASSIFY_District, bid)
		if _, err = c.Do("SET", key, data); err != nil {
			logging.Error("%v", err.Error())
			return
		}

		//所有板块下成份股
		keyall := fmt.Sprintf(REDISKEY_BLOCK_BOARD, protocol.REDIS_BLOCK_CLASSIFY_All, bid)
		if _, err = c.Do("SET", keyall, data); err != nil {
			logging.Error("%v", err.Error())
			return
		}

		//以类型分类的板块
		board := &protocol.Block{
			SetID:   bid,
			SetName: disMap[bid][0].Keyname,
		}

		boards1.List = append(boards1.List, board)
	}
	data1, err := proto.Marshal(boards1)
	if err != nil {
		logging.Error("%v", err.Error())
		return
	}

	key1 := fmt.Sprintf(REDISKEY_BLOCK_CLASSIFY, protocol.REDIS_BLOCK_CLASSIFY_District)
	if _, err = c.Do("SET", key1, data1); err != nil {
		logging.Error("%v", err.Error())
		return
	}

	//-------------------------------------------------------------------------------//
	var boards2 = &protocol.BlockList{}
	for bid, element := range conMap { //key,value: 某个概念下的成份股
		var secstr string
		for _, v := range element {
			secstr += "'" + int32Tostring(v.NSid) + "',"
		}

		secstr = strings.TrimRight(secstr, ",")

		//查数据库
		stock, err := fcmysql.NewFcSecuNameTab().GetComCodeList(secstr)
		if err != nil {
			logging.Error("%v", err.Error())
			return
		}

		var elms = &protocol.ElementList{}
		var sid int32
		for _, v := range stock {
			if strings.EqualFold(v.EXCHANGE.String, "001002") { //上海
				sid = 100*1000000 + stringToInt32(v.SYMBOL.String)
			} else if strings.EqualFold(v.EXCHANGE.String, "001003") { //深圳
				sid = 200*1000000 + stringToInt32(v.SYMBOL.String)
			} else {
				sid = stringToInt32(v.SYMBOL.String)
			}

			elm := &protocol.Element{
				NSid:    sid,
				Keyname: v.SENAME.String,
			}
			elms.List = append(elms.List, elm)
		}

		//以板块分类的成份股
		data, err := proto.Marshal(elms)
		if err != nil {
			logging.Error("%v", err.Error())
			return
		}
		key := fmt.Sprintf(REDISKEY_BLOCK_BOARD, protocol.REDIS_BLOCK_CLASSIFY_Concept, bid)
		if _, err = c.Do("SET", key, data); err != nil {
			logging.Error("%v", err.Error())
			return
		}
		//所有板块下成份股
		keyall := fmt.Sprintf(REDISKEY_BLOCK_BOARD, protocol.REDIS_BLOCK_CLASSIFY_All, bid)
		if _, err = c.Do("SET", keyall, data); err != nil {
			logging.Error("%v", err.Error())
			return
		}

		//以类型分类的板块
		board := &protocol.Block{
			SetID:   bid,
			SetName: conMap[bid][0].Keyname,
		}

		boards2.List = append(boards2.List, board)
	}

	data2, err := proto.Marshal(boards2)
	if err != nil {
		logging.Error("%v", err.Error())
		return
	}

	key2 := fmt.Sprintf(REDISKEY_BLOCK_CLASSIFY, protocol.REDIS_BLOCK_CLASSIFY_Concept)

	if _, err = c.Do("SET", key2, data2); err != nil {
		logging.Error("%v", err.Error())
		return
	}
	//-------------------------------------------------------------------------------//
	var boards3 = &protocol.BlockList{}
	var sid int32

	for bid, element := range indusMap { //key,value: 某个行业下的成份股

		var secstr string
		for _, v := range element {
			secstr += "'" + int32Tostring(v.NSid) + "',"
		}

		secstr = strings.TrimRight(secstr, ",")

		//查数据库
		stock, err := fcmysql.NewFcSecuNameTab().GetComCodeList(secstr)
		if err != nil {
			logging.Error("%v", err.Error())
			return
		}

		var elms = &protocol.ElementList{}
		for _, v := range stock {
			if strings.EqualFold(v.EXCHANGE.String, "001002") { //上海
				sid = 100*1000000 + stringToInt32(v.SYMBOL.String)
			} else if strings.EqualFold(v.EXCHANGE.String, "001003") { //深圳
				sid = 200*1000000 + stringToInt32(v.SYMBOL.String)
			} else {
				sid = stringToInt32(v.SYMBOL.String)
			}
			elm := &protocol.Element{
				NSid:    sid,
				Keyname: v.SENAME.String,
			}
			elms.List = append(elms.List, elm)

		}

		//以板块分类的成份股

		data, err := proto.Marshal(elms)
		if err != nil {
			logging.Error("%v", err.Error())
			return
		}
		key := fmt.Sprintf(REDISKEY_BLOCK_BOARD, protocol.REDIS_BLOCK_CLASSIFY_Industry, bid)
		if _, err = c.Do("SET", key, data); err != nil {
			logging.Error("%v", err.Error())
			return
		}
		//所有板块下成份股
		keyall := fmt.Sprintf(REDISKEY_BLOCK_BOARD, protocol.REDIS_BLOCK_CLASSIFY_All, bid)
		if _, err = c.Do("SET", keyall, data); err != nil {
			logging.Error("%v", err.Error())
			return
		}

		//以类型分类的板块
		board := &protocol.Block{
			SetID:   bid,
			SetName: indusMap[bid][0].Keyname,
		}

		boards3.List = append(boards3.List, board)

	}

	data3, err := proto.Marshal(boards3)
	if err != nil {
		logging.Error("%v", err.Error())
		return
	}

	key3 := fmt.Sprintf(REDISKEY_BLOCK_CLASSIFY, protocol.REDIS_BLOCK_CLASSIFY_Industry)
	if _, err = c.Do("SET", key3, data3); err != nil {
		logging.Error("%v", err.Error())
		return
	}

	//-------------------------------------------------------------------------------//
	boards1.List = append(boards1.List, boards2.List...)
	boards1.List = append(boards1.List, boards3.List...)

	data4, err := proto.Marshal(boards1)
	if err != nil {
		logging.Error("%v", err.Error())
		return
	}

	key4 := fmt.Sprintf(REDISKEY_BLOCK_CLASSIFY, protocol.REDIS_BLOCK_CLASSIFY_All)
	if _, err = c.Do("SET", key4, data4); err != nil {
		logging.Error("%v", err.Error())
		return
	}

	//-------------------------------------------------------处理证券查板块
	stockL, err := fcmysql.NewFcSecuNameTab().GetSecuNmList()
	if err != nil {
		logging.Error("select stcode table error:%v", err)
	}
	// 遍历所有股票
	for _, stite := range stockL {
		// 根据secode 查询所属板块
		blocks, err := fcmysql.NewTQ_COMP_BOARDMAP().GetBoadBySID(stite.SECODE.String)
		if err != nil {
			logging.Error("select TQ_COMP_BOARDMAP table error:%v", err)
		}
		var bls protocol.BlockList
		for _, blite := range blocks {
			var bk protocol.Block
			kcode := blite.KEYCODE.String
			if kcode[:1] == "C" && kcode[1:2] == "N" {
				kcode = kcode[2:]
			}
			kd, err := strconv.Atoi("81" + kcode)
			if err != nil {
				logging.Error("kcode  type conversion error:%v", err)
			}
			bk.SetID = int32(kd)
			bls.List = append(bls.List, &bk)
		}
		var sym = ""
		if stite.EXCHANGE.String == "001002" {
			sym = "100" + stite.SYMBOL.String
		} else if stite.EXCHANGE.String == "001003" {
			sym = "200" + stite.SYMBOL.String
		}
		symb, err := strconv.Atoi(sym)
		if err != nil {
			logging.Error("sym typeconversion error:%v", err)
		}
		bkkey := fmt.Sprintf(REDISKEY_STOCK_BLOCK, symb)

		bkdata, err := proto.Marshal(&bls)

		if err != nil {
			return
		}
		if _, err = c.Do("SET", bkkey, bkdata); err != nil {
			logging.Error("%v", err.Error())
			return
		}
	}
	//-------------------------------------------------------处理证券查板块
	end := time.Now()
	logging.Info("Update Kline historical data successed, and running time:%v", end.Sub(start))
}

func stringToInt32(str string) int32 {
	dd, err := strconv.Atoi(str)
	if err != nil {
		logging.Error("stringToInt32 error...")
	}
	return int32(dd)
}

func int32Tostring(dd int32) string {
	return strconv.Itoa(int(dd))
}
