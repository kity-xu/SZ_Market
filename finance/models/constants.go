package models

// App Setting
//---------------------------------------------------------------------------------
const (
	APP_NAME    = "market_finance"
	APP_VERSION = "0.0.3.0"
	APP_PID     = "market_finance"
)

// Context and Session Categorys
// --------------------------------------------------------------------------------

const (
	CONTEXT_END_DATE = "enddate" // 起始时间（默认当前时间）
	CONTEXT_COUNT    = "count"   // 条数（默认10条）
	CONTEXT_SCODE    = "scode"   // 证券代码
)

// F10财务分析接口URL请求相关参数
const (
	CONTEXT_TYPE    = "type"    // 报表类型
	CONTEXT_PERPAGE = "perpage" // 每页条数
	CONTEXT_PAGE    = "page"    // 当前页码
)
