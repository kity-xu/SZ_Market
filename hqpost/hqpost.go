package main

import (
	"haina.com/market/hqpost/controllers/market"
	"haina.com/share/app"
	"haina.com/share/logging"

	"haina.com/market/hqpost/config"
	. "haina.com/market/hqpost/models"

	"github.com/DeanThompson/ginpprof"
)

func main() {
	cfg := config.Default(APP_PID)

	// 项目初始化
	a := app.NewApp(APP_NAME, APP_VERSION)
	a.PidName = APP_PID
	a.WSPort = cfg.Serve.Port
	a.LogPort = cfg.Log.Port
	a.LogAddr = cfg.Log.Addr
	a.LogOn = cfg.Log.On
	a.DisableGzip = true
	a.Cors = cfg.Cors.AllowOrigin

	r := a.Init()

	// 监控性能
	ginpprof.Wrapper(r)

	//行情数据更新
	market.Update(cfg)

	logging.Error("%s", r.Run(cfg.Serve.Port))
}
