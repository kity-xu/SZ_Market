package market

import (
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"haina.com/market/hqpost/config"
	"haina.com/market/hqpost/controllers/market/kline"
	"haina.com/market/hqpost/controllers/market/minline"
	"haina.com/market/hqpost/controllers/market/sidcode"
	"haina.com/share/logging"
)

func Update(cfg *config.AppConfig) {
	/*********************开始时间************************/
	start := time.Now()
	var sids *[]int32
	var err error

	if os.Args[2] == "-of" {
		os.Args[2] = "-o"
		sids = GetAllSids(cfg.File.Finpath)
		logging.Info("start sids all by walk dir... len:%v", len(*sids))
	} else {
		sids, err = sidcode.GetSecurityTable()
		if err != nil {
			return
		}
		logging.Info("start sids by walk stock table... len:%v", len(*sids))
	}
	//---------------------日周月年线---------------------/
	//K线
	if len(os.Args) == 3 {
		min := minline.NewMinKline(sids, cfg)
		kline.InitPath(cfg)

		switch os.Args[2] {
		case "-m":
			logging.Info("--------min----------")
			min.InitMinLine()
			// 分钟线
			min.HMinLine_1()
			min.HMinLine_5()
			min.HMinLine_15()
			min.HMinLine_30()
			min.HMinLine_60()
		case "-o":
			logging.Info("--------oDay----------")
			kline.HisDayKline(sids)
			kline.HisWeekKline(sids)
			kline.HisMonthKline(sids)
			kline.HisYearKline(sids)
		case "-day":
			logging.Info("--------day----------")
			kline.HisDayKline(sids)
		case "-week":
			logging.Info("--------week----------")
			kline.HisWeekKline(sids)
		case "-month":
			logging.Info("--------month----------")
			kline.HisMonthKline(sids)
		case "-year":
			logging.Info("--------year----------")
			kline.HisYearKline(sids)
		case "-m01":
			logging.Info("--------min01----------")
			min.InitMinLine()
			min.HMinLine_1()
		case "-m05":
			logging.Info("--------min05----------")
			min.InitMinLine()
			min.HMinLine_5()
		case "-m15":
			logging.Info("--------min15----------")
			min.InitMinLine()
			min.HMinLine_15()
		case "-m30":
			logging.Info("--------min30----------")
			min.InitMinLine()
			min.HMinLine_30()
		case "-m60":
			logging.Info("--------min60----------")
			min.InitMinLine()
			min.HMinLine_60()
		default:
			return
		}

	} else {
		logging.Info("--------all----------")
		kline.InitPath(cfg)
		kline.HisDayKline(sids)
		kline.HisWeekKline(sids)
		kline.HisMonthKline(sids)
		kline.HisYearKline(sids)
		//--------------------------------------------------/

		//----------------------分钟线-----------------------/
		min := minline.NewMinKline(sids, cfg)
		min.InitMinLine()
		// 分钟线
		min.HMinLine_1()
		min.HMinLine_5()
		min.HMinLine_15()
		min.HMinLine_30()
		min.HMinLine_60()
		//--------------------------------------------------/
	}
	//----------------------板块指数线----------------------/
	//bline.NewBlockIndex().UpdateBlockIndexDayLine(cfg)
	//--------------------------------------------------/

	end := time.Now()
	logging.Info("Update Kline historical data successed, and running time:%v", end.Sub(start))

	//redistore.NewHKLine(REDISKEY_HQPOST_EXECUTED_TIME).HQpostExecutedTime()
	/*********************结束时间***********************/
}

// 第一次执行时遍历所有曾上市的股票id，
// 因为当前快照 sid < 所有曾上市的sid
//获取指定目录下的所有目录，不进入下一级目录搜索。
func ListDir(dirPth string) (dirs []int32, err error) {
	dirs = make([]int32, 0, 1000)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	for _, fi := range dir {
		if fi.IsDir() {
			sid, e := strconv.Atoi(fi.Name())
			if e != nil {
				logging.Error("Invalid sid by Walk... %v", fi.Name())
				continue
			}
			dirs = append(dirs, int32(sid))
		}
	}
	return dirs, nil
}

func GetAllSids(path string) *[]int32 {
	sh, err := ListDir(path + "/sh")
	if err != nil {
		logging.Error("Fist running... GetAllSids sh err | %v", err.Error())
		return nil
	}

	sz, err := ListDir(path + "/sz")
	if err != nil {
		logging.Error("Fist running... GetAllSids sz err | %v", err.Error())
		return nil
	}
	sh = append(sh, sz...)
	return &sh
}
