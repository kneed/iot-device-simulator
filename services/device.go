package services

import (
	"github.com/kneed/iot-device-simulator/controller/httpv1"
	"github.com/kneed/iot-device-simulator/db/models"
	"github.com/pkg/errors"
)


var Device = DeviceService{}

type DeviceService struct{}

// 获取设备列表
func (d *DeviceService) GetDevices(pageNum int, pageSize int, conditions interface{}, order string) ([]*models.Device, error) {
	devices, err := models.GetDevices(pageNum, pageSize, conditions, order)
	if err != nil{
		return nil, errors.Wrapf(err, "DeviceService GetDevices异常")
	}
	return devices, nil
}

// 创建设备
func (d *DeviceService) CreateDevice(deviceForm httpv1.CreateDeviceForm) (*models.Device, error) {
	var device = &models.Device{
		Name: deviceForm.Name,
		Type: deviceForm.Type,
		ServerIp: deviceForm.ServerIp,
		ServerPort: deviceForm.ServerPort,
	}
	device, err := models.CreateDevice(device)
	if err != nil {
		return nil, errors.Wrapf(err, "DeviceService CreateDevice异常")
	}
	return device, nil
}

// 统计设备数量
func (d *DeviceService) CountDevice(conditions interface{}) (int64, error) {
	count, err := models.CountDevice(conditions)
	if err != nil {
		return 0, errors.Wrapf(err, "DeviceServie CountDevice异常")
	}
	return count, nil
}