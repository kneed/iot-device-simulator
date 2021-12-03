package services

import (
	"github.com/kneed/iot-device-simulator/controller/form"
	"github.com/kneed/iot-device-simulator/db/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	Device            = DeviceService{}
	DeviceUpdateChan  = make(chan *models.Device, 5)
	DeviceDeleteChan  = make(chan *models.Device, 5)
	DeviceRestartChan = make(chan *models.Device, 5)
)

type DeviceService struct{}

// GetDevices 获取设备分页列表
func (d *DeviceService) GetDevices(pageNum int, pageSize int, conditions interface{}, order string) ([]*models.Device, error) {
	devices, err := models.GetDevices(pageNum, pageSize, conditions, order)
	if err != nil {
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
		Name:       deviceForm.Name,
		Type:       deviceForm.Type,
		ServerIp:   deviceForm.ServerIp,
		ServerPort: deviceForm.ServerPort,
		State:      deviceForm.State,
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

func (d *DeviceService) IsExistByID(deviceID int) (bool, error) {
	var (
		db     = models.GetDB()
		device models.Device
	)
	err := db.Select("id").Where("id = ?", deviceID).First(&device).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, errors.Wrap(err, "数据库查询device时错误")
	}
	if device.ID != 0 {
		return true, nil
	}
	return false, nil
}

func (d *DeviceService) UpdateDevice(deviceID int, data interface{}) (*models.Device, error) {
	device, err := models.UpdateDevice(deviceID, data)
	if err != nil {
		return nil, err
	}
	DeviceUpdateChan <- device
	return device, nil
}

func (d *DeviceService) DeviceRestart(deviceIDs []int) error{
	var (
		db = models.GetDB()
	)
	err := db.Table("device").Where("id IN ?", deviceIDs).Updates(map[string]interface{}{"state": "running"}).Error
	if err != nil {
		return errors.Wrap(err, "数据库更新device时出错")
	}
	for _, deviceID := range deviceIDs{
		device, err := models.GetDevice(deviceID)
		if err != nil {
			return errors.Wrap(err, "数据库查询device时出错")
		}
		DeviceRestartChan <- device
	}
	return nil
}
