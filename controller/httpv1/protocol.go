package httpv1

import (
	"github.com/gin-gonic/gin"
	form2 "github.com/kneed/iot-device-simulator/controller/form"
	"github.com/kneed/iot-device-simulator/pkg/app"
	"github.com/kneed/iot-device-simulator/services"
	log "github.com/sirupsen/logrus"
	"math"
	"strconv"
)

func GetProtocols(c *gin.Context) {
	var (
		g               = app.Gin{Ctx: c}
		queryConditions = make(app.Map)
		order           string
	)
	currentPage, err := strconv.Atoi(c.DefaultQuery("currentPage", "1"))
	if err != nil {
		g.Response(app.InvalidParams, nil)
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil {
		g.Response(app.InvalidParams, nil)
		return
	}
	if protocolType, isExist := c.GetQuery("type"); isExist {
		queryConditions["type"] = protocolType
	}
	if deviceId, isExist := c.GetQuery("device_id"); isExist {
		queryConditions["device_id"] = deviceId
	}
	if sorter, isExist := c.GetQuery("sorter"); isExist {
		order = sorter
	}
	protocols, err := services.Protocol.GetProtocols(currentPage, pageSize, queryConditions, order)
	if err != nil {
		log.Errorf("GetProtocols失败,%+v", err)
		g.Response(app.Error, nil)
		return
	}
	totalCount, err := services.Protocol.CountProtocol(queryConditions)
	if err != nil {
		log.Errorf("CountProtocol失败,%+v", err)
		g.Response(app.Error, nil)
		return
	}
	totalPage := int(math.Ceil(float64(totalCount) / float64(currentPage)))
	respData := app.ResponsePageBody{
		Total:       totalCount,
		TotalPage:   totalPage,
		PageSize:    pageSize,
		CurrentPage: currentPage,
		List:        protocols,
	}
	g.Response(app.Success, respData)
	return
}


func CreateProtocol(c *gin.Context) {
	var (
		g    = app.Gin{Ctx: c}
		form form2.CreateProtocolForm
	)
	if err := c.BindJSON(&form); err != nil {
		log.Errorf("CreateProtocol参数绑定错误,err:%+v", err)
		g.Response(app.InvalidParams, nil)
		return
	}
	protocol, err := services.Protocol.CreateProtocol(form)
	if err != nil {
		log.Errorf("%+v", err)
		g.Response(app.Error, nil)
		return
	}
	g.Response(app.Success, protocol)
	return
}
