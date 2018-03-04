package finchina

import (
	"encoding/json"
	"fmt"

	. "haina.com/market/f9/models"
	redigo "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	. "haina.com/share/models"
)

// TQ_OA_STCODE    证券内码表
// ---------------------------------------------------------------------
type TQ_OA_STCODE struct {
	Model    `db:"-"`
	COMPCODE dbr.NullString //公司代码(公司内码) 通过 SYMBOL 得到
	SECODE   dbr.NullString //分市场内码
	SYMBOL   dbr.NullString //股票代码
}

func NewTQ_OA_STCODE() *TQ_OA_STCODE {
	return &TQ_OA_STCODE{
		Model: Model{
			TableName: TABLE_TQ_OA_STCODE,
			Db:        MyCat,
		},
	}
}

func (this *TQ_OA_STCODE) getCompcode(symbol interface{}) error {
	var seg string
	switch symbol.(type) {
	case int, int32, int64:
		seg = fmt.Sprintf("%d", symbol)
	case string:
		seg = fmt.Sprintf("%s", symbol)
	default:
		return fmt.Errorf("Invalid symbol type")
	}

	key := fmt.Sprintf(REDISKEY_SYSMBOL_BASIC, seg)
	v, err := RedisCache.GetBytes(key)
	if err != nil {
		if err != redigo.ErrNil {
			logging.Error("Redis get %s: %s", key, err)
		}
		if len(seg) != 9 {
			return fmt.Errorf("Invalid symbol...")
		}
		var exchange string
		if seg[:3] == "100" {
			exchange = "001002"
		} else if seg[:3] == "200" {
			exchange = "001003"
		} else {
			exchange = ""
		}
		cond := fmt.Sprintf("EXCHANGE=%s and SETYPE in ('101','701') and SYMBOL=%s", exchange, seg[3:])
		err = this.Db.Select("*").From(this.TableName).Where(cond).Limit(1).LoadStruct(this)
		if err != nil {
			logging.Error("finchina db: getCompcode: %s", err)
			return err
		}
		if this.COMPCODE.Valid == false {
			logging.Error("finchina db: getCompcode: Query COMPCODE is NULL by SYMBOL='%s'", TABLE_TQ_OA_STCODE, symbol)
			return ERROR_COMPCODE_NULL
		}

		data, _ := json.Marshal(this)
		if err := RedisCache.Setex(key, 2*60*60, data); err != nil {
			logging.Error("Redis cache %s TTL %d: %s", key, 2*60*60, err)
			return err
		}
		logging.Info("Redis cache %s TTL %d", key, 2*60*60)
		return nil
	}

	json.Unmarshal(v, this)
	return nil
}

func (this *TQ_OA_STCODE) GetCompcode(symbol interface{}) error {
	return this.getCompcode(symbol)
}
