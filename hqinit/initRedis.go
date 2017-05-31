package main

import (
	"fmt"

	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

func ReadAndWriteRedis() {
	Redis171 := redis.NewRedisPool("47.93.102.171:61380", "0", "8dc40c2c4598ae5a", 3)
	Redis200 := redis.NewRedisPool("192.168.18.200:6379", "0", "8dc40c2c4598ae5a", 3)

	codes, err := Redis171.LRange("hq:st:nsid", 0, -1)
	if err != nil {
		logging.Error("%v", codes)
		return
	}

	for _, code := range codes {
		key := fmt.Sprintf("hq:st:min:%s", code)

		ls, _ := redis.LRange(key, 0, -1)

		if len(ls) < 1 {
			continue
		}

		for _, v := range ls {
			Redis200.Lpush(key, []byte(v))
		}

	}

}
