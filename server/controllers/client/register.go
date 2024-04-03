package client

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/vistart/project20240227/server/common"
)

func Register(c *gin.Context) {
	v, ok := c.Get("client")
	if !ok {
		return
	}
	client, ok := v.(*common.Client)
	if !ok {
		return
	}
	c.Stream(func(w io.Writer) bool {
		// Stream message to client from message channel
		channel := client.GetSessionChannel()
		if event, ok := <-channel; ok {
			if event == nil {
				return true
			}
			if eventD, ok := event.(*common.EventBase[any]); ok {
				c.SSEvent(common.EventCodeNameMap[eventD.Code], eventD.MarshalData())
			} else if eventD, ok := event.(*common.EventBase[common.EventCommandPowerData]); ok {
				c.SSEvent(common.EventCodeNameMap[eventD.Code], eventD.MarshalData())
			} else if eventD, ok := event.(*common.EventBase[common.EventMessageData]); ok {
				c.SSEvent(common.EventCodeNameMap[eventD.Code], eventD.MarshalData())
			} else if eventD, ok := event.(*common.EventBase[struct{}]); ok {
				c.SSEvent(common.EventCodeNameMap[eventD.Code], eventD.MarshalData()) // 删除后，中断连接。
				return false
			}
			return true
		}
		return false
	})
}
