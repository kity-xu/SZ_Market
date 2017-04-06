// F10 财务分析接口共用
package finchina

import (
	"database/sql"
	"fmt"

	. "haina.com/share/models"

	redigo "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

// ---------------------------------------------------------------------
type SymbolToCompcode struct {
	Model    `db:"-"`
	COMPCODE dbr.NullString //公司代码(公司内码) 通过 SYMBOL 得到
}

func NewSymbolToCompcode() *SymbolToCompcode {
	return &SymbolToCompcode{
		Model: Model{
			TableName: TABLE_TQ_OA_STCODE,
			Db:        MyCat,
		},
	}
}

func (this *SymbolToCompcode) getCompcode(symbol string) error {
	key := fmt.Sprintf(REDIS_SYMBOL_COMPCODE, symbol)
	v, err := redis.Get(key)
	if err != nil {
		if err != redigo.ErrNil {
			logging.Error("Redis get %s: %s", key, err)
		}
		err := this.Db.Select("COMPCODE").From(this.TableName).Where("SYMBOL=?", symbol).Limit(1).LoadStruct(this)
		if err != nil {
			logging.Error("finchina db: getCompcode: %s", err)
			return err
		}
		if this.COMPCODE.Valid == false {
			logging.Error("finchina db: getCompcode: select COMPCODE from %s where SYMBOL='%s'", TABLE_TQ_OA_STCODE, symbol)
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
