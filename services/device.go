package services

import (
	"github.com/kneed/iot-device-simulator/controller/form"
	"github.com/kneed/iot-device-simulator/db/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)


var Device = DeviceService{}

type DeviceService struct{}

// GetDevices 获取设备分页列表
func (d *DeviceService) GetDevices(pageNum int, pageSize int, conditions interface{}, order string) ([]*models.Device, error) {
	devices, err := models.GetDevices(pageNum, pageSize, conditions, order)
	if err != nil{
		return nil, errors.Wrapf(err, "DeviceService GetDevices异常")
	}
	return devices, nil
}

// GetAllDevices 获取所有设备
func (d *DeviceService) GetAllDevices(conditions interface{}) ([]*models.Device, error) {
	db := models.GetDB()
	var devices []*models.Device
	err := db.Where(conditions).Find(&devices).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(err, "GetAllDevices error:")
	}
	return devices, nil
}

// CreateDevice 创建设备
func (d *DeviceService) CreateDevice(deviceForm form.CreateDeviceForm) (*models.Device, error) {
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

// CountDevice 统计设备数量
func (d *DeviceService) CountDevice(conditions interface{}) (int64, error) {
	count, err := models.CountDevice(conditions)
	if err != nil {
		return 0, errors.Wrapf(err, "DeviceServie CountDevice异常")
	}
	return count, nil
}