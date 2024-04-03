package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	// 定义命令行参数
	inputPtr := flag.String("config", "./conf/type1device1.toml", "输入文件")

	// 解析命令行参数
	flag.Parse()

	if inputPtr == nil && len(*inputPtr) == 0 {
		fmt.Println("config file not specified")
		return
	}

	// 访问命令行参数的值
	config := LoadConfig(*inputPtr)
	apiSocket = config.Server.Socket
	Client := NewClient(config.Client.ID, config.Client.Type, config.Client.PowerFactor)
	// 启动注册逻辑并持续接收服务端发来的命令。
	exitChannel = make(chan bool)
	go Client.Register(exit)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		// 每秒向服务端报告一次瞬时功率。
		case <-ticker.C:
			if config.Server.ReportConsumption {
				Client.Report()
			}
		case <-exitChannel:
			ticker.Stop()
			return
		}
	}
}

var apiSocket = "localhost:59002"

var exitChannel chan bool

func exit() {
	exitChannel <- true
}
