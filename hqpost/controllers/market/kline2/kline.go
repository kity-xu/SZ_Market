package kline2

import (
	"ProtocolBuffer/projects/hqpost/go/protocol"
	"fmt"

	"haina.com/market/hqpost/models/kline"
	"haina.com/market/hqpost/models/lib"
	"haina.com/market/hqpost/models/redistore"
	"haina.com/share/logging"

	"haina.com/market/hqpost/config"
)

var cfg *config.AppConfig

func InitPath(c *config.AppConfig) {
	cfg = c
}

//------------------------------------------------------------BeyondDay线的第一次生成-----------------------------------------------------//
// 包括日线、周线、月线、年线
type BedayInterface interface {
	// 获取数据源
	getSrcFromDB()
	getDays()
	getWeeks()
	getMonths()
	getYears()

	// 历史数据 map[int32]protocol.KInfo
	sercurityInfo()

	// 生成历史K线
	makeKlines()
}

//单个股票
type BaseDayKline struct {
	Sid   int32 //股票SID
	Type  uint8 // stock? or index?
	IsNew uint8 // 是否新股

	Date      *[]int32                 //单个股票的历史日期
	WeekDays  *[][]int32               //单个股票的历史周天
	MonthDays *[][]int32               //单个股票的历史月天
	YearDays  *[][]int32               //单个股票的历史年天
	SigStock  map[int32]protocol.KInfo //单个股票的历史数据

	today *protocol.KInfo //单个股票的当天数据
}

func NewBaseDayKline(sid int32) *BaseDayKline {
	sbase, err := redistore.GetSecurityBase(sid)
	if err != nil {
		return &BaseDayKline{
			Sid:   sid,
			Type:  'X',
			IsNew: 'X',
		}
	}
	return &BaseDayKline{
		Sid:   sid,
		Type:  sbase.SzStatus[1],
		IsNew: sbase.SzSCName[0],
	}
}

// 获取证券历史数据（日线）from finchina
func (this *BaseDayKline) getSrcFromDB() {
	exchange, err := kline.GetExchange(this.Sid)
	if err != nil {
		return
	}
	var src_path string
	if this.Type == 'S' {
		src_path = fmt.Sprintf("%s%s%d/%s", cfg.File.Path, exchange, this.Sid, cfg.File.Day)
	} else if this.Type == 'I' {
		src_path = fmt.Sprintf("%s%s%d/%s", cfg.File.Path, exchange, this.Sid, cfg.File.Index)
	} else {
		logging.Debug("non-existent stock status: %s", this.Type)
		return
	}

	if !lib.IsFileExist(src_path) {
		logging.Error("%s", "non-existent data sources, finchina DB")
		return
	}
}

//------------------------------------------------------------beyondDay线的后续追加-----------------------------------------------------//
func NewBedaySubsequent(sid int32) *BedaySubsequent {
	return &BedaySubsequent{
		Sid: sid,
	}
}

type BedaySubsequent struct {
	Sid int32
}

// 获取当天的KInfo
func (this *BedaySubsequent) GetIntradayKInfo() {

}

// 更新hgs_file(海纳行情文件系统)
func (this *BedaySubsequent) UpdateFileStore() {

}

//-------------------------------------------------------------------------------------------------------------------------------------//
