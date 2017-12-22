package publish2

import (
	"encoding/json"

	"haina.com/share/logging"

	. "haina.com/market/hqpublish/models"
)

func GetResFromCache(key string, res interface{}) (interface{}, error) {
	data, err := RedisCache.GetBytes(key)
	if err != nil {
		logging.Debug("RedisCache: RedisCache get error |%v", err)
		return nil, err
	}

	if err = json.Unmarshal(data, res); err != nil {
		logging.Debug("RedisCache: json.Unmarshal error | %v", err)
		return nil, err
	}
	return res, nil
}

func SetResToCache(key string, src interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		logging.Debug("RedisCache: json Marshal error| %v", err)
	}

	if err = RedisCache.Setex(key, TTL.Min60, data); err != nil {
		logging.Debug("%v", err)
		return err
	}
	return nil
}
