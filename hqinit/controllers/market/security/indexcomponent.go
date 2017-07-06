// 整理指数成分股信息 生成xml
package security

import (
	"haina.com/market/hqinit/config"
	"haina.com/market/hqinit/controllers"
	"haina.com/market/hqinit/servers"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

func UpdateIndexComponent(cfg *config.AppConfig) {

	lib.CheckDir(cfg.File.Path)
	_, err := controllers.OpenFile(cfg.File.IndexComponentPath)

	if err != nil {
		logging.Info("create indexcomponent.xml error ！")
	}

	// 生成指数xml文件
	servers.NewIndexComponentXML().CreateIndexComponentXML(cfg)
}
