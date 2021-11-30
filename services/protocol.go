package services

import (
	"github.com/kneed/iot-device-simulator/controller/form"
	"github.com/kneed/iot-device-simulator/db/models"
	"github.com/pkg/errors"
)

var Protocol = ProtocolService{}

type ProtocolService struct{}

// 获取协议列表
func (d *ProtocolService) GetProtocols(pageNum int, pageSize int, conditions interface{}, order string) ([]*models.Protocol, error) {
	protocols, err := models.GetProtocols(pageNum, pageSize, conditions, order)
	if err != nil{
		return nil, errors.Wrapf(err, "ProtocolService GetProtocols异常")
	}
	return protocols, nil
}

// 创建协议
func (d *ProtocolService) CreateProtocol(protocolForm form.CreateProtocolForm) (*models.Protocol, error) {
	var protocol = &models.Protocol{
		DeviceId: protocolForm.DeviceId,
		Name: protocolForm.Name,
		Content: protocolForm.Content,
		Qos: *protocolForm.Qos,
		Type: *protocolForm.Type,
		SubTopic: protocolForm.SubTopic,
		PubTopic: protocolForm.PubTopic,
		Strategy: protocolForm.Strategy,
	}
	protocol, err := models.CreateProtocol(protocol)
	if err != nil {
		return nil, errors.Wrapf(err, "ProtocolService CreateProtocol异常")
	}
	return protocol, nil
}

func (d *ProtocolService) CountProtocol(conditions interface{}) (int64, error) {
	count, err := models.CountProtocol(conditions)
	if err != nil {
		return 0, errors.Wrapf(err, "ProtocolServie CountProtocol异常")
	}
	return count, nil
}