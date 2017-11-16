package kline

import (
	"ProtocolBuffer/projects/hqpost/go/protocol"
	"fmt"

	"haina.com/share/lib"

	cl "haina.com/market/hqpost/controllers"
	"haina.com/market/hqpost/models/filestore"
	"haina.com/market/hqpost/models/kline"
	"haina.com/market/hqpost/models/redistore"
	"haina.com/share/logging"
)

// HisDayKline ... 日K
func HisDayKline(sids *[]int32) {

	for _, sid := range *sids {
		// 获取当天快照
		today, err := GetIntradayKInfo(sid)
		if err != nil {
			logging.Error("严重错误！！！，程序被迫停止执行 |%v ", err)
			return
		}

		this := NewBaseDayKline(sid)

		spath, srcExist := kline.IsExistFileInHGSFileStore(cfg, cfg.File.Day, sid)
		if !srcExist { // hgs_file 文件系统不存在该目录 --> hgs_file/sh/day
			dpath, _ := kline.IsExistdirInThirdDB(cfg, sid)
			thirdFile := ThirdFilepath(dpath)
			if !lib.IsFileExist(thirdFile) {
				// 此时hgs_file和thirdDB中都不存在该sid的数据(可能的情况是：sid无效；股票新上市)
				if this.IsNew != 'N' {
					logging.Error("Error: Invalid sid | %d", sid)
					continue
				}
				logging.Info("这是一支新股？")
				// 新股 AddHgsFile()
				if e := filestore.AppendFile(kline.HGSFilepath(spath), today); e != nil {
					logging.Error("Error: Append File Err | %v", e)
					continue
				}
			} else {
				// 此时数据不在hgs_file中，而是存于ThirdDB (可能的情况是：没做数据的第一次生成; 数据的第一次生成遗漏了该数据; 前一天上市的股票没有更新到hgs_file,此种情况的前提是hqdata每天更新)
				// 以上无论那种情况都需要做该数据的第一次生成
				// 数据做第一次生成 InitHgsFile()
				if err = CreateSingleHgsFile(spath, thirdFile, today); err != nil {
					continue
				}
			}
		} else {
			// hgs_file中有该sid数据，说明有历史数据存在, 则只需后续追加今天Kinfo即可
			if e := filestore.AppendFile(kline.HGSFilepath(spath), today); e != nil {
				logging.Error("Error: Append File Err | %v", e)
				continue
			}
		}
	}
}

// --------------------------------------------------------------------functhions----------------------------------------------------------------//
// GetIntradayKInfo ... 获取当天的KInfo
func GetIntradayKInfo(sid int32) (*protocol.KInfo, error) {
	key := fmt.Sprintf(cl.REDISKEY_SECURITY_SNAP, sid)
	return redistore.GetStockSnapshotObj(key)
}

// CreateSingleHgsFile ... 生成该sid的hgs_file
func CreateSingleHgsFile(srcpath, despath string, today *protocol.KInfo) error {
	klist, err := filestore.ReadSrcFileStore(despath)
	if err != nil {
		logging.Error("%s", err)
		return err
	}
	if today != nil {
		klist.List = append(klist.List, today)
	}
	hgsfile := kline.HGSFilepath(srcpath)

	if err = filestore.WiteHainaFileStore(hgsfile, klist); err != nil {
		logging.Error("%s", err)
		return err
	}
	return nil
}

// ThirdFilepatth ... 得到数据源文件路劲
func ThirdFilepath(despath string) string {
	var desfile string
	desfile = fmt.Sprintf("%s/%s", despath, cfg.File.Finday)
	if !lib.IsFileExist(desfile) {
		desfile = fmt.Sprintf("%s/%s", despath, cfg.File.Index)
	}
	return desfile
}
