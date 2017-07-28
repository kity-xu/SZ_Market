package main

import (
	"time"

	"haina.com/market/hqtools/stockindex/control"
	"haina.com/share/logging"
)

func main() {
	start := time.Now()

	control.Operation()

	end := time.Now()
	logging.Info("Update Kline historical data successed, and running time:%v", end.Sub(start))
}
