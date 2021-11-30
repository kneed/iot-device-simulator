package simulator

import (
	"fmt"
	mqtt2 "github.com/eclipse/paho.mqtt.golang"
	"github.com/kneed/iot-device-simulator/db/models"
	"github.com/kneed/iot-device-simulator/pkg/mqtt"
	"github.com/kneed/iot-device-simulator/services"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"sync"
)

var DeviceMqttClientMap = sync.Map{}

// StartSimulator 读取数据库的数据,交给goroutine运行
func StartSimulator() {
	devices, err := services.Device.GetAllDevices(map[string]interface{}{})
	if err != nil {
		log.Error(err)
		return
	}
	for _, device := range devices {
		fmt.Println(device.ID)
		handleDeviceProtocols(device)
	}
}

// 处理设备协议
func handleDeviceProtocols(device *models.Device) {
	protocols, err := models.GetProtocolsByDeviceId(device.ID)
	if err != nil {
		log.Error(err)
		return
	}
	for _, protocol := range protocols {
		newProtocol := protocol
		if protocol.Type == 0 {
			go deviceAutoPublish(device, newProtocol)
		} else if protocol.Type == 1 {
			go deviceSubscribe(device, newProtocol)
		}
	}
}

// 读取协议并订阅相关topic
func deviceSubscribe(device *models.Device, protocol *models.Protocol) {
	key := fmt.Sprintf("device_%d", device.ID)
	var mqttClient mqtt2.Client
	value, ok := DeviceMqttClientMap.Load(key)
	if !ok {
		newMqttClient, err := mqtt.NewMqttClient(device.ServerIp, device.ServerPort, key)
		if err != nil {
			log.Errorf("订阅mqtt失败. Host:%s, IP:%s", device.ServerIp, device.ServerPort)
			return
		}
		DeviceMqttClientMap.Store(key, newMqttClient)
		mqttClient = newMqttClient
	} else {
		mqttClient = value.(mqtt2.Client)
	}
	mqttClient.(mqtt2.Client).Subscribe(
		protocol.SubTopic,
		byte(protocol.Qos),
		mqtt.NewMessageHandlerWithProtocol(mqttClient.(mqtt2.Client), *protocol),
	)

}

// 定时发送协议
func deviceAutoPublish(device *models.Device, protocol *models.Protocol) {
	c := cron.New()
	pubFunc := func(){
		var mqttClient mqtt2.Client
		key := fmt.Sprintf("device_%d", device.ID)
		value, ok := DeviceMqttClientMap.Load(key)
		if !ok {
			newMqttClient, err := mqtt.NewMqttClient(device.ServerIp, device.ServerPort, key)
			if err != nil {
				log.Errorf("创建mqttClient失败. Host:%s, IP:%s", device.ServerIp, device.ServerPort)
				return
			}
			DeviceMqttClientMap.Store(key, newMqttClient)
			mqttClient = newMqttClient
		} else {
			mqttClient = value.(mqtt2.Client)
		}
		mqtt.Publish(mqttClient.(mqtt2.Client), protocol.PubTopic, protocol.Qos, protocol.Content)
	}
	_, err := c.AddFunc("@every 1m", pubFunc)
	if err != nil {
		log.Errorf("定时任务创建失败. error:%s", err)
		return
	}
	c.Start()
}
