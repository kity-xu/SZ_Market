package kline

import (
	"ProtocolBuffer/projects/hqpost/go/protocol"
	"fmt"
	"strings"

	"haina.com/market/hqpost/models/filestore"
	"haina.com/market/hqpost/models/lib"
	"haina.com/market/hqpost/models/redistore"

	"haina.com/market/hqpost/config"
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
	rstore := redistore.NewHKLine(REDISKEY_SECURITY_HDAY)

	/******************************沪深所有股票*************************************/
	for _, sid := range *this.sids {
		var (
			e, err   error
			filename string         //文件名
			exchange string         //股票交易所
			issrc    bool   = false //判断是否需要去读源文件
			klist           = &protocol.KInfoTable{}
			sigList  SingleSecurity
			date     []int32
		)

		dmap := make(map[int32]protocol.KInfo)

		if sid/100000000 == 1 { //ascii 字符
			exchange = SH
		} else if sid/100000000 == 2 {
			exchange = SZ
		} else {
			logging.Error("%s", "Invalid file name...")
			return
		}

		/**********************************************filename*******************************************************/
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

		if issrc { //读源文件的逻辑操作(表示FileStore不存在)
			srcfile := fmt.Sprintf("%s%s%d/%s", cfg.File.Finpath, exchange, sid, cfg.File.Finday) //src文件store路劲
			if !lib.IsFileExist(srcfile) {
				srcindex := fmt.Sprintf("%s%s%d/%s", cfg.File.Finpath, exchange, sid, cfg.File.Findex)
				if !lib.IsFileExist(srcindex) { //新增的K线（个股或指数新上市）
					tag := redistore.IsNSidIndex(sid)
					if tag == 1 {
						filename = hnfile
						lib.CreateFile(filename)

						//创建周、月、年文件路劲，用于新增股票（或指数）的追加
						lib.CreateFile(fmt.Sprintf("%s%s%d/%s", cfg.File.Path, exchange, sid, cfg.File.Week))
						lib.CreateFile(fmt.Sprintf("%s%s%d/%s", cfg.File.Path, exchange, sid, cfg.File.Month))
						lib.CreateFile(fmt.Sprintf("%s%s%d/%s", cfg.File.Path, exchange, sid, cfg.File.Year))
						goto LABEL
					} else if tag == 7 {
						filename = fmt.Sprintf("%s%s%d/%s", cfg.File.Path, exchange, sid, cfg.File.Index)
						lib.CreateFile(filename)

						//创建周、月、年文件路劲，用于新增股票（或指数）的追加
						lib.CreateFile(fmt.Sprintf("%s%s%d/%s", cfg.File.Path, exchange, sid, cfg.File.Week))
						lib.CreateFile(fmt.Sprintf("%s%s%d/%s", cfg.File.Path, exchange, sid, cfg.File.Month))
						lib.CreateFile(fmt.Sprintf("%s%s%d/%s", cfg.File.Path, exchange, sid, cfg.File.Year))
						goto LABEL
					} else {
						filename = ""
						continue
					}
				} else {
					filename = srcindex
				}
			} else {
				filename = srcfile
			}

			//解析源文件数据
			klist, err = filestore.ReadSrcFileStore(filename)
			if err != nil {
				continue
			}

			//从源搬到haina FileStore
			ss := strings.Split(filename, "/")
			filename, _ = filestore.CheckFileSoteDir(sid, cfg.File.Path, ss[len(ss)-1])
			filestore.WiteHainaFileStore(filename, klist)

			//redis做第一次生成
			for _, v := range klist.List {
				if err := rstore.LPushHisKLine(sid, v); err != nil {
					logging.Error("%v", err.Error())
					return
				}
			}

		} else {
			//读haina FileStore
			klist, e = filestore.ReadHainaFileStore(filename)
			if e != nil {
				logging.Error("%v", e.Error())
				return
			}
		}
		if len(klist.List) < 1 {
			logging.Error("SID:%v---No historical data...", sid)
			continue
		}

		//map SingleSecurity结构
		for _, v := range klist.List {
			date = append(date, v.NTime)
			dmap[v.NTime] = *v
		}

		lib.GetASCStruct(&klist.List) //按时间升序排序

	LABEL:
		today, e := GetTodayDayLine(sid) //得到今天的k线
		if e == nil && today != nil {    //获取当天数据没毛病
			sigList.today = today

			//追加到文件
			if err := filestore.AppendFile(filename, today); err != nil {
				logging.Error("%v", err.Error())
				continue
			}

			//追加redis
			if err := rstore.AppendTodayLine(sid, today); err != nil {
				logging.Error("%v", err.Error())
				continue
			}
		} else {
			//logging.Debug("获取%v当天K线信息失败...", sid)
		}
		sigList.Sid = sid
		sigList.Date = date
		sigList.SigStock = dmap

		seList = append(seList, sigList)
	}
	this.list.Securitys = &seList
}

//获取今天分钟线生成的日线
func GetTodayDayLine(sid int32) (*protocol.KInfo, error) {
	min, err := redistore.NewMinKLine(REDISKEY_SECURITY_MIN).GetMinKLineToday(sid)
	if err != nil {
		logging.Debug("%v", err.Error())
		return nil, err
	}
	var tmp protocol.KInfo //pb类型

	var (
		i          int = 0
		AvgPxTotal uint32
	)

	lib.GetASCStruct(min) //按时间升序排序
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
	tmp.NTime = filestore.GetDateToday()      //时间
	tmp.NOpenPx = (*min)[0].NOpenPx           //开盘价
	tmp.NPreCPx = (*min)[len(*min)-1].NPreCPx //昨收价
	tmp.NLastPx = (*min)[len(*min)-1].NLastPx //最新价
	tmp.NAvgPx = AvgPxTotal / uint32(i+1)     //平均价

	return &tmp, nil
}
