package publish

import (
	"ProtocolBuffer/format/kline"
	"sort"
)

// 自定义的 Reverse 类型
type Reverse struct {
	sort.Interface // 这样，Reverse可以接纳任何实现了sort.Interface的对象
}

// Reverse 只是将其中的 Inferface.Less 的顺序对调了一下
func (r Reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}

type KTable []*kline.KInfo

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
func GetASCStruct(a *[]*kline.KInfo) {
	sort.Sort(sort.Reverse(KTable(*a)))
}

//降序
func GetSECStruct(a *[]*kline.KInfo) {
	sort.Sort(KTable(*a))
}
