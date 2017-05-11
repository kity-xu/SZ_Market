// 证券内码表
package finchina

import (
	"database/sql"
	"fmt"
	"strings"

	. "haina.com/share/models"

	redigo "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

// TQ_OA_STCODE    证券内码表
// ---------------------------------------------------------------------
type TQ_OA_STCODE struct {
	Model    `db:"-"`
	COMPCODE dbr.NullString //公司代码(公司内码) 通过 SYMBOL 得到
}

func NewTQ_OA_STCODE() *TQ_OA_STCODE {
	return &TQ_OA_STCODE{
		Model: Model{
			TableName: TABLE_TQ_OA_STCODE,
			Db:        MyCat,
		},
	}
}

func (this *TQ_OA_STCODE) getCompcode(symbol string, market string) error {
	m := strings.ToUpper(market)
	switch m {
	case "SH", "SZ":
	default:
		return ErrMarket
	}
	seg := fmt.Sprintf("%s.%s", symbol, m)
	key := fmt.Sprintf(REDIS_SYMBOL_COMPCODE, seg)
	v, err := redis.Get(key)
	if err != nil {
		if err != redigo.ErrNil {
			logging.Error("Redis get %s: %s", key, err)
		}

		var cond string
		switch m {
		case "SH": // 001002 上海证券交易所
			cond = "EXCHANGE='001002'"
		case "SZ": // 001003 深圳证券交易所
			cond = "EXCHANGE='001003'"
		}
		cond += " and SETYPE='101' and SYMBOL=" + symbol

		err = this.Db.Select("*").From(this.TableName).Where(cond).Limit(1).LoadStruct(this)
		if err != nil {
			logging.Error("finchina db: getCompcode: %s", err)
			return err
		}
		if this.COMPCODE.Valid == false {
			logging.Error("finchina db: getCompcode: Query COMPCODE is NULL by SYMBOL='%s'", TABLE_TQ_OA_STCODE, symbol)
			return ErrNullComp
		}
		if err := redis.Setex(key, REDIS_TTL, []byte(this.COMPCODE.String)); err != nil {
			logging.Error("Redis cache %s TTL %d: %s", key, REDIS_TTL, err)
			return err
		}
		logging.Info("Redis cache %s TTL %d", key, REDIS_TTL)
		return nil
	}

	this.COMPCODE = dbr.NullString{
		NullString: sql.NullString{
			String: string(v),
			Valid:  true,
		},
	}

	return nil
}
func (this *TQ_OA_STCODE) GetCompcode(symbol string, market string) error {
	return this.getCompcode(symbol, market)
}
