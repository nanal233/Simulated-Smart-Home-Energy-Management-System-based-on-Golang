package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/vistart/project20240227/server/common"
	controllerClient "github.com/vistart/project20240227/server/controllers/client"
	controllerUserClient "github.com/vistart/project20240227/server/controllers/user/client"
)

func main() {
	// 定义命令行参数
	inputPtr := flag.String("config", "./conf/server1.toml", "输入文件")

	// 解析命令行参数
	flag.Parse()

	if inputPtr == nil && len(*inputPtr) == 0 {
		fmt.Println("config file not specified")
		return
	}

	// 访问命令行参数的值
	config := common.LoadConfig(*inputPtr)
	common.PrepareDatabase(config.Database.DSN)
	router := gin.Default()

	common.GlobalSessionManager = common.NewSessionManager()
	go common.GlobalSessionManager.Serve()
	go common.GlobalSessionManager.BroadcastTimestamp(config.BroadcastTimestamp.BroadcastTimestampInterval)

	bindRouter(router)
	router.Run(fmt.Sprintf(":%d", config.Port))
}

func bindRouter(e *gin.Engine) {
	client := e.Group("/client")
	// 客户端向服务端注册，服务端向客户端发送命令。Server-sent event模式。
	client.POST("/register", controllerClient.Authorize, common.GlobalSessionManager.SetHeadersHandler(), common.GlobalSessionManager.NewSessionChannelHandler(), controllerClient.Register)
	// 客户端向服务端报告状态。
	client.POST("/report", controllerClient.Authorize, common.GlobalSessionManager.GetClientHandler(), controllerClient.Report)

	// 用户客户端相关
	userClient := e.Group("/user/client")
	// 向客户端发送命令
	userClient.POST("/command", controllerUserClient.Authorize, common.GlobalSessionManager.GetClientHandler(), controllerUserClient.SendCommand)
	// 客户端列表。
	userClient.GET("/list", controllerUserClient.BindPageSize, controllerUserClient.List)
	// 获取某个客户端信息。
	userClient.GET("/info", controllerUserClient.Authorize, controllerUserClient.GetInfo)
	// 编辑某个客户端信息。
	userClient.POST("/info", controllerUserClient.Authorize, common.GlobalSessionManager.GetClientHandler(), controllerUserClient.EditInfo)
	// 删除某个客户端信息。
	userClient.DELETE("/info", controllerUserClient.Authorize, common.GlobalSessionManager.GetClientHandler(), controllerUserClient.DeleteInfo)
	// 获取某个客户端的能耗列表。
	userClient.GET("/consumption", controllerUserClient.Authorize, controllerUserClient.BindPageSize, controllerUserClient.GetConsumptions)
	// 获取某个客户端的能耗模式历史。
	userClient.GET("/command", controllerUserClient.Authorize, controllerUserClient.BindPageSize, controllerUserClient.GetPowerModeHistories)

	// 能耗模式
	userPowerMode := e.Group("/user/power_mode")
	// 能耗模式列表。
	userPowerMode.GET("/list")
	// 获取某个能耗模式。
	userPowerMode.GET("")
	// 添加/编辑某个能耗模式。
	userPowerMode.POST("")
	// 删除某个能耗模式。
	userPowerMode.DELETE("")

	// 能耗模式命令相关
	userPowerModeCommand := userPowerMode.Group("/command")

	// 获取能耗模式执行历史。
	userPowerModeCommand.GET("")

	// 为指定能耗模式添加命令
	userPowerModeCommand.POST("")

	// 删除指定能耗模式的具体命令。
	userPowerModeCommand.DELETE("")

	// 执行指定能耗模式。
	userPowerModeCommand.POST("/execute")

	// 查询指定能耗模式执行历史。
	userPowerModeCommand.GET("/executions")
}
