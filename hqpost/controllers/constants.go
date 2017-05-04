package controllers

const (
	/// 证券代码表
	REDISKEY_MARKET_SECURITY_TABLE = "hq:market:sts:%d" ///<证券市场代码表(参数：MarketID)  (hq-init写入)
	REDISKEY_SECURITY_NAME_ID      = "hq:st:name:%d"    ///<证券代码(参数：sid) (hq-init写入)
	REDISKEY_SECURITY_NAME_CODE    = "hq:st:name:%s"    ///<证券代码(参数：scode) (hq-init写入)
)
