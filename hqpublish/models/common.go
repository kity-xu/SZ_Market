package models

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"haina.com/share/logging"
)

// 自定义的 Reverse 类型
type Reverse struct {
	sort.Interface // 这样，Reverse可以接纳任何实现了sort.Interface的对象
}

// Reverse 只是将其中的 Inferface.Less 的顺序对调了一下
func (r Reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}

type KTable []*protocol.KInfo

func (this KTable) Len() int {
	return len(this)
}

func (this KTable) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this KTable) Less(i, j int) bool {
	return this[j].NTime < this[i].NTime
}

//升序
func GetASCStruct(a *[]*protocol.KInfo) {
	sort.Sort(sort.Reverse(KTable(*a)))
}

//降序
func GetSECStruct(a *[]*protocol.KInfo) {
	sort.Sort(KTable(*a))
}

//获取当前时间20170101
func GetCurrentTime() int32 {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	date, _ := strconv.Atoi(tm.Format("20060102"))
	return int32(date)
}

//获取当前时间201701010100
func GetCurrentTimeHM() int {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	ss := tm.Format("200601021504")

	dd, err := strconv.Atoi(ss)
	if err != nil {
		logging.Error("%v", err.Error())
		return 0
	}
	return dd
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
