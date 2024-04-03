package common

import (
	"errors"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type ConfigDatabase struct {
	DSN string `toml:"dsn"`
}

type ConfigSessionManager struct {
	BroadcastTimestampInterval int64 `toml:"broadcast_timestamp_interval"`
}

type Config struct {
	Port               uint16               `toml:"port"`
	Database           ConfigDatabase       `toml:"database"`
	BroadcastTimestamp ConfigSessionManager `toml:"session_manager"`
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
	if config.BroadcastTimestamp.BroadcastTimestampInterval <= 0 {
		panic(errors.New("broadcast timestamp interval is zero"))
	}
	return &config
}
