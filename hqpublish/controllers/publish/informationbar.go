// 证券快照 获取信息栏数据
package publish

import (
	. "haina.com/market/hqpublish/controllers"
	"haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish"

	"haina.com/share/logging"

	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"github.com/gin-gonic/gin"
)

type InfoBar struct{}

func NewInfoBar() *InfoBar {
	return &InfoBar{}
}

func (this *InfoBar) POST(c *gin.Context) {
	replayfmt := c.Query(models.CONTEXT_FORMAT)
	if len(replayfmt) == 0 {
		replayfmt = "pb" // 默认
	}

	switch replayfmt {
	case "json":
		this.PostJson(c)
	case "pb":
		this.PostPB(c)
	default:
		return
	}
}

func (this *InfoBar) PostJson(c *gin.Context) {
	var req protocol.RequestSnapshot

	req.SID = 100000001 // 上证指数
	datash, err := publish.NewIndexSnapshot().GetIndexSnapshotObj(&req)
	if err != nil {
		logging.Error("%v", err)
		WriteJson(c, 40002, nil)
		return
	}
	req.SID = 200399001 // 深圳成指
	datasz, err := publish.NewIndexSnapshot().GetIndexSnapshotObj(&req)
	if err != nil {
		logging.Error("%v", err)
		WriteJson(c, 40002, nil)
		return
	}
	req.SID = 200399006 // 创业板
	datacy, err := publish.NewIndexSnapshot().GetIndexSnapshotObj(&req)

	if err != nil {
		logging.Error("%v", err)
		WriteJson(c, 40002, nil)
		return
	}

	payload := &protocol.PayloadInfoBar{
		List: make([]*protocol.Infobar, 0, 10),
	}
	payload.List = append(payload.List, DataTreating(1, datash))
	payload.List = append(payload.List, DataTreating(2, datasz))
	payload.List = append(payload.List, DataTreating(3, datacy))

	WriteJson(c, 200, payload)
}

func (this *InfoBar) PostPB(c *gin.Context) {
	var req protocol.RequestSnapshot

	req.SID = 100000001 // 上证指数
	datash, err := publish.NewIndexSnapshot().GetIndexSnapshotObj(&req)
	if err != nil {
		logging.Error("%v", err)
		WriteJson(c, 40002, nil)
		return
	}
	req.SID = 200399001 // 深圳成指
	datasz, err := publish.NewIndexSnapshot().GetIndexSnapshotObj(&req)
	if err != nil {
		logging.Error("%v", err)
		WriteJson(c, 40002, nil)
		return
	}
	req.SID = 200399006 // 创业板
	datacy, err := publish.NewIndexSnapshot().GetIndexSnapshotObj(&req)
	if err != nil {
		logging.Error("%v", err)
		WriteJson(c, 40002, nil)
		return
	}
	payload := &protocol.PayloadInfoBar{
		List: make([]*protocol.Infobar, 0, 10),
	}
	payload.List = append(payload.List, DataTreating(1, datash))
	payload.List = append(payload.List, DataTreating(2, datasz))
	payload.List = append(payload.List, DataTreating(3, datacy))

	WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_INFOBAR, payload)
}

// 处理证券快照
func DataTreating(ind int, pst *protocol.IndexSnapshot) *protocol.Infobar {
	var sname = ""
	if ind == 1 {
		sname = "上证指数"
	}
	if ind == 2 {
		sname = "深圳成指"
	}
	if ind == 3 {
		sname = "创业板指"
	}

	return &protocol.Infobar{
		NSID:       pst.NSID,
		SzSName:    sname,
		NLastPx:    pst.NLastPx,
		LlVolume:   pst.LlVolume,
		NPxChg:     pst.NPxChg,
		PxChgRatio: pst.NPxChgRatio,
	}
}
