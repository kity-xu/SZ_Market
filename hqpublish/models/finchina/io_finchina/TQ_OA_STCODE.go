// 证券内码表
package io_finchina

import (
	"database/sql"
	"fmt"

	. "haina.com/market/hqpublish/models"
	. "haina.com/share/models"

	redigo "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"

	"haina.com/market/hqpublish/models/finchina"
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
			TableName: finchina.TABLE_TQ_OA_STCODE,
			Db:        MyCat,
		},
	}
}

func (this *TQ_OA_STCODE) getCompcode(symbol string, market string) error {
	switch market {
	case "100", "200":
	default:
		return finchina.ErrMarket
	}
	seg := fmt.Sprintf("%s%s", market, symbol)
	key := fmt.Sprintf(finchina.REDIS_SYMBOL_COMPCODE, seg)
	v, err := RedisCache.Get(key)
	if err != nil {
		if err != redigo.ErrNil {
			logging.Error("Redis get %s: %s", key, err)
		}

		var cond string
		switch market {
		case "100": // 001002 上海证券交易所
			cond = "EXCHANGE='001002'"
		case "200": // 001003 深圳证券交易所
			cond = "EXCHANGE='001003'"
		}
		cond += " and SETYPE='101' and SYMBOL=" + symbol

		err = this.Db.Select("*").From(this.TableName).Where(cond).Limit(1).LoadStruct(this)
		if err != nil {
			logging.Error("finchina db: getCompcode: %s", err)
			return err
		}
		if this.COMPCODE.Valid == false {
			logging.Error("finchina db: getCompcode: Query COMPCODE is NULL by SYMBOL='%s'", finchina.TABLE_TQ_OA_STCODE, symbol)
			return finchina.ErrNullComp
		}
		if err := RedisCache.Setex(key, finchina.REDIS_TTL, []byte(this.COMPCODE.String)); err != nil {
			logging.Error("Redis cache %s TTL %d: %s", key, finchina.REDIS_TTL, err)
			return err
		}
		logging.Info("Redis cache %s TTL %d", key, finchina.REDIS_TTL)
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
