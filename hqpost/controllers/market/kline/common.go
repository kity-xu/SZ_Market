package kline

import "errors"

/*****************************************Const***************************************/
//Exchange 交易所
const (
	SH = "/sh/"
	SZ = "/sz/"
)

//// K线历史部分
const (
	REDISKEY_SECURITY_MIN    = "hq:st:min:%d"    ///<证券分钟线数据(参数：sid) (calc写入)
	REDISKEY_SECURITY_HDAY   = "hq:st:hday:%d"   ///<证券历史日K线数据(参数：sid) (hq-post写入)
	REDISKEY_SECURITY_HWEEK  = "hq:st:hweek:%d"  ///<证券周K线(参数：sid)
	REDISKEY_SECURITY_HMONTH = "hq:st:hmonth:%d" ///<证券月K线(参数：sid)
	REDISKEY_SECURITY_HYEAR  = "hq:st:hyear:%d"  ///<证券年K线(参数：sid)
)

///
var (
	NOTFOUND_DAYLINE_IN_HGSFILE = errors.New("not found dayline in hgs filestore")
)

//----------------------------------------------------------------funtions--------------------------------------------------------------------//
