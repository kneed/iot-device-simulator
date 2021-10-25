package httpv1

import (
	"github.com/gin-gonic/gin"
	"github.com/kneed/iot-device-simulator/pkg/app"
	"github.com/kneed/iot-device-simulator/services"
	log "github.com/sirupsen/logrus"
	"math"
	"strconv"
)

// @Summary 获取设备列表
// @tags 设备接口
// @Produce  json
// @Param type query string false "Type"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/devices [get]
func GetDevices(c *gin.Context) {
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
	if deviceType, isExist := c.GetQuery("type"); isExist {
		queryConditions["type"] = deviceType
	}
	if sorter, isExist := c.GetQuery("sorter"); isExist {
		order = sorter
	}
	devices, err := services.Device.GetDevices(currentPage, pageSize, queryConditions, order)
	if err != nil {
		log.Errorf("GetDevices失败,%s", err)
		g.Response(app.Error, nil)
		return
	}
	totalCount, err := services.Device.CountDevice(queryConditions)
	if err != nil {
		log.Errorf("CountDevice失败,%s", err)
		g.Response(app.Error, nil)
		return
	}
	totalPage := int(math.Ceil(float64(totalCount) / float64(currentPage)))
	respData := app.ResponsePageBody{
		Total:       totalCount,
		TotalPage:   totalPage,
		PageSize:    pageSize,
		CurrentPage: currentPage,
		List:        devices,
	}
	g.Response(app.Success, respData)
	return
}

// @Summary 创建设备
// @tags 设备接口
// @Produce  json
// @Param name body string true "Name"
// @Param type body int false "Type"
// @Param server_ip body string true "ServerIp"
// @Param server_port body string true "ServerPort"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/devices [post]
func CreateDevice(c *gin.Context) {
	var (
		g    = app.Gin{Ctx: c}
		form CreateDeviceForm
	)
	if err := c.BindJSON(&form); err != nil {
		log.Error("CreateDevice参数绑定错误,err:%s", err)
		g.Response(app.InvalidParams, nil)
		return
	}
	device, err := services.Device.CreateDevice(form)
	if err != nil {
		log.Error(err)
		g.Response(app.Error, nil)
		return
	}
	g.Response(app.Success, device)
	return

}
