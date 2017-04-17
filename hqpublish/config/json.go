package config

import (
	"encoding/json"
	io "io/ioutil"
	//"github.com/bitly/go-simplejson"
)

type JsonConfig struct {
}

func NewJsonConfig() *JsonConfig {
	return &JsonConfig{}
}

func (self *JsonConfig) Load(filename string, v interface{}) error {
	data, err := io.ReadFile(filename)
	if err != nil {
		return err
	}
	datajson := []byte(data)
	err = json.Unmarshal(datajson, &v)
	if err != nil {
		return err
	}
	return nil
}
