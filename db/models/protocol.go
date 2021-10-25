package models

import (
	"github.com/pkg/errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Protocol struct {
	Model

	ProtocolId uint           `json:"device_id"`
	Name       string         `json:"name"`
	Content    datatypes.JSON `json:"content"`
	Qos        int            `json:"qos"`
	Type       int            `json:"type"` // 0--自发, 1--响应
	SubTopic   string         `json:"sub_topic"`
	PubTopic   string         `json:"pub_topic"`
	StrategyId int            `json:"strategy_id"`
}

func CreateProtocol(device *Protocol) (*Protocol, error) {
	if err := db.Create(device).Error; err != nil {
		return nil, errors.Wrap(err, "CreateProtocol异常")
	}
	return device, nil
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
	var devices []*Protocol
	err := db.Where(maps).Order(order).Scopes(Paginate(pageNum, pageSize)).Find(&devices).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(err, "GetProtocols异常")
	}
	return devices, nil
}

func GetProtocol(id int) (*Protocol, error) {
	var device Protocol
	err := db.Where("id = ? AND deleted_at = ?", id, 0).First(&device).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(err, "GetProtocol异常")
	}
	return &device, nil
}

func DeleteProtocol(id int) error {
	if err := db.Where("id = ?", id).Delete(&Protocol{}).Error; err != nil {
		return errors.Wrap(err, "DeleteProtocol异常")
	}
	return nil
}

func CountProtocol(conditions interface{}) (int64, error) {
	var count int64
	if err := db.Model(&Protocol{}).Where(conditions).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}