package models

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)


const(
	DeviceRunning  = "running"
	DeviceIdle  = "idle"
	DeviceError  = "error"
)

type Device struct {
	Model

	Name       string `json:"name"`
	Type       string `json:"type"`
	ServerIp   string `json:"server_ip"`
	ServerPort string `json:"server_port"`
	State      string `json:"state"` // [running, idle, error]
}

func CreateDevice(device *Device) (*Device, error) {
	if err := db.Create(device).Error; err != nil {
		return nil, errors.Wrap(err, "CreateDevice异常")
	}
	return device, nil
}

func UpdateDevice(id int, data interface{}) (*Device, error) {
	if err := db.Model(&Device{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return nil, errors.Wrap(err, "UpdateDevice异常")
	}
	device, err := GetDevice(id)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func GetDevices(pageNum int, pageSize int, conditions interface{}, order string) ([]*Device, error) {
	var devices []*Device
	err := db.Where(conditions).Order(order).Scopes(Paginate(pageNum, pageSize)).Find(&devices).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(err, "GetDevices异常")
	}
	return devices, nil
}

func GetDevice(id int) (*Device, error) {
	var device Device
	err := db.Where("id = ?", id).First(&device).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(err, "GetDevice异常")
	}
	return &device, nil
}

func DeleteDevice(id int) error {
	if err := db.Where("id = ?", id).Delete(&Device{}).Error; err != nil {
		return errors.Wrap(err, "DeleteDevice异常")
	}
	return nil
}

func CountDevice(conditions interface{}) (int64, error) {
	var count int64
	if err := db.Model(&Device{}).Where(conditions).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}


