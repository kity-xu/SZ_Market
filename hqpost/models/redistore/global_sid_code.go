package redistore

import (
	"fmt"
	"strconv"

	"haina.com/share/logging"

	"haina.com/market/hqpost/models"

	"ProtocolBuffer/projects/hqpost/go/protocol"

	"github.com/golang/protobuf/proto"
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

type GlobalSid struct {
	Model `db:"-"`
}

func NewGlobalSid(key string) *GlobalSid {
	return &GlobalSid{
		Model: Model{
			CacheKey: key,
		},
	}
}

func (this *GlobalSid) GetGlobalSidFromRedis() (*[]int32, error) {
	sids, err := redis.LRange(this.CacheKey, 0, -1)
	if err != nil {
		return nil, err
	}
	if len(sids) < 1 {
		return nil, fmt.Errorf("sids list is null...")
	}

	var NSids []int32
	for _, sid := range sids {
		nsid, err := strconv.Atoi(sid)
		if err != nil { //此处出错误是因为出现了非数字字符
			return nil, fmt.Errorf("The sid is not numeric types...%s", sid)
		}
		NSids = append(NSids, int32(nsid))
	}
	return &NSids, nil
}

func GetSecurityStatus(sid int32) int {
	key := fmt.Sprintf(models.REDISKEY_SECURITY_NAME_ID, sid)
	data, err := models.RedisStore.GetBytes(key)
	if err != nil {
		logging.Error(err.Error())
		return -1
	}

	stock := &protocol.SecurityName{}
	if err = proto.Unmarshal(data, stock); err != nil {
		logging.Error(err.Error())
		return -1
	}

	if stock.SzSType[1] == 'S' { //股票
		return 'S'
	} else if stock.SzSType[1] == 'I' { //指数
		return 'I'
	} else {
		return 0
	}
}
