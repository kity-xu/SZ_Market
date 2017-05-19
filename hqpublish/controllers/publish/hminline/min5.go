package hminline

import (
	"haina.com/market/hqpublish/models/publish"
	"haina.com/share/logging"
)

func TestMinLine() {
	data, err := publish.NewHMinKLine(REDISKEY_SECURITY_HMIN5).GetHMinKLine(int32(200300026))
	if err != nil {
		logging.Error("%v", err)
		return
	}

	for _, v := range data.List {
		logging.Info("Date----%v", v.Date)
		for _, kinfo := range v.List {
			logging.Info("kinfo----%v", kinfo)
		}
	}
}
