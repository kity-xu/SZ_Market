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

	lib.CheckDir("/opt/develop/hqs/filestore/security/")
	//lib.CheckDir("E:/hqfile/")
	_, err := controllers.OpenFile("/opt/develop/hqs/filestore/security/securitiesplate.xml")
	//_, err := controllers.OpenFile("E:/hqfile/securitiesplate.xml")

	if err != nil {
		logging.Info("创建文件失败！")
	}
	// 生成板块xml文件
	servers.NewStockBlockXML().CreateStockblockXML()
	// 入redis
	servers.NewStockBlockRedis().Block()
}
