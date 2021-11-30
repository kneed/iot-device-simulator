package models

import (
	"github.com/pkg/errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Protocol struct {
	Model

	DeviceId int           `json:"device_id"`
	Name     string         `json:"name"`
	Content  datatypes.JSON `json:"content"`
	Qos      int            `json:"qos"`
	Type     int            `json:"type"` // 0--自发, 1--响应
	SubTopic string         `json:"sub_topic"`
	PubTopic string         `json:"pub_topic"`
	Strategy datatypes.JSON `json:"strategy"`
}

func CreateProtocol(protocol *Protocol) (*Protocol, error) {
	if err := db.Create(protocol).Error; err != nil {
		return nil, errors.Wrap(err, "CreateProtocol异常")
	}
	return protocol, nil
}

func UpdateProtocol(id int, data interface{}) (*Protocol, error) {
	if err := db.Model(&Protocol{}).Where("id = ? AND deleted_at = ?", id, nil).Updates(data).Error; err != nil {
		return nil, errors.Wrap(err, "UpdateProtocol异常")
	}
	protocol, err := GetProtocol(id)
	if err != nil {
		return nil, err
	}
	return protocol, nil
}

func GetProtocols(pageNum int, pageSize int, maps interface{}, order string) ([]*Protocol, error) {
	var protocols []*Protocol
	err := db.Where(maps).Order(order).Scopes(Paginate(pageNum, pageSize)).Find(&protocols).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(err, "GetProtocols异常")
	}
	return protocols, nil
}

func GetProtocol(id int) (*Protocol, error) {
	var protocol Protocol
	err := db.Where("id = ? AND deleted_at = ?", id, nil).First(&protocol).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(err, "GetProtocol异常")
	}
	return &protocol, nil
}

func DeleteProtocol(id int) error {
	if err := db.Where("id = ?", id).Delete(&Protocol{}).Error; err != nil {
		return errors.Wrap(err, "DeleteProtocol异常")
	}
	return nil
}

func GetProtocolsByDeviceId(deviceId int) ([]*Protocol, error) {
	var protocols []*Protocol
	err := db.Where("device_id = ?", deviceId, nil).Find(&protocols).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(err, "GetProtocolsByDeviceId异常")
	}
	return protocols, nil
}

func CountProtocol(conditions interface{}) (int64, error) {
	var count int64
	if err := db.Model(&Protocol{}).Where(conditions).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
