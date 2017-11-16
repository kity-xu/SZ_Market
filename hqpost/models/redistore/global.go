package redistore

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"strings"

	"haina.com/share/logging"

	"haina.com/market/hqpost/models"

	"ProtocolBuffer/projects/hqpost/go/protocol"

	. "haina.com/share/models"
	"haina.com/share/store/redis"
)

const (
	SECURITY_CODE_LEN = 24 ///< 证券代码长度
	SECURITY_NAME_LEN = 40 ///< 证券名称长度
	SECURITY_DESC_LEN = 8  ///< 英文简称
	INDUSTRY_CODE_LEN = 8  ///< 行业代码
	SECURITY_ISIN_LEN = 16 ///< 证券国际代码信息
)

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

// GlobalSid ...
type GlobalSid struct {
	Model `db:"-"`
}

// NewGlobalSid ...
func NewGlobalSid(key string) *GlobalSid {
	return &GlobalSid{
		Model: Model{
			CacheKey: key,
		},
	}
}

// GetGlobalSidFromRedis 股票代码表
func (this *GlobalSid) GetGlobalSidFromRedis() (*[]int32, error) {
	keys, err := redis.Keys(this.CacheKey)
	if err != nil {
		return nil, err
	}
	if len(keys) < 1 {
		return nil, fmt.Errorf("keys list is null...")
	}

	var NSids []int32
	for _, key := range keys {
		sid := strings.Split(key, ":")[3]
		nsid, _ := strconv.Atoi(sid)
		NSids = append(NSids, int32(nsid))
	}
	return &NSids, nil
}

// GetSecurityBase 获取股票基本信息
func GetSecurityBase(sid int32) (*TagSecurityName, error) {
	key := fmt.Sprintf(models.REDISKEY_SECURITY_NAME_ID, sid)
	data, err := models.RedisStore.GetBytes(key)
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	stock := &TagSecurityName{}
	if err = binary.Read(bytes.NewBuffer(data), binary.LittleEndian, stock); err != nil && err != io.EOF {
		logging.Error(err.Error())
		return nil, err
	}
	//stock := &protocol.SecurityName{}
	//if err = proto.Unmarshal(data, stock); err != nil {
	//	logging.Error(err.Error())
	//	return nil, err
	//}
	return stock, nil
}

// GetSecurityMarketStatus ... 取单个市场状态
func GetSecurityMarketStatus(mid int32) (*protocol.MarketStatus, error) {
	key := fmt.Sprintf("hq:market:%d", mid)
	bin, err := models.RedisStore.GetBytes(key)
	if err != nil {
		return nil, err
	}
	var obj protocol.MarketStatus
	buffer := bytes.NewBuffer([]byte(bin))
	if err := binary.Read(buffer, binary.LittleEndian, &obj); err != nil && err != io.EOF {
		return nil, err
	}
	return &obj, nil
}

// TradeDateByMarketStatus ... 交易日
func TradeDateByMarketStatus(mid int32) int32 {
	market, err := GetSecurityMarketStatus(mid)
	if err != nil {
		logging.Error("GetSecurityMarketStatus Err | %v", err)
		return 0
	}
	return market.NTradeDate
}
