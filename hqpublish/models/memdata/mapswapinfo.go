package memdata

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"

	ctrl "haina.com/market/hqpublish/controllers"
	"haina.com/market/hqpublish/models/publish/security"
	//. "haina.com/share/models"

	hsgrr "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
	"haina.com/share/store/redis"

	"ProtocolBuffer/projects/hqpublish/go/protocol"
)

var (
	//_ = MapToStruct
	_ = redis.Init
	//_ = GetCache
	_ = ctrl.MakeRespDataByBytes
	_ = errors.New
	_ = fmt.Println
	_ = hsgrr.Bytes
	_ = logging.Info
	_ = bytes.NewBuffer
	_ = binary.Read
	_ = io.ReadFull
)

type SwapInfo struct {
	Type   string
	Status string
}

type SwapInfoMap map[int32]*SwapInfo

//------------------------------------------------------------------------------
var (
	swapInfoMap      = make(SwapInfoMap)
	swapInfoMapMutex sync.RWMutex
)

func SwapInfoMap_Wlock() {
	swapInfoMapMutex.Lock()
}
func SwapInfoMap_Wunlock() {
	swapInfoMapMutex.Unlock()
}
func SwapInfoMap_Rlock() {
	swapInfoMapMutex.RLock()
}
func SwapInfoMap_Runlock() {
	swapInfoMapMutex.RUnlock()
}

////////////////////////////////////////////////////////////////////////////////

func HandleSwapInfo(sid int32) *SwapInfo {
	info := SafeGetSwapInfo(sid)
	if info == nil {
		return nil
	}
	//SafeSetSwapInfo(sid, info)
	return info
}

func GetSwapInfo(sid int32) *SwapInfo {
	find, ok := swapInfoMap[sid]
	if ok {
		return find
	}

	info := SwapInfo{}

	req := protocol.RequestSingleSecurity{
		SID: sid,
	}
	pay, err := security.NewSecurityInfo().GetSecurityBasicInfo(&req)
	if err != nil {
		logging.Error("%v", err)
		return nil
	} else {
		info.Status = pay.SNInfo.SzStatus
		info.Type = pay.SNInfo.SzSType
	}

	return &info
}

func SafeGetSwapInfo(sid int32) *SwapInfo {
	SwapInfoMap_Rlock()
	defer SwapInfoMap_Runlock()
	return GetSwapInfo(sid)
}

func SetSwapInfo(sid int32, info *SwapInfo) *SwapInfo {
	if info == nil {
		return nil
	}
	swapInfoMap[sid] = info
	return info
}
func SafeSetSwapInfo(sid int32, info *SwapInfo) *SwapInfo {
	SwapInfoMap_Wlock()
	defer SwapInfoMap_Wunlock()
	return SetSwapInfo(sid, info)
}

func EmptySwapInfoMap() {
	// 方法一：手工遍历
	//	for key := range swapInfoMap {
	//		delete(swapInfoMap, key)
	//	}

	// 方法二：利用GC
	swapInfoMap = make(SwapInfoMap)
}

func SafeEmptySwapInfoMap() {
	SwapInfoMap_Wlock()
	defer SwapInfoMap_Wunlock()
	EmptySwapInfoMap()
}
