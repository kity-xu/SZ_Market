//攻击力度&攻击人气
package publish2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	. "haina.com/market/hqpublish/models"
	"haina.com/share/logging"
)

type GJLDRQ struct {
	NSID  int32
	Ntime int32
	NGjld int32
	NGjrq int32
}

type GJLDRQJson struct {
	NSID  int32 `json:"nsid"`
	Ntime int32 `json:"ntime"`
	NGjld int32 `json:"nGjld"`
	NGjrq int32 `json:"nGjrq"`
}

func NewGJLDRQ() *GJLDRQ {
	return &GJLDRQ{}
}

func (GJLDRQ) GetGJLDRQ(sid int32) (*[]GJLDRQJson, error) {
	key := fmt.Sprintf("hq:gjrq:min:%d", sid)
	list, err := RedisStore.LRange(key, 0, -1)
	if len(list) == 0 && err != nil {
		logging.Error("%s", err.Error())
		return nil, err
	}

	var ldrqs []GJLDRQJson
	for _, v := range list {
		ld := &GJLDRQ{}
		if err = binary.Read(bytes.NewBuffer([]byte(v)), binary.LittleEndian, ld); err != nil && err != io.EOF {
			logging.Error("%v", err)
			return nil, err
		}
		lj := GJLDRQJson{
			NSID:  ld.NSID,
			Ntime: ld.Ntime,
			NGjld: ld.NGjld,
			NGjrq: ld.NGjrq,
		}
		ldrqs = append(ldrqs, lj)
	}
	return &ldrqs, nil
	//
	//l := len(arr)
	//ldrq := &GJLDRQ{}
	//size := binary.Size(ldrq)
	//
	//var ldrqs []GJLDRQJson
	//for i := 0; i < l; i += size {
	//	ld := &GJLDRQ{}
	//	if err = binary.Read(bytes.NewBuffer(data[i:i+size]), binary.LittleEndian, ld); err != nil || err != io.EOF {
	//		logging.Error("%v", err)
	//		return nil, err
	//	}
	//	lj := GJLDRQJson{
	//		NSID:  ld.NSID,
	//		Ntime: ld.Ntime,
	//		NGjld: ld.NGjld,
	//		NGjrq: ld.NGjrq,
	//	}
	//	ldrqs = append(ldrqs, lj)
	//}
	//return &ldrqs, nil
}
