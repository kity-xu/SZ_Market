package filestore

import (
	"fmt"

	"ProtocolBuffer/projects/hqpost/go/protocol"
	"errors"
	"strconv"
	"strings"
	"time"
)

//int型时间转Time类型（最小单位 天）
func IntToTime(date int) time.Time {
	swap := date % 10000
	year := date / 10000
	month := swap / 100
	day := swap % 100
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

//查询某一日期所在周的周日
func DateAdd(date int32) (time.Time, error) {
	var sat time.Time
	swap := date % 10000
	year := int(date / 10000)
	month := swap / 100
	day := int(swap % 100)

	//logging.Debug("%d-%d-%d", year, month, day)

	baseTime := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	weekday := baseTime.Weekday().String()

	var basedate string
	if strings.EqualFold(weekday, "Monday") {
		basedate = fmt.Sprintf("%d%s", 24*6, "h")

	} else if strings.EqualFold(weekday, "Tuesday") {
		basedate = fmt.Sprintf("%d%s", 24*5, "h")

	} else if strings.EqualFold(weekday, "Wednesday") {
		basedate = fmt.Sprintf("%d%s", 24*4, "h")

	} else if strings.EqualFold(weekday, "Thursday") {
		basedate = fmt.Sprintf("%d%s", 24*3, "h")

	} else if strings.EqualFold(weekday, "Friday") {
		basedate = fmt.Sprintf("%d%s", 24*2, "h")

	} else if strings.EqualFold(weekday, "Saturday") {
		basedate = fmt.Sprintf("%d%s", 24*1, "h")

	} else {
		//logging.Error("SID:%v------Invalid trade date...%v", sid, date)
		return sat, errors.New("周日有交易..")
	}

	dd, _ := time.ParseDuration(basedate)
	sat = baseTime.Add(dd) //Saturday（星期日）
	return sat, nil
}

//int型时间转Time类型（最小单位 月份）
func IntToMonth(date int) time.Time {
	swap := date % 10000
	year := date / 10000
	month := swap / 100

	return time.Date(year, time.Month(month), 0, 0, 0, 0, 0, time.UTC)
}

//比较历史最新和当天
func compareKInfo(tmp *protocol.KInfo, today *protocol.KInfo) protocol.KInfo {
	var swap protocol.KInfo

	swap.NSID = tmp.NSID
	swap.NTime = tmp.NTime
	swap.NPreCPx = today.NPreCPx
	swap.NOpenPx = tmp.NOpenPx
	if tmp.NHighPx > today.NHighPx {
		swap.NHighPx = tmp.NHighPx
	} else {
		swap.NHighPx = today.NHighPx
	}
	if tmp.NLowPx > today.NLowPx {
		swap.NLowPx = today.NLowPx
	} else {
		swap.NLowPx = tmp.NLowPx
	}
	swap.NLastPx = today.NLastPx
	swap.LlVolume = today.LlVolume + tmp.LlVolume
	swap.LlValue = today.LlValue + tmp.LlValue
	swap.NAvgPx = (today.NAvgPx + tmp.NAvgPx) / 2
	return swap
}

func GetDateToday() int32 {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	date, _ := strconv.Atoi(tm.Format("20060102"))
	return int32(date)
}
