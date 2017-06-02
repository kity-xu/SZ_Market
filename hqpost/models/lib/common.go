package lib

import (
	"ProtocolBuffer/projects/hqpost/go/protocol"
	"sort"
	"strconv"
	"time"
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

type IntSlice []int32

func (p IntSlice) Len() int           { return len(p) }
func (p IntSlice) Less(i, j int) bool { return p[j] < p[i] }
func (p IntSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

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

func GetASCIntArray(a []int32) {
	sort.Sort(sort.Reverse(IntSlice(a)))
}

func GetASCStruct(a *[]*protocol.KInfo) {
	sort.Sort(sort.Reverse(KTable(*a)))
}

//降序
func GetSECIntArray(a []int32) {
	sort.Sort(IntSlice(a))
}

func GetSECStruct(a *[]*protocol.KInfo) {
	sort.Sort(KTable(*a))
}

//获取当前时间20170101
func GetCurrentTime() int {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	date, _ := strconv.Atoi(tm.Format("20060102"))
	return date
}
