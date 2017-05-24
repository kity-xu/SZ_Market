package controllers

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"

	"haina.com/share/lib"
	"haina.com/share/logging"

	"ProtocolBuffer/projects/hqpublish/go/protocol"
)

var _ = protocol.HAINA_PUBLISH_CMD_MUST_BE_NULL

// 接收 http Body 数据
func GetRequestData(c *gin.Context, bufsize int) ([]byte, error) {
	temp := make([]byte, bufsize)
	n, err := c.Request.Body.Read(temp)
	if err != nil && err != io.EOF {
		logging.Error("Body Read: %v", err)
		return nil, err
	}
	//logging.Info("\x1b[0;31mRequest Body len %d\x1b[0m", n)
	logging.Info("Request Body len %d", n)
	return temp[:n], nil
}

// 接收格式为Json的http包体并解析
func RecvAndUnmarshalJson(c *gin.Context, bufsize int, v interface{}) (int, error) {
	buf, err := GetRequestData(c, bufsize)
	if err != nil && err != io.EOF {
		return 40004, err
	}
	if err := json.Unmarshal(buf, v); err != nil {
		return 40004, err
	}
	logging.Info("Request Data: %+v", v)
	return 0, nil
}

// 接收格式为PB的http包体并解析
func RecvAndUnmarshalPB(c *gin.Context, bufsize int, pb proto.Message) (int, error) {
	buf, err := GetRequestData(c, bufsize)
	if err != nil && err != io.EOF {
		return 40004, err
	}
	if err := proto.Unmarshal(buf, pb); err != nil {
		return 40004, err
	}
	logging.Info("Request Data: %+v", pb)
	return 0, nil
}

// 打包返回数据(通用)
// c  Code for Response
// t  Type for protocol.HAINA_PUBLISH_CMD
// bs Protocol Buffer encoded []byte
func MakeRespDataByBytes(c int, t protocol.HAINA_PUBLISH_CMD, bs []byte) ([]byte, error) {
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

	if ret_code == 200 && bs != nil {
		if err := binary.Write(buffer, binary.LittleEndian, bs); err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}

// 构造返回数据(protobuf)
// c  Code for Response
// t  Type for protocol.HAINA_PUBLISH_CMD
// pb Protocol Buffer's Message struct
func MakeRespDataByPB(c int, t protocol.HAINA_PUBLISH_CMD, pb proto.Message) ([]byte, error) {
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

func WriteJson(c *gin.Context, code int, data interface{}) {
	lib.WriteString(c, code, data)
}
func WriteDataJson(c *gin.Context, jsonb []byte) {
	c.Data(200, "application/json; charset=utf-8", jsonb)
}
func WriteDataBytes(c *gin.Context, _type protocol.HAINA_PUBLISH_CMD, payload []byte) {
	data, err := MakeRespDataByBytes(200, _type, payload)
	if err != nil {
		logging.Error("MakeRespDataByBytes: %v", err)
		return
	}
	lib.WriteData(c, data)
}
func WriteDataPB(c *gin.Context, _type protocol.HAINA_PUBLISH_CMD, pb proto.Message) {
	data, err := MakeRespDataByPB(200, _type, pb)
	if err != nil {
		logging.Error("MakeRespDataByPB: %v", err)
		return
	}
	lib.WriteData(c, data)
}
func WriteDataErrCode(c *gin.Context, code int) {
	data, err := MakeRespDataByBytes(code, 0, nil)
	if err != nil {
		logging.Error("MakeRespDataByBytes: %v", err)
		return
	}
	lib.WriteData(c, data)
}
