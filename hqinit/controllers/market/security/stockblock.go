// 整理证券板块信息 生成xml和入redis
package security

import (
	"haina.com/market/hqinit/config"
	"haina.com/market/hqinit/controllers"
	"haina.com/market/hqinit/servers"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

func UpdateStockBlockSet(cfg *config.AppConfig) {

	lib.CheckDir(cfg.File.Path)
	//lib.CheckDir(cfg.File.Path)
	_, err := controllers.OpenFile(cfg.File.Securitiesplate)
	//_, err := controllers.OpenFile(cfg.File.Securitiesplate)

	if err != nil {
		logging.Info("创建文件失败！")
	}
	// 生成板块xml文件
	//servers.NewStockBlockXML().CreateStockblockXML(cfg)
	// 入redis
	servers.NewStockBlockRedis().Block()
}
