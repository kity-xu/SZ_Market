// 公告信息
package publish

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"
	"strconv"

	"haina.com/market/hqpublish/models/fcmysql"
	"haina.com/share/logging"
)

type NoticeinfoL struct {
}

func NewNoticeinfoL() *NoticeinfoL {
	return &NoticeinfoL{}
}

// 公告信息
func (this *NoticeinfoL) GetNoticeInfoL(req *protocol.RequestNoticeInfo) (*protocol.PayloadNoticeInfo, error) {

	var psb protocol.PayloadNoticeInfo

	sid := strconv.Itoa(int(req.SID))
	// 根据sid查询公司内码
	sid = sid[3:]
	stc, err := fcmysql.NewTQ_OA_STCODE().GetStcodeInfo(sid)
	if err != nil {
		logging.Info("select　TQ_OA_STCODE　error%v", err)
	}
	// 根据证券id获取公告信息
	noif, err := fcmysql.NewTQ_SK_ANNOUNCEMT().GetNoticeInfo(stc.COMPCODE.String)
	if err != nil {
		logging.Info("mysql select error %v", err)
	}

	psb.SID = req.SID
	psb.Num = int32(len(noif))
	for _, ite := range noif {
		var noti protocol.NoticeInfoB
		noti.NDeclardate = ite.DECLAREDATE
		noti.SzHeadline = ite.ANNTITLE.String
		noti.SzWebtake = ite.ANNTEXT.String
		noti.SzNoticeType = ite.ANNTYPE.String
		isif := ite.LEVEL1.String
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
		psb.List = append(psb.List, &noti)

	}
	logging.Info("----------")

	return &psb, nil
}
