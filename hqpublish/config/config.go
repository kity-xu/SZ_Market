package config

import (
	"os"
	"os/exec"
	"path/filepath"

	"haina.com/share/lib"
)

var config *AppConfig

type AppSettings struct {
	AllowOrigin string     `xml:"allowOrigin"`
	EncryFactor string     `xml:"encryFactor"`
	Environment string     `xml:"environment"`
	Listen      string     `xml:"listen"`
	Projects    []Projects `xml:"projects"`
}

type Projects struct {
	AppId      string `xml:"appId"`
	ConfigFile string `xml:"configFile"`
}

type AccessKeys struct {
	ID     string `xml:"id"`
	Secret string `xml:"secret"`
	AESKey string `xml:"aesKey"`
}

// Database
type Database struct {
	DriverName string `xml:"driverName"`
	DataSource string `xml:"dataSource"`
}

// MongoDB
type MongoStore struct {
	Source string `xml:"source"`
}

// Redis
type RedisStore struct {
	Addr    string `xml:"addr"`
	Auth    string `xml:"auth"`
	Db      string `xml:"db"`
	Timeout int    `xml:"timeout"`
}

// Email
type EmailSetting struct {
	Addr     string `xml:"addr"`
	Password string `xml:"password"`
	Server   string `xml:"server"`
	Port     string `xml:"port"`
}

type SessionSetting struct {
	On           bool   `xml:"on"`
	ProviderName string `xml:"providerName"`
	Config       string `xml:"config"`
}
type MnsSetting struct {
	Url             string    `xml:"url"`
	AccessKeyId     string    `xml:"accessKeyId"`
	AccessKeySecret string    `xml:"accessKeySecret"`
	Queues          QueueName `xml:"queues"`
}
type QueueName struct {
	SmartCall string `xml:"exchangeSmartCall"`
}

type LogServer struct {
	On   bool   `xml:"on"`
	Addr string `xml:"addr"`
	Port string `xml:"port"`
}

type AppConfig struct {
	AccessKeys AccessKeys     `xml:"accessKeys"`
	Cors       CorsSetting    `xml:"cors"`
	Db         Database       `xml:"database"`
	Email      EmailSetting   `xml:"emailSetting"`
	Mns        MnsSetting     `xml:"mns"`
	Mongo      MongoStore     `xml:"mongoStore"`
	Redis      RedisStore     `xml:"redisStore"`
	RedisCache RedisStore     `xml:"redisCache"`
	Serve      ListenAndServe `xml:"listenAndServe"`
	Session    SessionSetting `xml:"session"`
	Settings   AppSettings    `xml:"appSettings"`
	Log        LogServer      `xml:"logServer"`
}

type CorsSetting struct {
	AllowOrigin []string `xml:"allowOrigin"`
}

type ListenAndServe struct {
	Port    string `xml:"port"`
	LogPort string `xml:"logport"`
}

func Default(appID string) *AppConfig {
	if config == nil {
		var cfg AppConfig
		lib.LoadConfig(appID, &cfg)
		config = &cfg
	}
	return config
}

func Reload() {
	config = nil
}

// ------------------------------------------------------------------------

func getCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	return filepath.Dir(file)
}
