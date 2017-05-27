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

////////////////////////////////////////////////////////////////////////////////
// 接收 http Body 数据
// 返回：
//   []byte 接受到的数据，最大尺寸 bufsize
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

////////////////////////////////////////////////////////////////////////////////
// 接收格式为Json的http包体并解析, v 为 protocol.RequestXXX 地址
// 返回：
//   int 错误码
func RecvAndUnmarshalJson(c *gin.Context, bufsize int, v interface{}) (int, error) {
	buf, err := GetRequestData(c, bufsize)
	if err != nil && err != io.EOF {
		return 40004, err
	}
	if err := json.Unmarshal(buf, v); err != nil {
		return 40004, err
	}
	logging.Info("Parsed Data %+v", v)
	return 0, nil
}

////////////////////////////////////////////////////////////////////////////////
// 接收格式为PB的http包体并解析, pb 为 protocol.RequestXXX 地址
// 返回：
//   int 错误码
func RecvAndUnmarshalPB(c *gin.Context, bufsize int, pb proto.Message) (int, error) {
	buf, err := GetRequestData(c, bufsize)
	if err != nil && err != io.EOF {
		return 40004, err
	}
	if err := proto.Unmarshal(buf, pb); err != nil {
		return 40004, err
	}
	logging.Info("Parsed Data: %+v", pb)
	return 0, nil
}

////////////////////////////////////////////////////////////////////////////////
// pb模式应答数据封包(pb模式通用)
////////////////////////////////////////////////////////////////////////////////
// 封包后数据包格式 Header(Code + Type) + Payload([]byte)
// Header   格式 8 Byte: 4 Byte 小端整数 Code + 4 Byte 小端整数 Type
// Payload  格式 n Byte: n >= 0, Protocol Buffer序列化编码)
//   c  Code for Response
//   t  Type for protocol.HAINA_PUBLISH_CMD
//   bs Protocol Buffer encoded []byte
// 返回：
//   []byte 已经完成封包的数据
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

////////////////////////////////////////////////////////////////////////////////
// pb模式应答数据封包
//   c  Code for Response
//   t  Type for protocol.HAINA_PUBLISH_CMD
//   pb Protocol Buffer's Message struct
// 返回：
//   []byte 已经完成封包的数据
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

////////////////////////////////////////////////////////////////////////////////
// json模式打包返回数据
//   c    Code for Response
//   obj  Object for data
// 返回：
//   []byte 已经完成序列化的json数据
func MakeRespJson(c int, obj interface{}) ([]byte, error) {
	res := map[string]interface{}{"code": c}
	if obj != nil {
		res["data"] = obj
	}
	js, err := json.Marshal(&res)
	if err != nil {
		return nil, err
	}
	return js, nil
}

////////////////////////////////////////////////////////////////////////////////
// json模式应答数据http方式发送
// bin 为 head(8 Byte) + payload(pb []byte)
func WriteDataBinary(c *gin.Context, bin []byte) {
	lib.WriteData(c, bin)
}

////////////////////////////////////////////////////////////////////////////////
// json模式序列化数据http方式发送
//   jsb/jss 为含 code 的JSON序列化数据
func WriteDataJson(c *gin.Context, jsb []byte) {
	c.Data(200, "application/json; charset=utf-8", jsb)
}
func WriteJsonBytes(c *gin.Context, jsb []byte) {
	WriteDataJson(c, jsb)
}
func WriteJsonString(c *gin.Context, jss string) {
	WriteDataJson(c, []byte(jss))
}

////////////////////////////////////////////////////////////////////////////////
// json模式应答数据序列化并http方式发送
//   code JSON code 返回码
//   data JSON data 对象
func WriteJson(c *gin.Context, code int, data interface{}) {
	lib.WriteString(c, code, data)
}

////////////////////////////////////////////////////////////////////////////////
// pb模式应答数据封包并http方式发送
//   _type   : payload 类型
//   payload : protocol.PayloadXXXX -> proto.Marshal 编码后返回的 []byte
func WriteDataBytes(c *gin.Context, _type protocol.HAINA_PUBLISH_CMD, payload []byte) {
	data, err := MakeRespDataByBytes(200, _type, payload)
	if err != nil {
		logging.Error("MakeRespDataByBytes: %v", err)
		return
	}
	lib.WriteData(c, data)
}

////////////////////////////////////////////////////////////////////////////////
// pb模式应答数据封包并http方式发送
//   _type : pb 类型
//   pb    : protocol.PayloadXXXX
func WriteDataPB(c *gin.Context, _type protocol.HAINA_PUBLISH_CMD, pb proto.Message) {
	data, err := MakeRespDataByPB(200, _type, pb)
	if err != nil {
		logging.Error("MakeRespDataByPB: %v", err)
		return
	}
	lib.WriteData(c, data)
}

////////////////////////////////////////////////////////////////////////////////
// pb模式应答数据封包并http方式发送
// 数据包格式 Header(8 Byte), 没有 Payload
func WriteDataErrCode(c *gin.Context, code int) {
	data, err := MakeRespDataByBytes(code, 0, nil)
	if err != nil {
		logging.Error("MakeRespDataByBytes: %v", err)
		return
	}
	lib.WriteData(c, data)
}
