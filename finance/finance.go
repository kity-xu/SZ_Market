package main

import (
	"haina.com/share/app"
	"haina.com/share/logging"

	"haina.com/market/finance/config"
	. "haina.com/market/finance/models"
	"haina.com/market/finance/routes"

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
	a.SessionOn = cfg.Session.On
	a.SessionProviderName = cfg.Session.ProviderName
	a.SessionConfig = cfg.Session.Config
	a.DisableGzip = true
	a.Cors = cfg.Cors.AllowOrigin

	r := a.Init()

	// 路由注册
	routes.Register(r)

	// 监控性能
	ginpprof.Wrapper(r)

	logging.Error("%s", r.Run(cfg.Serve.Port))
}
