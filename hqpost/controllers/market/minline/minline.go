package minline

import (
	"ProtocolBuffer/format/kline"
	"fmt"

	"haina.com/market/hqpost/config"

	"github.com/golang/protobuf/proto"
	"haina.com/market/hqpost/models"
	"haina.com/market/hqpost/models/redis_minline"
	redigo "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
)

func NewMinKline(sids *[]int32, cg *config.AppConfig) *MinKline {
	cfg = cg
	return &MinKline{
		sids: sids,
	}
}

//所有股票当天的分钟线数据
func (this *MinKline) HMinLine_1() {
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
				return
			}
		}

		//min 单个股票当天的历史分钟线
		/// do something
		All = append(All, min)

	}
	this.list.All = &All
}

//单个股票当天分钟线基本数据
func (this *MinKline) getBasicMinDataToday(sid int32) (*SingleMin, error) {
	var kinfo kline.HMinLineDay

	day := &SingleMin{
		Time: make([]int32, 0), //单个股票的历史日期
		Min:  make(map[int32]kline.KInfo),
	}
	dmin, err := redis_minline.NewMinKLine(REDISKEY_SECURITY_MIN).GetMinKLineToday(sid)
	if err != nil {
		return nil, err
	}

	for _, v := range *dmin {
		//logging.Debug("sid:%v---v.Time:%v", sid, v.NTime)
		day.Min[v.NTime] = *v
		day.Time = append(day.Time, v.NTime)

		//v 分钟线, 在此处生成今天的1分钟线数据
		kinfo.List = append(kinfo.List, v)
	}
	kinfo.Date = GetDateToday()
	/*********************1分钟线操作**************************************/

	this.mergeMin(sid, REDISKEY_SECURITY_HMIN, kinfo)
	/***********************1分钟历史线操作OVER*******************************/

	day.Sid = sid

	models.GetASCIntArray(day.Time) //升序排序time
	generateMinLineTimes(day)
	return day, nil
}

