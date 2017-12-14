// 市场周期分类的资金流向
package szdb

import (
	"haina.com/share/logging"

	"haina.com/share/models"
)

type SZ_HQ_MARKETFUNDFLOW_PERIOD struct {
	SZ_HQ_MARKETFUNDFLOW
	models.Model `db:"-" `
}

func NewSZ_HQ_MARKETFUNDFLOW_PERIOD() *SZ_HQ_MARKETFUNDFLOW_PERIOD {
	return &SZ_HQ_MARKETFUNDFLOW_PERIOD{
		Model: models.Model{
			CacheKey:  "redis_key",
			TableName: TABLE_SZ_HQ_MARKETFUNDFLOW_PERIOD,
			Db:        models.DBSZ,
		},
	}
}

func (s *SZ_HQ_MARKETFUNDFLOW_PERIOD) GetMarketFundFlow(count int32, marketID int32) ([]SZ_HQ_MARKETFUNDFLOW_PERIOD, error) {
	var capflow []SZ_HQ_MARKETFUNDFLOW_PERIOD

	exps := map[string]interface{}{
		"MARKETID=?": marketID,
	}
	builder := s.Db.Select("MARKETID,TRADEDATE,HUGEBUYVALUE,BIGBUYVALUE,HUGESELLVALUE,BIGSELLVALUE").From(s.TableName)
	_, err := s.SelectWhere(builder, exps).OrderBy("TRADEDATE desc").Limit(uint64(count)).LoadStructs(&capflow)
	if err != nil {
		logging.Error("%s", err.Error())
		return capflow, err
	}
	return capflow, err
}
