package publish

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"

	. "haina.com/share/models"

	"ProtocolBuffer/projects/hqpublish/go/protocol"

	. "haina.com/market/hqpublish/models"
	"haina.com/share/logging"
)

type BlockIndex struct {
	Model `db:"-"`
}

func NewBlockIndex(redis_key string) *BlockIndex {
	return &BlockIndex{
		Model: Model{
			CacheKey: redis_key,
		},
	}
}

//protocol.PayloadBlockindex ------GetPayloadBlockindex

func (this *BlockIndex) GetPayloadBlockindex(req *protocol.RequestBlockindex) (*protocol.PayloadBlockindex, error) {
	table, err := ReadBlockIndexFile(FStore.Bindex, req.BlockID)
	if err != nil {
		logging.Error("%v", err.Error())
		return nil, err
	}

	return &protocol.PayloadBlockindex{
		BlockID: req.BlockID,
		Num:     int32(len(table)),
		KList:   table,
	}, nil
}

func ReadBlockIndexFile(path string, bid int32) ([]*protocol.KInfo, error) {
	logging.Debug("path:%v", path)
	filepath := fmt.Sprintf("%s%d", path, bid)

	data, err := ioutil.ReadFile(filepath)
	if err != nil || len(data) == 0 {
		logging.Error("read file error...")
		return nil, err
	}

	kinfo := &protocol.KInfo{}

	size := binary.Size(kinfo)

	table := make([]*protocol.KInfo, 0, 1024)
	for i := 0; i < len(data); i += size {
		index := &protocol.KInfo{}
		if err = binary.Read(bytes.NewBuffer(data[i:i+size]), binary.LittleEndian, index); err != nil && err != io.EOF {
			return nil, err
		}
		table = append(table, index)
	}
	return table, err
}
