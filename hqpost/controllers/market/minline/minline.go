package minline

import (
	"ProtocolBuffer/projects/hqpost/go/protocol"

	"haina.com/market/hqpost/config"
	"haina.com/market/hqpost/models"
	"haina.com/market/hqpost/models/lib"
	"haina.com/market/hqpost/models/redistore"
	"haina.com/share/logging"
)

func NewMinKline(sids *[]int32, cg *config.AppConfig) *MinKline {
	cfg = cg
	return &MinKline{
		sids: sids,
	}
}

var (
	rstore1 *redistore.HMinKLine
	mstore  *redistore.MinKLine
)

//所有股票当天的分钟线数据
func (this *MinKline) HMinLine_1() {
	rstore1 = redistore.NewHMinKLine(REDISKEY_SECURITY_HMIN)
	mstore = redistore.NewMinKLine(REDISKEY_SECURITY_MIN)
	this.initBasicMinData()
}

//初始化分钟线基本数据
func (this *MinKline) initBasicMinData() {
	var All []*SingleMin

	for _, sid := range *this.sids {
		min, err := this.getBasicMinDataToday(sid)
		if err != nil {
			if err == models.ERROR_REDIS_LIST_NULL { //双方nsid表不同导致的错误（有一方没有更新）
				continue
			} else { //其他严重错误
				logging.Error("%v", err)
				return
			}
		}
		All = append(All, min)
	}
	this.list.All = &All
}

//单个股票当天分钟线基本数据
func (this *MinKline) getBasicMinDataToday(sid int32) (*SingleMin, error) {
	day := &SingleMin{
		Time: make([]int32, 0), //单个股票的历史日期
		Min:  make(map[int32]protocol.KInfo),
	}
	dmin, err := mstore.GetMinKLineToday(sid)
	if err != nil {
		return nil, err
	}

	for _, v := range *dmin {
		day.Min[v.NTime] = *v
		day.Time = append(day.Time, v.NTime)
	}

	/*********************1分钟线操作**************************************/

	this.mergeMin(sid, rstore1, dmin)
	/***********************1分钟历史线操作OVER*******************************/

	day.Sid = sid

	lib.GetASCIntArray(day.Time) //升序排序time
	generateMinLineTimes(day)

	return day, nil
}

//当天分钟线并入历史
func (this *MinKline) mergeMin(sid int32, rs *redistore.HMinKLine, dmin *[]*protocol.KInfo) {
	if rs.Ktype != REDISKEY_SECURITY_HMIN { //如果不是历史1分钟线的话，进行redis操作（也就是说minline_1不进redis）
		kinfo := protocol.HMinLineDay{
			Date: GetDateToday(),
			List: *dmin,
		}

		if err := rs.LPushHMinKLine(sid, &kinfo); err != nil {
			return
		}
	}

	//文件
	switch rs.Ktype {
	case REDISKEY_SECURITY_HMIN:
		if e := AppendFile(sid, cfg.File.Min, dmin); e != nil {
			logging.Error("%v", e.Error())
			return
		}
		break
	case REDISKEY_SECURITY_HMIN5:
		if e := AppendFile(sid, cfg.File.Min5, dmin); e != nil {
			logging.Error("%v", e.Error())
			return
		}
		break
	case REDISKEY_SECURITY_HMIN15:
		if e := AppendFile(sid, cfg.File.Min15, dmin); e != nil {
			logging.Error("%v", e.Error())
			return
		}
		break
	case REDISKEY_SECURITY_HMIN30:
		if e := AppendFile(sid, cfg.File.Min30, dmin); e != nil {
			logging.Error("%v", e.Error())
			return
		}
		break
	case REDISKEY_SECURITY_HMIN60:
		if e := AppendFile(sid, cfg.File.Min60, dmin); e != nil {
			logging.Error("%v", e.Error())
			return
		}
		break
	default:
		return
	}

}

//生成分钟线时间（5、15、30、60）[][]into2
func generateMinLineTimes(day *SingleMin) {
	var minbuf_5, minbuf_15, minbuf_30, minbuf_60 []int32
	var time_5, time_15, time_30, time_60 [][]int32

	st_05 := MIN_START + 5  // 0935
	st_15 := MIN_START + 15 // 0945
	st_30 := 1000           // 930 + 30
	st_60 := 1030           // 1030
	for i, min := range day.Time {
		length := len(day.Time)

		//5
		if min%10000 <= int32(st_05) {
			minbuf_5 = append(minbuf_5, min)
			if i == length-1 { //别遗漏最后一条
				time_5 = append(time_5, minbuf_5)
			}
		} else { //生成了一个n分钟
			time_5 = append(time_5, minbuf_5)

			minbuf_5 = nil //缓冲置空（清除上一次的缓冲数据）

			minbuf_5 = append(minbuf_5, min) //本次加进缓冲
			st_05 = st_05 + 5                //更新临界时间
			if st_05%100 == 60 {
				st_05 = (st_05/100 + 1) * 100
			}
			if 1130 < st_05 && st_05 < 1300 {
				st_05 = 1300 + 5
			}
		}

		//15
		if min%10000 <= int32(st_15) {
			minbuf_15 = append(minbuf_15, min)
			if i == length-1 { //别遗漏最后一条
				time_15 = append(time_15, minbuf_15)
			}
		} else {
			time_15 = append(time_15, minbuf_15)
			minbuf_15 = nil
			minbuf_15 = append(minbuf_15, min)
			st_15 = st_15 + 15 //更新临界时间

			if st_15%100 == 60 {
				st_15 = (st_15/100 + 1) * 100
			}
			if 1130 < st_15 && st_15 < 1300 {
				st_15 = 1300 + 15
			}
		}

		//30
		if min%10000 <= int32(st_30) {
			minbuf_30 = append(minbuf_30, min)
			if i == length-1 { //别遗漏最后一条
				time_30 = append(time_30, minbuf_30)
			}
		} else {
			time_30 = append(time_30, minbuf_30)
			minbuf_30 = nil
			minbuf_30 = append(minbuf_30, min)

			st_30 = st_30 + 30 //更新临界时间

			if st_30%100 == 60 {
				st_30 = (st_30/100 + 1) * 100
			}
			if 1130 < st_30 && st_30 < 1300 {
				st_30 = 1300 + 30
			}
		}

		//60
		if min%10000 <= int32(st_60) {
			minbuf_60 = append(minbuf_60, min)
			if i == length-1 { //别遗漏最后一条
				time_60 = append(time_60, minbuf_60)
			}
		} else {
			time_60 = append(time_60, minbuf_60)
			minbuf_60 = nil
			minbuf_60 = append(minbuf_60, min)

			st_60 = st_60 + 100
			if 1130 < st_60 && st_60 < 1300 {
				st_60 = 1300 + 100
			}
		}
	}

	day.Time_5 = &time_5
	day.Time_15 = &time_15
	day.Time_30 = &time_30
	day.Time_60 = &time_60
	return
}
