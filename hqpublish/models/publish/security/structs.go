package security

const (
	SECURITY_CODE_LEN = 24 ///< 证券代码长度
	SECURITY_NAME_LEN = 40 ///< 证券名称长度
	SECURITY_DESC_LEN = 8  ///< 英文简称
	INDUSTRY_CODE_LEN = 8  ///< 行业代码
	SECURITY_ISIN_LEN = 16 ///< 证券国际代码信息
)

// 股票代码表结构(hq:st:name)
type TagSecurityName struct {
	NSID        int32
	NMarket     int32                   // 市场类型
	SzSType     [4]byte                 // 证券类型										len:4
	SzStatus    [4]byte                 // 证券状态										len:4
	SzSCode     [SECURITY_CODE_LEN]byte // 证券代码: 600036.SH							len:SECURITY_CODE_LEN
	SzSymbol    [SECURITY_CODE_LEN]byte // 证券原始: 600036								len:SECURITY_CODE_LEN
	SzISIN      [SECURITY_ISIN_LEN]byte // 证券国际代码信息								    len:SECURITY_ISIN_LEN
	SzSName     [SECURITY_NAME_LEN]byte // 证券名称 (超过24字节部分被省略)					len:SECURITY_NAME_LEN
	SzSCName    [SECURITY_NAME_LEN]byte // 证券简体中文名称 (美股、港股超过40字节部分被省略		len:SECURITY_NAME_LEN
	SzDESC      [SECURITY_DESC_LEN]byte // 英文简称										len:SECURITY_DESC_LEN
	SzPhonetic  [SECURITY_CODE_LEN]byte // 拼音											len:SECURITY_CODE_LEN
	SzCUR       [4]byte                 // 币种											len:4
	SzIndusCode [INDUSTRY_CODE_LEN]byte // 行业代码										len:INDUSTRY_CODE_LEN
}

// / 股票静态数据 (REDISKEY_SECURITY_STATIC)
type StockStatic struct {
	NSID               int32
	SzSType            [4]byte
	SzStatus           [4]byte
	NListDate          int32
	NLastTradeDate     int32
	NDelistDate        int32
	LlCircuShare       int64
	LlTotalShare       int64
	LlLast5Volume      int64
	NEPS               int32
	LlTotalProperty    int64
	LlFlowProperty     int64
	NAVPS              int32
	LlMainIncoming     int64
	LlMainProfit       int64
	LlTotalProfit      int64
	LlNetProfit        int64
	NHolders           int32
	NReportDate        int32
	NQuickMovingRatio  int32
	NCurrentRatio      int32
	NEUndisProfit      int32
	NFlowLiab          int64
	NTotalLiabilities  int64
	NTotalHolderEquity int64
	NCapitalReserve    int64
	NIncomeInvestments int64
}

func ByteNToString(src interface{}) string {
	var des []byte
	switch src.(type) {
	case [4]byte:
		for _, v := range src.([4]byte) {
			if v == 0 {
				break
			}
			des = append(des, v)
		}
	case [8]byte:
		for _, v := range src.([8]byte) {
			if v == 0 {
				break
			}
			des = append(des, v)
		}
	case [16]byte:
		for _, v := range src.([16]byte) {
			if v == 0 {
				break
			}
			des = append(des, v)
		}
	case [24]byte:
		for _, v := range src.([24]byte) {
			if v == 0 {
				break
			}
			des = append(des, v)
		}
	case [40]byte:
		for _, v := range src.([40]byte) {
			if v == 0 {
				break
			}
			des = append(des, v)
		}
	}
	return string(des)
}
