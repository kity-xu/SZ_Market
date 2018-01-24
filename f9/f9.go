package main

import (
	"github.com/DeanThompson/ginpprof"
	"haina.com/share/app"
	"haina.com/share/logging"

	"haina.com/market/f9/config"
	"haina.com/market/f9/models"
	"haina.com/market/f9/routes"
)

func main() {
	cfg := config.Default(models.APP_PID)

	a := app.NewApp(models.APP_NAME, models.APP_VERSION)
	a.PidName = models.APP_PID
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
