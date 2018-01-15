// 历史周、月、年资金流向
package szdb

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"

	"haina.com/share/models"
)

type SZ_HQ_SECURITYFUNDFLOW_PERIOD struct {
	models.Model `db:"-" `
	SID          int32 //证券ID
	PERIODID     int32 //周期ID
	//TRADEMONTH     int32           //月份
	//TRADEWEEK      int32           //该月第几周
	LASTDATE     int32           //日期
	ENTRYTIME    string          //库更新时间点
	HUGEBUYVALUE dbr.NullFloat64 //#特大买单成交额(元)
	BIGBUYVALUE  dbr.NullFloat64 //#大买单成交额(元)
	//MIDDLEBUYVALUE dbr.NullFloat64 //#中买单成交额(元)
	//SMALLBUYVALUE  dbr.NullFloat64 //#小买单成交额(元)

	HUGESELLVALUE dbr.NullFloat64 //#特大卖单成交额(元)
	BIGSELLVALUE  dbr.NullFloat64 //#大卖单成交额(元)
	//MIDDLESELLVALUE dbr.NullFloat64 //#中卖单成交额(元)
	//SMALLSELLVALUE  dbr.NullFloat64 //#小卖单成交额(元)
}

func NewSZ_HQ_SECURITYFUNDFLOW_PERIOD() *SZ_HQ_SECURITYFUNDFLOW_PERIOD {
	return &SZ_HQ_SECURITYFUNDFLOW_PERIOD{
		Model: models.Model{
			CacheKey:  "redis_key",
			TableName: TABLE_SZ_HQ_SECURITYFUNDFLOW_PERIOD,
			Db:        models.DBSZ,
		},
	}
}

// 股票周、月、年的资金流向
func (s *SZ_HQ_SECURITYFUNDFLOW_PERIOD) GetSecurityFundFlowPeriod(sid int32, periodID int32) ([]SZ_HQ_SECURITYFUNDFLOW_PERIOD, error) {
	var capflow []SZ_HQ_SECURITYFUNDFLOW_PERIOD

	exps := map[string]interface{}{
		"SID=?":      sid,
		"PERIODID=?": periodID,
	}
	builder := s.Db.Select("*").From(s.TableName)
	_, err := s.SelectWhere(builder, exps).OrderBy("ENTRYDATE desc").LoadStructs(&capflow)
	if err != nil {
		logging.Error("%s", err.Error())
		return capflow, err
	}
	return capflow, err
}
