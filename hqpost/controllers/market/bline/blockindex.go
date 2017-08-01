package bline

import (
	"haina.com/market/hqpost/models/redistore"

	"haina.com/market/hqpost/config"
	"haina.com/share/logging"
)

type BlockIndex struct {
}

func NewBlockIndex() *BlockIndex {
	return &BlockIndex{}
}

func (this *BlockIndex) UpdateBlockIndexDayLine(cfg *config.AppConfig) {
	err := redistore.NewMBlockIndex("hq:bk:snap:*").UpdateBlockIndexLine(cfg.File.BlockPath)
	if err != nil {
		logging.Error("Update block index Dayline failed...%v", err.Error())
		return
	}
	logging.Info("Update block index Dayline successed...")
}
