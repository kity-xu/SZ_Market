// 快照
package publish

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	//	. "haina.com/market/hqpublish/models"
	. "haina.com/share/models"

	pro "ProtocolBuffer/projects/hqpublish/go/protocol"

	"haina.com/market/hqpublish/models/memdata"
	hsgrr "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
)

var _ = fmt.Println
var _ = hsgrr.Bytes
var _ = logging.Info
var _ = bytes.NewBuffer
var _ = binary.Read
var _ = io.ReadFull

type Snapshot struct {
	Model `db:"-"`
}

func NewSnapshot() *Snapshot {
	return &Snapshot{
		Model: Model{
			CacheKey: REDISKEY_SECURITY_SNAP,
		},
	}
}

// 获取快照
func (this *Snapshot) GetSnapshot(req *pro.RequestSnapshot) (*pro.PayloadStockSnapshot, *pro.PayloadIndexSnapshot, error) {
	info := memdata.HandleSwapInfo(req.SID)
	if info == nil {
		return nil, nil, errors.New(fmt.Sprintf("SID %d Snapshot data no found", req.SID))
	}

	var ferr error

	switch t := info.Type[1]; t {
	case 'S': // 股票
		stock, err := NewStockSnapshot().GetStockSnapshotObj(req)
		if err == nil {
			return stock, nil, nil
		}
		ferr = err
	case 'I': // 指数
		indexob, err := NewIndexSnapshot().GetIndexSnapshotObj(req)
		if err == nil {
			index := pro.PayloadIndexSnapshot{
				SID:      req.SID,
				SnapInfo: indexob,
			}
			return nil, &index, nil
		}
		ferr = err
	default:
		err := errors.New(fmt.Sprintf("Unknown Type %v", t))
		logging.Warning("Unknown Type %v", err)
		return nil, nil, err
	}

	return nil, nil, ferr
}
