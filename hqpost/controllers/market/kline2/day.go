package kline2

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

// HisDayKline ... 日K线
func HisDayKline(sids *[]int32) {

	for _, sid := range *sids {
		// 获取当天快照
		today, err := GetIntradayKInfo(sid) // TODO 成交量和成交额
		if err != nil || today == nil {
			continue
		}

		this := NewBaseDayKline(sid)

		spath, srcExist := kline.IsExistdirInHGSFileStore(cfg, cfg.File.Day, sid)
		if !srcExist { // hgs_file 文件系统不存在该目录 --> hgs_file/sh/day

			dpath, desExist := kline.IsExistdirInThirdDB(cfg, sid)
			if !desExist {
				// 此时hgs_file和thirdDB中都不存在该sid的数据(可能的情况是：sid无效；股票新上市)
				if this.IsNew != 'N' {
					logging.Error("Error: Invalid sid | %s", sid)
					continue
				}
				// 新股 AddHgsFile()
				if e := filestore.AppendFile(kline.HGSFilepath(spath, sid), today); e != nil {
					logging.Error("Error: Append File Err | %v", e)
					continue
				}
			}
			// 此时数据不在hgs_file中，而是存于ThirdDB (可能的情况是：没做数据的第一次生成; 数据的第一次生成遗漏了该数据; 前一天上市的股票没有更新到hgs_file,此种情况的前提是hqdata每天更新)
			// 以上无论那种情况都需要做该数据的第一次生成
			// 数据做第一次生成 InitHgsFile()
			if err = CreateSingleHgsFile(spath, dpath, sid); err != nil {
				continue
			}
		}
		// hgs_file中有该sid数据，说明有历史数据存在, 则只需后续追加今天Kinfo即可
		if e := filestore.AppendFile(kline.HGSFilepath(spath, sid), today); e != nil {
			logging.Error("Error: Append File Err | %v", e)
			continue
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
func CreateSingleHgsFile(srcpath, despath string, sid int32) error {
	desfile := ThirdFilepath(despath)
	klist, err := filestore.ReadSrcFileStore(desfile)
	if err != nil {
		logging.Error("%s", err)
		return err
	}

	hgsfile := kline.HGSFilepath(srcpath, sid)

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
	if !lib.IsFileExist(desfile5) {
		desfile = fmt.Sprintf("%s/%s", despath, cfg.File.Findex)
	}
	return desfile
}
