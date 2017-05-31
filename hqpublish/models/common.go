package models

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"
	"io"
	"sort"
	"strconv"
	"time"

	"haina.com/share/logging"

	"github.com/gin-gonic/gin"
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

func GetRequestData(c *gin.Context) ([]byte, error) {
	temp := make([]byte, 1024)
	n, err := c.Request.Body.Read(temp)
	if err != nil && err != io.EOF {
		logging.Error("Body Read: %v", err)
		return nil, err
	}
	//logging.Info("\nBody len %d\n%s", n, temp[:n])
	logging.Info("Body len %d", n)
	return temp[:n], nil

}
