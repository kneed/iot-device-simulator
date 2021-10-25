package httpv1

import (
	"github.com/gin-gonic/gin"
	"github.com/kneed/iot-device-simulator/pkg/app"
	"github.com/kneed/iot-device-simulator/services"
	log "github.com/sirupsen/logrus"
	"math"
	"strconv"
)

// @Summary 获取协议列表
// @tags 协议接口
// @Produce  json
// @Param type query string false "Type"
// @Param device_id query string false "Type"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/protocols [get]
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
		log.Errorf("GetProtocols失败,%s", err)
		g.Response(app.Error, nil)
		return
	}
	totalCount, err := services.Protocol.CountProtocol(queryConditions)
	if err != nil {
		log.Errorf("CountProtocol失败,%s", err)
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

// @Summary 创建协议
// @tags 协议接口
// @Produce  json
// @Param name body string true "Name"
// @Param type body int false "Type"
// @Param server_ip body string true "ServerIp"
// @Param server_port body string true "ServerPort"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/protocols [post]
func CreateProtocol(c *gin.Context) {
	var (
		g    = app.Gin{Ctx: c}
		form CreateProtocolForm
	)
	if err := c.BindJSON(&form); err != nil {
		log.Error("CreateProtocol参数绑定错误,err:%s", err)
		g.Response(app.InvalidParams, nil)
		return
	}
	protocol, err := services.Protocol.CreateProtocol(form)
	if err != nil {
		log.Error(err)
		g.Response(app.Error, nil)
		return
	}
	g.Response(app.Success, protocol)
	return

}
