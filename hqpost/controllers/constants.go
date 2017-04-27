package controllers

//Exchange 交易所
const (
	SH = "/sh/"
	SZ = "/sz/"
)

//// K线历史部分
const (
	REDISKEY_SECURITY_HDAY   = "hq:st:hday:%d"   ///<证券历史日K线数据(参数：sid) (hq-post写入)
	REDISKEY_SECURITY_HWEEK  = "hq:st:hweek:%d"  ///<证券周K线(参数：sid)
	REDISKEY_SECURITY_HMONTH = "hq:st:hmonth:%d" ///<证券月K线(参数：sid)
	REDISKEY_SECURITY_HYEAR  = "hq:st:hyear:%d"  ///<证券年K线(参数：sid)
	REDISKEY_SECURITY_HMIN   = "hq:st:hmin:%d"   ///<<证券历史分钟线数据(参数：sid) (hq-post写入)
	REDISKEY_SECURITY_HMIN5  = "hq:st:hmin5:%d"  ///<证券5分钟K线(参数：sid)
	REDISKEY_SECURITY_HMIN15 = "hq:st:hmin15:%d" ///<证券15分钟K线(参数：sid)
	REDISKEY_SECURITY_HMIN30 = "hq:st:hmin30:%d" ///<证券30分钟K线(参数：sid)
	REDISKEY_SECURITY_HMIN60 = "hq:st:hmin60:%d" ///<证券60分钟K线(参数：sid)
)

//个股信息
type StockSingle struct {
	SID    int32  // 证券ID
	Time   int32  // 时间 unix time
	PreCPx int32  // 昨收价 * 10000
	OpenPx int32  // 开盘价 * 10000
	HighPx int32  // 最高价 * 10000
	LowPx  int32  // 最低价 * 10000
	LastPx int32  // 最新价 * 10000
	Volume int64  // 成交量
	Value  int64  // 成交额 * 10000
	AvgPx  uint32 // 平均价 * 10000

}
