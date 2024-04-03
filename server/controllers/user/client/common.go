package client

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vistart/project20240227/server/common"
)

type ResponseList struct {
	Data  any   `json:"data"`
	Count int64 `json:"count"`
}

// Authorize 表示客户端通用验证逻辑。
// 客户端需要在请求的 header 中提交 client-id, client-type, authorization 字段值。
// 如果验证通过，则继续后续逻辑。如果验证不通过，则直接中止。
func Authorize(c *gin.Context) {
	clientIDPost := c.PostForm("client_id")
	clientIDQuery := c.Query("client_id")
	clientID := ""
	if len(clientIDPost) > 0 {
		clientID = clientIDPost
	} else if len(clientIDQuery) > 0 {
		clientID = clientIDQuery
	} else {
		c.String(http.StatusBadRequest, common.ErrRequestBadClientID{}.Error())
		c.Abort()
		return
	}
	c.Set("client-id", clientID)
	c.Next()
}

type RequestPageParams struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

func (p *RequestPageParams) String() string {
	// 定义输出字符串
	var s string

	// 添加字段值
	s += fmt.Sprintf("page=%d ", p.Page)
	s += fmt.Sprintf("size=%d ", p.Size)

	// 返回输出字符串
	return s
}

func NewRequestPageParams(page *int, size *int) *RequestPageParams {
	p := 1
	if page != nil && *page > 1 {
		p = *page
	}
	s := 10
	if size != nil && *size > 1 {
		s = *size
	}
	return &RequestPageParams{
		Page: p, Size: s,
	}
}

func BindPageSize(c *gin.Context) {
	params := NewRequestPageParams(nil, nil)
	if err := c.ShouldBindQuery(params); err == nil {
		log.Println(params.String())
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	if params.Size <= 0 || params.Size > 100 { // 设上限和下限。
		params.Size = 100
	}
	c.Set("page_size", params)
	c.Next()
}
