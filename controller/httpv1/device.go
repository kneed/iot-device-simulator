package httpv1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kneed/iot-device-simulator/controller/form"
	"github.com/kneed/iot-device-simulator/pkg/app"
	"github.com/kneed/iot-device-simulator/services"
	log "github.com/sirupsen/logrus"
	"math"
	"strconv"
)

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
		log.Errorf("GetDevices失败,%+v", err)
		g.Response(app.Error, nil)
		return
	}
	totalCount, err := services.Device.CountDevice(queryConditions)
	if err != nil {
		log.Errorf("CountDevice失败,%+v", err)
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


func CreateDevice(c *gin.Context) {
	var (
		g    = app.Gin{Ctx: c}
		deviceForm form.CreateDeviceForm
	)
	if err := c.BindJSON(&deviceForm); err != nil {
		log.Error("CreateDevice参数绑定错误,err:%+v", err)
		g.Response(app.InvalidParams, nil)
		return
	}
	device, err := services.Device.CreateDevice(deviceForm)
	if err != nil {
		log.Errorf("%+v",err)
		g.Response(app.Error, nil)
		return
	}
	g.Response(app.Success, device)
	return

}

func PatchDevice(c *gin.Context) {
	var (
		g             = app.Gin{Ctx: c}
		deviceContent = map[string]interface{}{}
	)
	err := c.ShouldBindJSON(&deviceContent)
	if err != nil {
		log.Errorf("%+v",err)
		g.Response(app.InvalidParams, nil)
		return
	}
	deviceID:= c.Param("device_id")
	fmt.Println(deviceID)
	deviceIDInt,_ := strconv.Atoi(deviceID)
	isExist, err := services.Device.IsExistByID(deviceIDInt)
	if err != nil {
		log.Errorf("%+v",err)
		g.Response(app.Error, nil)
		return
	}
	if !isExist {
		g.Response(app.ObjectNotExist, nil)
		return
	}
	device, err := services.Device.UpdateDevice(deviceIDInt, deviceContent)
	if err != nil {
		log.Errorf("%+v", err)
		g.Response(app.Error, nil)
		return
	}
	g.Response(app.Success, device)
	return
}

func RestartDevice(c *gin.Context) {
	var (
		g             = app.Gin{Ctx: c}
		deviceContent = form.RestartDeviceForm{}
		deviceIDs     []int
	)
	err := c.ShouldBindJSON(&deviceContent)
	if err != nil {
		log.Errorf("%+v",err)
		g.Response(app.InvalidParams, nil)
		return
	}
	deviceIDs = deviceContent.DeviceIDS
	err = services.Device.DeviceRestart(deviceIDs)
	if err != nil {
		log.Errorf("重启模拟器失败. erorr:%+v", err)
		g.Response(app.Error, nil)
		return
	}
	g.Response(app.Success, nil)
	return
}