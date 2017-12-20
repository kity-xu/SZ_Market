// 证券内码表
package finchina

import (
	"database/sql"
	"fmt"

	. "haina.com/market/hqpublish/models"
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
	key := fmt.Sprintf(REDIS_SYMBOL_COMPCODE, seg)
	v, err := RedisCache.Get(key)
	if err != nil {
		if err != redigo.ErrNil {
			logging.Error("Redis get %s: %s", key, err)
		}

		var cond string
		cond = "EXCHANGE in ('001003','001002')"
		symstr := "0"
		if len(seg) > 6 {
			symstr = seg[3:]
		}
		cond += " and SETYPE='101' and SYMBOL=" + symstr

		err = this.Db.Select("*").From(this.TableName).Where(cond).Limit(1).LoadStruct(this)
		if err != nil {
			logging.Error("finchina db: getCompcode: %s", err)
			return err
		}
		if this.COMPCODE.Valid == false {
			logging.Error("finchina db: getCompcode: Query COMPCODE is NULL by SYMBOL='%s'", TABLE_TQ_OA_STCODE, symbol)
			return ErrNullComp
		}
		if err := RedisCache.Setex(key, REDIS_TTL, []byte(this.COMPCODE.String)); err != nil {
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

func (this *TQ_OA_STCODE) GetCompcode(symbol interface{}) error {
	return this.getCompcode(symbol)
}