//当天分钟线并入历史
func (this *MinKline) mergeMin(sid int32, minkey string, kinfo kline.HMinLineDay) {
	var data []byte

	if minkey != REDISKEY_SECURITY_HMIN { //如果不是历史1分钟线的话，进行redis操作（也就是说minline_1不进redis）
		//redis
		hmin, err := getHMinKline(sid, minkey)
		if err != nil {
			return
		}

		key := fmt.Sprintf(minkey, sid)
		if hmin != nil { //已存在历史数据
			hmin.List = append(hmin.List, &kinfo) //并入历史

			data, err = proto.Marshal(hmin)
			if err != nil {
				logging.Error("历史分钟线PB出错...%v", err.Error())
				return
			}
			if err := redis_minline.WriteHMinLine(key, data); err != nil {
				return
			}

		} else { //不存在历史数据
			newmin := kline.HMinTable{}
			newmin.List = append(newmin.List, &kinfo) //并入历史

			data, err = proto.Marshal(&newmin)
			if err != nil {
				logging.Error("历史分钟线PB出错...%v", err.Error())
				return
			}

			if err := redis_minline.WriteHMinLine(key, data); err != nil {
				return
			}
		}
	}

	//文件
	hdata := kline.HMinTable{}
	switch minkey {
	case REDISKEY_SECURITY_HMIN:
		mindata, err := KlineReadFile(sid, cfg.File.Min)
		if err != nil {
			logging.Error("%v", err.Error())
			return
		}
		if mindata != nil { //有历史文件存在
			if err = proto.Unmarshal(mindata, &hdata); err != nil {
				logging.Error("%v", err.Error())
				return
			}
		}

		hdata.List = append(hdata.List, &kinfo)
		bs, err := proto.Marshal(&hdata)
		if err != nil {
			logging.Error("%v", err.Error())
			return
		}
		KlineWriteFile(sid, cfg.File.Min, &bs)
		break
	case REDISKEY_SECURITY_HMIN5:
		mindata, err := KlineReadFile(sid, cfg.File.Min5)
		if err != nil {
			logging.Error("%v", err.Error())
			return
		}
		if mindata != nil { //有历史文件存在
			if err = proto.Unmarshal(mindata, &hdata); err != nil {
				logging.Error("%v", err.Error())
				return
			}
		}

		hdata.List = append(hdata.List, &kinfo)
		bs, err := proto.Marshal(&hdata)
		if err != nil {
			logging.Error("%v", err.Error())
			return
		}
		KlineWriteFile(sid, cfg.File.Min5, &bs)
		break
	case REDISKEY_SECURITY_HMIN15:
		mindata, err := KlineReadFile(sid, cfg.File.Min15)
		if err != nil {
			logging.Error("%v", err.Error())
			return
		}
		if mindata != nil { //有历史文件存在
			if err = proto.Unmarshal(mindata, &hdata); err != nil {
				logging.Error("%v", err.Error())
				return
			}
		}

		hdata.List = append(hdata.List, &kinfo)
		bs, err := proto.Marshal(&hdata)
		if err != nil {
			logging.Error("%v", err.Error())
			return
		}
		KlineWriteFile(sid, cfg.File.Min15, &bs)
		break
	case REDISKEY_SECURITY_HMIN30:
		mindata, err := KlineReadFile(sid, cfg.File.Min30)
		if err != nil {
			logging.Error("%v", err.Error())
			return
		}
		if mindata != nil { //有历史文件存在
			if err = proto.Unmarshal(mindata, &hdata); err != nil {
				logging.Error("%v", err.Error())
				return
			}
		}

		hdata.List = append(hdata.List, &kinfo)
		bs, err := proto.Marshal(&hdata)
		if err != nil {
			logging.Error("%v", err.Error())
			return
		}
		KlineWriteFile(sid, cfg.File.Min30, &bs)
		break
	case REDISKEY_SECURITY_HMIN60:
		mindata, err := KlineReadFile(sid, cfg.File.Min60)
		if err != nil {
			logging.Error("%v", err.Error())
			return
		}
		if mindata != nil { //有历史文件存在
			if err = proto.Unmarshal(mindata, &hdata); err != nil {
				logging.Error("%v", err.Error())
				return
			}
		}

		hdata.List = append(hdata.List, &kinfo)
		bs, err := proto.Marshal(&hdata)
		if err != nil {
			logging.Error("%v", err.Error())
			return
		}
		KlineWriteFile(sid, cfg.File.Min60, &bs)
		break
	default:
		break
	}

}

//获取历史分钟k线
func getHMinKline(sid int32, minkey string) (*kline.HMinTable, error) {
	//获取历史分钟k线
	hmin, err := redis_minline.NewHMinKLine(minkey).GetHMinKLine(sid)
	if err != nil {
		if err != redigo.ErrNil && err != models.ERROR_REDIS_LIST_NULL {
			//logging.Error("%v", err.Error())
			return nil, err
		}

		//可能是第一次获取历史数据为空，也可能是获取失败。
		//logging.Error("%v", err.Error())
		return nil, nil
	}
	return hmin, nil
}

//生成分钟线时间（5、15、30、60）[][]into2
func generateMinLineTimes(day *SingleMin) {
	//logging.Debug("****day_time:%v*******", day.Time)

	var minbuf_5, minbuf_15, minbuf_30, minbuf_60 []int32
	var time_5, time_15, time_30, time_60 [][]int32

	st_05 := MIN_START + 5  // 0935
	st_15 := MIN_START + 15 // 0945
	st_30 := 1000           // 930 + 30
	st_60 := 1030           // 1030
	for _, min := range day.Time {

		//5
		if min <= int32(st_05) {
			minbuf_5 = append(minbuf_5, min)
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
		if min <= int32(st_15) {
			minbuf_15 = append(minbuf_15, min)
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
		if min <= int32(st_30) {
			minbuf_30 = append(minbuf_30, min)
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
		if min <= int32(st_60) {
			minbuf_60 = append(minbuf_60, min)
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
