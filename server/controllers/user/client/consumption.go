package client

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vistart/project20240227/server/common"
	"github.com/vistart/project20240227/server/models"
)

type ResponseGetConsumptionsData struct {
	Consumptions []models.ClientConsumption `json:"consumptions"`
}

func GetConsumptions(c *gin.Context) {
	clientID, ok := c.Get("client-id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, "client id not specified")
	}
	p, ok := c.Get("page_size")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, "page and size not specified")
	}
	paramPageSize := p.(*RequestPageParams)

	client, err := models.GetClient(common.DB, clientID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseList{
			Data: ResponseGetConsumptionsData{Consumptions: nil},
		})
		return
	}

	consumptions, count, err := client.GetConsumptions(common.DB, paramPageSize.Page, paramPageSize.Size)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseList{
			Data: ResponseGetConsumptionsData{Consumptions: nil},
		})
		return
	}

	response := ResponseList{
		Data:  ResponseGetConsumptionsData{Consumptions: consumptions},
		Count: count,
	}
	c.JSON(http.StatusOK, response)
}
