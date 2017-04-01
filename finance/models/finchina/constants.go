package finchina

const (
	TABLE_TQ_OA_STCODE     = "TQ_OA_STCODE"     //证券内码表
	TABLE_TQ_SK_DIVIDENTS  = "TQ_SK_DIVIDENTS"  //分红情况表
	TABLE_TQ_SK_PROADDISS  = "TQ_SK_PROADDISS"  //上市公司增发情况表
	TABLE_TQ_SK_PROPLACING = "TQ_SK_PROPLACING" //上市公司配股情况表
	TABLE_TQ_SK_SHAREHDCHG = "TQ_SK_SHAREHDCHG"
)

// market SCHEMA
const (
	TABLE_TQ_SK_SHAREHOLDERNUM = "TQ_SK_SHAREHOLDERNUM" // 股东户数统计表
	TABLE_TQ_SK_OTSHOLDER      = "TQ_SK_OTSHOLDER"      // 流通股东信息表
	TABLE_TQ_SK_SHAREHOLDER    = "TQ_SK_SHAREHOLDER"    // 股东名单信息表
	TABLE_TQ_SK_SHARESTRUCHG   = "TQ_SK_SHARESTRUCHG"   // 股本结构变化
	TABLE_TQ_SK_LCPERSON       = "TQ_SK_LCPERSON"       // 上市公司董事名单
)

// Context and Session Categorys
// --------------------------------------------------------------------------------

const (
	CONTEXT_END_DATE     = "enddate"      // 起始时间（默认当前时间）
	CONTEXT_COUNT        = "count"        // 条数（默认10条）
	CONTEXT_SECURITYCODE = "securitycode" // 证卷代码
)
