package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"haina.com/share/lib"
	"haina.com/share/logging"

	"github.com/gin-gonic/gin"
)

type klist struct {
	NSid     int64 `json:"nSID"`
	NTime    int64 `json:"nTime"`
	NPreCPx  int64 `json:"nPreCPx"`
	NOpenPx  int64 `json:"nOpenPx"`
	NHighPx  int64 `json:"nHighPx"`
	NLastPx  int64 `json:"nLastPx"`
	LlVolume int64 `json:"llVolume"`
	LlValue  int64 `json:"llValue"`
	NAvgPx   int64 `json:"nAvgPx"`
}

func GetJson(c *gin.Context) {

	url := "https://hbmk.0606.com.cn/api/hq/kline?format=json"
	//	post := "{\"SID\":100600000,\"Type\":11,\"TimeIndex\":0,\"Num\":2,\"Direct\":0}"

	posts := map[string]interface{}{
		"SID":       100600000,
		"Type":      11,
		"TimeIndex": 0,
		"Num":       2,
		"Direct":    0,
	}

	postjson, _ := json.Marshal(posts)

	logging.Info("type===%v", reflect.TypeOf(postjson))

	//	logging.Info("postjson===%v", string(postjson))

	//	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postjson)) // req.Header.Set("X-Custom-Header", "myvalue")
	//	req.Header.Set("Content-Type", "application/json")
	//	client := &http.Client{}
	//	resp, err := client.Do(req)
	//	if err != nil {
	//		panic(err)
	//	}

	//	defer resp.Body.Close()
	//	fmt.Println("response status", resp.Status)
	//	fmt.Println("response header", resp.Header)
	//	body, err := ioutil.ReadAll(resp.Body)
	//	if err != nil {
	//		//handle error
	//	}

	//	logging.Info(string(body))

	var _param struct {
		Code int64 `json:"code"`
		Data struct {
			SID   int64   `json:"sid"`
			Type  int64   `json:"type"`
			Total int64   `json:"total"`
			Begin int64   `json:"begin"`
			Num   int64   `json:"num"`
			KList []klist `json:"klist"`
		} `json:"data"`
	}

	errS := httpPostJson(url, posts, &_param)

	//errS := json.Unmarshal([]byte(body), &_param)
	if errS != nil {
		logging.Info(errS.Error())

		return
	} else {
		logging.Info("Code====%v", _param.Code)
		logging.Info("SID====%v", _param.Data.SID)
		logging.Info("SID====%v", _param.Data.KList)
	}
	lib.WriteString(c, 200, "")
}

func httpPostJson(url string, postJson interface{}, _param interface{}) error {
	postJsons, _ := json.Marshal(postJson)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postJsons))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println("response status", resp.Status)
	fmt.Println("response header", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(body), _param)
	return err
}
