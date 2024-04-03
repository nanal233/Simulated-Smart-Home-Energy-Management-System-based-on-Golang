package client

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vistart/project20240227/server/common"
	"github.com/vistart/project20240227/server/models"
)

// Authorize 表示客户端通用验证逻辑。
// 客户端需要在请求的 header 中提交 client-id, client-type, authorization 字段值。
// 如果验证通过，则继续后续逻辑。如果验证不通过，则直接中止。
func Authorize(c *gin.Context) {
	clientID := c.GetHeader(common.RequestClientID)
	clientType, err := strconv.Atoi(c.GetHeader(common.RequestClientType))
	if err != nil {
		c.String(http.StatusBadRequest, common.ErrRequestBadClientType{}.Error())
		c.Abort()
		return
	}
	authorization := c.GetHeader(common.RequestAuthorization)

	auth := common.RequestClientAuthorization{
		ClientID:      clientID,
		ClientType:    models.ClientType(clientType),
		Authorization: authorization,
	}
	if err = auth.Auth(); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		c.Abort()
		return
	}
	c.Set("client-id", clientID)
	c.Set("client-type", models.ClientType(clientType))
	c.Next()
}
