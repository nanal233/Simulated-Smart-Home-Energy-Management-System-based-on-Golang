package client

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vistart/project20240227/server/common"
	"github.com/vistart/project20240227/server/models"
)

type RequestListParams struct {
	Type models.ClientType `form:"type"`
}

func (p *RequestListParams) String() string {
	// 定义输出字符串
	var s string

	// 添加字段值
	s += fmt.Sprintf("type=%d ", p.Type)

	// 返回输出字符串
	return s
}

type ResponseListData struct {
	Clients []models.Client `json:"clients"`
}

func List(c *gin.Context) {
	params := RequestListParams{}
	if err := c.ShouldBindQuery(&params); err == nil {
		log.Println(params.String())
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	p, ok := c.Get("page_size")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, "page and size not specified")
	}
	paramPageSize := p.(*RequestPageParams)

	clients, count, err := models.GetClients(common.DB, paramPageSize.Page, paramPageSize.Size, params.Type)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// 查询每个客户端是否活跃。

	for i := 0; i < len(clients); i++ {
		clients[i].IsActive = common.GlobalSessionManager.GetIsActive(clients[i].ID)
	}

	response := ResponseList{
		Data:  ResponseListData{Clients: clients},
		Count: count,
	}
	c.JSON(http.StatusOK, response)
}
