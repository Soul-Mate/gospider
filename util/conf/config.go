package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type FilterConfig struct {
	Usage    bool        `json:"usage"`
	Length   interface{} `json:"length"`
	Name     string      `json:"name"`
	Dump     bool        `json:"dump"`
	DumpPath string      `json:"dump_path"`
}

type HttpConfig struct {
	Timeout int `json:"client_timeout"` // request timeout
}

type DownloadConfig struct {
	Concurrent int         `json:"concurrent"` // download concurrent
	Delay      interface{} `json:"delay"`      // download delay time
}

type ScheduleConfig struct {
	Name string `json:"name"`
}

type RedisConfig struct {
	Addr           string `json:"addr"`
	DB             int    `json:"db"`
	ConnectionNums int    `json:"connection_nums"`
}

type LogConfig struct {
	Enable bool   `json:"enable"`
	Level  string `json:"level"`
	File   string `json:"file"`
}

type ConnectionConfig struct {
	Redis RedisConfig `json:"redis"`
}

type SpiderConfig struct {
	Http       HttpConfig       `json:"http"`
	Downloader DownloadConfig   `json:"downloader"`
	Items      Items            `json:"items"`
	Filter     FilterConfig     `json:"filter"`
	Schedule   ScheduleConfig   `json:"schedule"`
	Connection ConnectionConfig `json:"connection"`
	Log        LogConfig        `json:"log"`
}

var (
	http = HttpConfig{
		Timeout: 5,
	}

	download = DownloadConfig{
		Concurrent: 100,
		Delay:      5,
	}

	scheduleConfig = ScheduleConfig{
		Name: "local",
	}

	connectionConfig = ConnectionConfig{
		Redis: RedisConfig{
			Addr: "127.0.0.1:6379",
			DB:   1,
		},
	}

	defaultConfig = &SpiderConfig{
		Http:       http,
		Downloader: download,
		Items:      nil,
		Schedule:   scheduleConfig,
		Connection: connectionConfig,
	}

	GlobalSharedConfig *SpiderConfig
)

func InitConfig(path string) {
	if path == "" {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		path = dir + "/config.json"
	}
	_, err := os.Stat(path)

	if err != nil {
		log.Fatal(err)
	}

	if os.IsNotExist(err) {
		GlobalSharedConfig = defaultConfig
		return
	}

	stream, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}
	context, err := ioutil.ReadAll(stream)

	if err != nil {
		log.Fatal(err)
	}

	conf := new(SpiderConfig)

	if err = json.Unmarshal(context, conf); err != nil {
		log.Fatal(err)
	}

	GlobalSharedConfig = conf
}
