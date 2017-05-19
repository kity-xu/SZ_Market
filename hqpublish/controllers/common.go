package controllers

import (
	"bytes"
	"encoding/binary"

	"github.com/golang/protobuf/proto"
)

// 打包返回数据(通用)
// c Code
// t Type
// s Payload
func MakeRespDataByBytes(c int, t int, s []byte) ([]byte, error) {
	ret_code := int32(c)
	ret_type := int32(t)

	if ret_code != 200 {
		ret_type = 0
	}

	buffer := bytes.NewBuffer(nil)
	if err := binary.Write(buffer, binary.LittleEndian, &ret_code); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.LittleEndian, &ret_type); err != nil {
		return nil, err
	}
	if ret_code == 200 && s != nil {
		if err := binary.Write(buffer, binary.LittleEndian, s); err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}

// 构造返回数据(protobuf)
// c  Code
// t  Type
// pb Payload Protobuf Message
func MakeRespDataByPB(c int, t int, pb proto.Message) ([]byte, error) {
	var s []byte
	var e error
	if c == 200 && pb != nil {
		s, e = proto.Marshal(pb)
		if e != nil {
			return nil, e
		}
	}
	return MakeRespDataByBytes(c, t, s)
}
