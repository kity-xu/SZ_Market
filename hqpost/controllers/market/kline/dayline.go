package kline

import (
	pbk "ProtocolBuffer/format/kline"
	"fmt"
	"io/ioutil"
	"strings"

	"haina.com/market/hqpost/models"

	"haina.com/market/hqpost/config"
	"haina.com/market/hqpost/models/redis_minline"

	"haina.com/share/lib"
	"haina.com/share/store/redis"

	"github.com/golang/protobuf/proto"
	"haina.com/share/logging"
)

func NewSecurityKLine(sids *[]int32, cg *config.AppConfig) *Security {
	cfg = cg
	return &Security{
		sids: sids,
	}
}

func (this *Security) DayLine() {
	var seList []SingleSecurity

	/******************************沪深所有股票*************************************/
	for _, sid := range *this.sids {
		var filename string //文件名
		var exchange string //股票交易所

		var issrc bool = false //判断是否需要去读源文件

		//PB
		var klist pbk.KInfoTable

		//History of Single-Security
		var sigList SingleSecurity
		var date []int32
		dmap := make(map[int32]pbk.KInfo)

		if sid/100000000 == 1 { //ascii 字符
			exchange = SH
		} else if sid/100000000 == 2 {
			exchange = SZ
		} else {
			logging.Error("%s", "Invalid file name...")
			return
		}

		/**********************************************filename*******************************************************/
		//这里要做一个逻辑判断
		// 1. 先判断haina历史文件是否存在，不存在去读源文件做第一次生成
		hnfile := fmt.Sprintf("%s%s%d/%s", cfg.File.Path, exchange, sid, cfg.File.Day) //haina文件store路劲
		if !lib.IsFileExist(hnfile) {                                                  //haina store dk.dat不存在
			hnindex := fmt.Sprintf("%s%s%d/%s", cfg.File.Path, exchange, sid, cfg.File.Index)
			if !lib.IsFileExist(hnindex) { //haina store index.dat不存在
				issrc = true //说明需要去读源文件
			} else {
				filename = hnindex
			}

		} else {
			filename = hnfile
		}

		if issrc { //读源文件的逻辑操作
			srcfile := fmt.Sprintf("%s%s%d/%s", cfg.File.Finpath, exchange, sid, cfg.File.Finday) //src文件store路劲
			if !lib.IsFileExist(srcfile) {
				srcindex := fmt.Sprintf("%s%s%d/%s", cfg.File.Finpath, exchange, sid, cfg.File.Findex)
				if !lib.IsFileExist(srcindex) {
					logging.Error("Cannot find file(haina filestore && finchina filestore)...")
					logging.Error("SID:%v源文件不存在", sid)
					filename = ""
					continue
				} else {
					filename = srcindex
				}
			} else {
				filename = srcfile
			}
		}

		/**********************************************filename*******************************************************/

		//读文件
		fd, err := ioutil.ReadFile(filename)
		if err != nil {
			logging.Error("%v", err.Error())
			return
		}
		//解PB
		if err = proto.Unmarshal(fd, &klist); err != nil {
			logging.Error("%v", err.Error())
			return
		}

		//map SingleSecurity结构
		for _, v := range klist.List {
			date = append(date, v.NTime)
			dmap[v.NTime] = *v
		}

		//得到今天的k线
		tm := getDateToday()

		//昨收价 lastPx
		models.GetASCStruct(&klist.List) //按时间升序排序
		today, e := GetTodayDayLine(sid, klist.List[len(klist.List)-1].NLastPx)
		if e == nil && today != nil {
			klist.List = append(klist.List, today) //追加
			dmap[tm] = *today
		}

		//转PB
		data, err := proto.Marshal(&klist)
		if err != nil {
			logging.Error("Encode protocbuf of day Line error...%v", err.Error())
			return
		}

		/*******************入文件******************************/
		ss := strings.Split(filename, "/")
		if e := KlineWriteFile(sid, ss[len(ss)-1], &data); e != nil {
			logging.Error("%v", err.Error())
			return
		}

		//入redis
		key := fmt.Sprintf(REDISKEY_SECURITY_HDAY, sid)
		if err := redis.Set(key, data); err != nil {
			logging.Fatal("%v", err)
		}

		sigList.Sid = sid
		sigList.Date = date
		sigList.SigStock = dmap

		seList = append(seList, sigList)

		/*-------------------------------------end------------------------------------*/

		//logging.Debug("The historical data of each stock number:%v", count)
	}

	this.list.Securitys = &seList
	/*-----------------------------------------end----------------------------------*/
}

//获取今天分钟线生成的日线
func GetTodayDayLine(sid int32, lastPx int32) (*pbk.KInfo, error) {
	min, err := redis_minline.NewMinKLine(REDISKEY_SECURITY_MIN).GetMinKLineToday(sid)
	if min == nil || err != nil {
		//logging.Error("%v", err.Error())
		return nil, err
	}
	var tmp pbk.KInfo //pb类型

	var (
		i          int = 0
		AvgPxTotal uint32
	)

	models.GetASCStruct(min) //按时间升序排序
	for _, v := range *min {
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
	tmp.NSID = sid
	tmp.NTime = getDateToday()      //时间
	tmp.NOpenPx = (*min)[0].NOpenPx //开盘价
	tmp.NPreCPx = lastPx            //昨收价

	tmp.NLastPx = (*min)[len(*min)-1].NLastPx //最新价
	tmp.NAvgPx = AvgPxTotal / uint32(i+1)     //平均价

	return &tmp, nil
}
