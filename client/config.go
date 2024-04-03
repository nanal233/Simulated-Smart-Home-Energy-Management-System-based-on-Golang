package main

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type ConfigClient struct {
	ID          string `toml:"id"`           // 客户端编号
	Type        int    `toml:"type"`         // 客户端类型
	PowerFactor int    `toml:"power_factor"` // 客户端起始功率
}

type ConfigServer struct {
	Socket            string `toml:"socket"` // 服务端套接字
	ReportConsumption bool   `toml:"report_consumption"`
}

type Config struct {
	Client ConfigClient `toml:"client"`
	Server ConfigServer `toml:"server"`
}

func LoadConfig(name string) *Config {
	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 解析Toml配置
	config := Config{}
	if err := toml.NewDecoder(file).Decode(&config); err != nil {
		panic(err)
	}
	return &config
}
