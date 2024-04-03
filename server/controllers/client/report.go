package client

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vistart/project20240227/server/common"
)

func Report(c *gin.Context) {
	client, existed := c.Get("client")
	if !existed {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "client not found")
		return
	}
	m := client.(*common.Client)
	consumption, existed := c.GetPostForm("consumption")
	cF, _ := strconv.ParseFloat(consumption, 32)
	recordedAt, existed := c.GetPostForm("recorded_at")
	recordedAtInt, _ := strconv.ParseInt(recordedAt, 10, 32)
	_, err := m.ReceiveReportConsumption(float32(cF), time.Unix(recordedAtInt, 0))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "success")
	return
}
