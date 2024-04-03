package client

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vistart/project20240227/server/common"
	"github.com/vistart/project20240227/server/models"
)

func GetInfo(c *gin.Context) {
	clientID, _ := c.Get("client-id")
	client, err := models.GetClient(common.DB, clientID.(string))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "client not found")
		return
	}
	client.IsActive = common.GlobalSessionManager.GetIsActive(client.ID)
	c.JSON(http.StatusOK, client)
}

func EditInfo(c *gin.Context) {
	clientID, _ := c.Get("client-id")
	client, err := models.GetClient(common.DB, clientID.(string))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "client not found")
		return
	}
	name := c.PostForm("name")
	if len(name) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, "name not specified")
		return
	}
	if len(name) > 255 {
		c.AbortWithStatusJSON(http.StatusBadRequest, "name too long")
		return
	}
	updateName, err := client.UpdateName(common.DB, name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	if updateName == 0 {
		c.JSON(http.StatusOK, "client name not changed")
	} else {
		c.JSON(http.StatusOK, "success")
	}
}

func DeleteInfo(c *gin.Context) {
	clientID, _ := c.Get("client-id")
	client, err := models.GetClient(common.DB, clientID.(string))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "client not found")
		return
	}

	total, err := client.Delete(common.DB)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	// 删除后，检查GlobalSessionManager是否当前正在连接。如果是，则主动断开。
	if common.GlobalSessionManager.GetIsActive(client.ID) {
		clientSession, ok := c.Get("client")
		if ok {
			channel := clientSession.(*common.Client).GetSessionChannel()
			channel <- common.NewEventDisconnect()
		}
	}
	if total == 0 {
		c.JSON(http.StatusOK, "client not deleted")
	} else {
		c.JSON(http.StatusOK, "success")
	}
}
