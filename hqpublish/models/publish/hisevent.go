// 公告信息集合
package publish

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"
	"strconv"
	"strings"

	. "haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/fcmysql"
	"haina.com/share/logging"
)

type HisEventinfo struct {
}

func NewHisEventinfo() *HisEventinfo {
	return &HisEventinfo{}
}

// 公告信息
func (this *HisEventinfo) GetHisevent(req *protocol.RequestHisevent) (*protocol.PayloadHisevent, error) {

	// 根据证券id获取公告信息
	hid := strconv.Itoa(int(req.HiseventID))
	hsi, err := fcmysql.NewTQ_SK_ANNOUNCEMT().GetHisEvent(hid)
	if err != nil {
		logging.Info("mysql select error %v", err)
		return nil, err
	}

	var noti protocol.PayloadHisevent

	noti.NHiseventID = req.HiseventID
	noti.NDeclardate = hsi.DECLAREDATE
	noti.SzHeadline = hsi.ANNTITLE.String
	// 如果公告内容等于 “公告内容详见附件” 需查询公告目录表
	anncmt, err := fcmysql.NewTQ_OA_ANNTFILE().GetAnntfile(hsi.ANNOUNCEMTID.String)
	if err != nil {
		logging.Info("查询公告目录 error %v", err)
	}
	//if hsi.ANNTEXT.String == "公告内容详见附件" {
	str := strings.Replace(anncmt.FILELINK.String, `\`, "/", -1)
	urlstr := FCat.Url + str
	noti.SzWebtake = urlstr
	// 公告内容类型
	noti.SzAcsyType = "url"
	//	} else {
	//		noti.SzWebtake = hsi.ANNTEXT.String
	//		// 公告内容类型
	//		noti.SzAcsyType = "txt"
	//	}

	noti.SzNoticeType = hsi.ANNTYPE.String
	isif := hsi.LEVEL1.String
	if len(isif) < 1 {
		isif = "no"
	}
	switch isif {
	case "no":
		noti.NWhether = 12
	case "回报规划":
		noti.NWhether = 12
	case "回购":
		noti.NWhether = 11
	case "获取认证":
		noti.NWhether = 12
	case "基本资料变动":
		noti.NWhether = 11
	case "减持":
		noti.NWhether = 11
	case "借款":
		noti.NWhether = 11
	case "其它":
		noti.NWhether = 12
	case "权证发行(上市)":
		noti.NWhether = 11
	case "权证开盘参考价":
		noti.NWhether = 12
	case "权证行权":
		noti.NWhether = 12
	case "权证终止上市":
		noti.NWhether = 11
	case "日期变动":
		noti.NWhether = 11
	case "实际控制人变更":
		noti.NWhether = 11
	case "收购/出售股权（资产）":
		noti.NWhether = 11
	case "税率变动":
		noti.NWhether = 11
	case "诉讼仲裁":
		noti.NWhether = 11
	case "停牌":
		noti.NWhether = 11
	case "投资设立(参股)公司":
		noti.NWhether = 11
	case "投资项目":
		noti.NWhether = 11
	case "违规":
		noti.NWhether = 11
	case "委托理财":
		noti.NWhether = 12
	case "委员会成员变动":
		noti.NWhether = 11
	case "新创设权证":
		noti.NWhether = 11
	case "信托":
		noti.NWhether = 11
	case "行权价格(比例)调整":
		noti.NWhether = 11
	case "要约收购":
		noti.NWhether = 11
	case "业绩预测":
		noti.NWhether = 12
	case "再融资预案":
		noti.NWhether = 11
	case "暂停上市风险":
		noti.NWhether = 11
	case "增持":
		noti.NWhether = 11
	case "增持解锁":
		noti.NWhether = 11
	case "质押":
		noti.NWhether = 11
	case "中介机构变动":
		noti.NWhether = 11
	case "终止上市风险":
		noti.NWhether = 11
	case "重大合同":
		noti.NWhether = 11
	case "重大事故":
		noti.NWhether = 11
	case "注销权证":
		noti.NWhether = 11
	case "追加限售":
		noti.NWhether = 11
	case "资产(债务)重组":
		noti.NWhether = 11
	case "资金占用":
		noti.NWhether = 11
	case "短期融资券":
		noti.NWhether = 11

	}

	return &noti, nil
}
