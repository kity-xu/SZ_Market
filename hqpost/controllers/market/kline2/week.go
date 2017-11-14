package kline2

import (
	"ProtocolBuffer/projects/hqpost/go/protocol"
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"

	"io"

	"errors"
	"fmt"

	"haina.com/market/hqpost/config"
	"haina.com/market/hqpost/models/filestore"
	"haina.com/market/hqpost/models/kline"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type BaseLine struct {
	sid      int32                     //股票SID
	date     []int32                   //单个股票的历史日期
	sigStock map[int32]*protocol.KInfo //单个股票的历史数据
	kindDays *[][]int32
}

// 生成周线
//this.WeekDay
func HisWeekKline(sids *[]int32) {
	base := new(BaseLine)
	for _, sid := range *sids {
		base = nil
		weekName := filePath(cfg, cfg.File.Week, sid)
		if !lib.IsFileExist(weekName) { // 不存在或其他做第一从生成 TODO
			if err := base.CreateWeekLine(sid, weekName, today); err != nil {
				continue
			}
		} else { // 追加周线
			filestore.UpdateWeekLineToFile(weekName, today)
		}
	}
}

// 读day文件生成week
func (this *BaseLine) CreateWeekLine(sid int32, weekFile string, today *protocol.KInfo) error {
	this.sid = sid
	sigle := make(map[int32]*protocol.KInfo)

	spath, is := kline.IsExistdirInHGSFileStore(cfg, cfg.File.Day, sid) // 目录是否存在
	if !is {
		logging.Error("Create WeekLine: DayLine no Exist in hgs_file")
		return errors.New("error")
	}
	dayPath := fmt.Sprintf("%s/%d.dat", spath, sid)
	bs, err := ioutil.ReadFile(dayPath)
	if err != nil || len(bs) == 0 {
		logging.Error("Create WeekLine: Read dayLine error|%v", err)
		return err
	}

	buff := &protocol.KInfo{}
	size := binary.Size(buff)
	for i := 0; i < len(bs); i += size {
		buffer := &protocol.KInfo{}
		if err = binary.Read(bytes.NewBuffer(bs[i:size+i]), binary.LittleEndian, buffer); err != nil && err != io.EOF {
			logging.Error("Create WeekLine: binary read dayline error|%v", err)
			return err
		}
		this.date = append(this.date, buffer.NTime)
		sigle[buffer.NTime] = buffer
	}
	this.sigStock = sigle
	this.GetSecurityWeekDay()
	wTable := this.produceWeeprotocol()
	filestore.MaybeBelongAWeek(wTable, today) //第一次生成的时候 如果同属一周加入当天数据
	err = filestore.WiteHainaFileStore(weekFile, wTable)
	if err != nil {
		logging.Error("WiteHainaFileStore error | %v", err)
	}
	return nil
}

// 生成周线日期
func (this *BaseLine) GetSecurityWeekDay() {
	if len(this.date) < 1 {
		logging.Error("Create WeekLine:%v---No historical data...", v.Sid)
		return
	}
	var wday [][]int32
	sat, _ := filestore.DateAdd(this.date[0]) //该股票第一个交易日所在周的周日（周六可能会有交易）

	var dates []int32
	for j, date := range this.date {
		if filestore.IntToTime(int(date)).Before(sat) {
			dates = append(dates, date)
			if j == int(len(this.date)-1) { //执行到最后一个
				wday = append(wday, dates)
			}
		} else {
			wday = append(wday, dates)

			sat, _ = filestore.DateAdd(date)
			dates = nil
			dates = append(dates, date)
		}
	}
	this.kindDays = &wday
}

// 生成周K线
func (this *BaseLine) produceWeeprotocol() *protocol.KInfoTable {
	var tmps []protocol.KInfo
	var klist protocol.KInfoTable

	for _, week := range this.kindDays { //每一周
		tmp := protocol.KInfo{}

		var (
			i          int
			day        int32
			AvgPxTotal uint32
		)
		for i, day = range week { //每一天
			stockday := this.sigStock[day]
			if tmp.NHighPx < stockday.NHighPx || tmp.NHighPx == 0 { //最高价
				tmp.NHighPx = stockday.NHighPx
			}
			if tmp.NLowPx > stockday.NLowPx || tmp.NLowPx == 0 { //最低价
				tmp.NLowPx = stockday.NLowPx
			}
			tmp.LlVolume += stockday.LlVolume //成交量
			tmp.LlValue += stockday.LlValue   //成交额
			AvgPxTotal += stockday.NAvgPx
		}

		tmp.NSID = this.sid
		tmp.NTime = this.sigStock[week[0]].NTime     //时间（取每周第一天）
		tmp.NOpenPx = this.sigStock[week[0]].NOpenPx //开盘价（每周第一天的开盘价）
		tmp.NPreCPx = this.sigStock[week[0]].NPreCPx //本周一的昨收
		tmp.NLastPx = this.sigStock[week[i]].NLastPx //最新价
		tmp.NAvgPx = AvgPxTotal / uint32(i+1)        //平均价
		tmps = append(tmps, tmp)
		logging.Debug("周线是:%v", tmps)
		klist.List = append(klist.List, &tmp)
	}
	return &klist
}

func filePath(cfg *config.AppConfig, kind string, sid int32) string {
	dir, _ := kline.IsExistdirInHGSFileStore(cfg, kind, sid)
	lib.CheckDir(dir)
	return fmt.Sprintf("%s/%d.dat", dir, sid)
}
