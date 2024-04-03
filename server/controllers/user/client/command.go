package client

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/vistart/project20240227/server/common"
	"github.com/vistart/project20240227/server/models"
)

type RequestSendCommandParams struct {
	Command string `form:"command"`
	Data    string `form:"data"`
}

func (p *RequestSendCommandParams) String() string {
	// 定义输出字符串
	var s string

	// 添加字段值
	s += fmt.Sprintf("command=%s ", p.Command)
	s += fmt.Sprintf("data=%s", p.Data)

	// 返回输出字符串
	return s
}

func (p *RequestSendCommandParams) Check() error {
	if len(p.Command) == 0 {
		return errors.New("empty command")
	}
	_, err := strconv.ParseInt(p.Data, 10, 64)
	if err != nil {
		return err
	}
	return nil
}

func sendCommandPower(client *common.Client, params *RequestSendCommandParams) error {
	data, _ := strconv.Atoi(params.Data)

	command := common.NewEventCommandPower(data)
	now := time.Now()
	client.SendToSessionChannel(command)
	// 发送命令后，将刚才发送的内容记录到数据表中。
	clientModel, _ := models.GetClient(common.DB, client.ID())
	if clientModel != nil {
		_, err := clientModel.InsertNewCommandExecution(common.DB, command.Code, command.MarshalData(), &now)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}
	return nil
}

func SendCommand(c *gin.Context) {
	params := RequestSendCommandParams{}
	if err := c.MustBindWith(&params, binding.Form); err == nil {
		log.Println(params.String())
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := params.Check(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	clientParam, ok := c.Get("client")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "invalid client")
		return
	}

	client := clientParam.(*common.Client)
	var err error

	command := "command-" + params.Command

	switch command {
	case common.EventNameCommandPower:
		err = sendCommandPower(client, &params)
	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, "command not supported")
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "Success")
}

func GetCommands(c *gin.Context) {

}
