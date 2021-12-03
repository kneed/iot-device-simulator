package simulator

import (
	"context"
	"encoding/json"
	"fmt"
	mqtt2 "github.com/eclipse/paho.mqtt.golang"
	"github.com/kneed/iot-device-simulator/db/models"
	"github.com/kneed/iot-device-simulator/pkg/mqtt"
	"github.com/kneed/iot-device-simulator/services"
	"github.com/kneed/iot-device-simulator/utils"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var DeviceMqttClientMap = sync.Map{}
var DeviceCtxMap = sync.Map{}

// StartSimulator 读取数据库的数据,交给goroutine运行
func StartSimulator() {
	queryCondition := map[string]interface{}{
		"state": "running",
	}
	devices, err := services.Device.GetAllDevices(queryCondition)
	if err != nil {
		log.Error(err)
		return
	}
	for _, device := range devices {
		runDevice(device)
	}
	go handleDeviceUpdate(services.DeviceUpdateChan, services.DeviceDeleteChan, services.DeviceRestartChan)
}

// 处理设备变更
func handleDeviceUpdate(deviceUpdateChan, deviceDeleteChan, deviceRestartChan chan *models.Device) {

	defer utils.WrapRecover()

	var device *models.Device
	for {
		select {
		case device = <- deviceUpdateChan:
			cancel, ok := DeviceCtxMap.Load(device.ID)
			if device.State == models.DeviceIdle || device.State == models.DeviceError  {
				if ok {
					key := fmt.Sprintf("device_%d", device.ID)
					cancel.(context.CancelFunc)()
					DeviceCtxMap.Delete(device.ID)
					client, ok := DeviceMqttClientMap.Load(key)
					if ok {
						client.(mqtt2.Client).Disconnect(1000)
						DeviceMqttClientMap.Delete(key)
					}
				}
			}else if device.State == models.DeviceRunning {
				if !ok {
					runDevice(device)
				}
			}
		case device = <- deviceDeleteChan:
			cancel, ok := DeviceCtxMap.Load(device.ID)
			if ok {
				cancel.(context.CancelFunc)()
				DeviceCtxMap.Delete(device.ID)
				key := fmt.Sprintf("device_%d", device.ID)
				client, ok := DeviceMqttClientMap.Load(key)
				if ok {
					client.(mqtt2.Client).Disconnect(1000)
					DeviceMqttClientMap.Delete(key)
				}
			}
		case device = <- deviceRestartChan:
			cancel, ok := DeviceCtxMap.Load(device.ID)
			if ok {
				cancel.(context.CancelFunc)()
				DeviceCtxMap.Delete(device.ID)
				key := fmt.Sprintf("device_%d", device.ID)
				client, ok := DeviceMqttClientMap.Load(key)
				if ok {
					client.(mqtt2.Client).Disconnect(1000)
					DeviceMqttClientMap.Delete(key)
				}
			}
			runDevice(device)
		}
	}
}


// 处理设备协议
func runDevice(device *models.Device) {
	protocols, err := models.GetProtocolsByDeviceId(device.ID)
	if err != nil {
		log.Error(err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	DeviceCtxMap.Store(device.ID, cancel)
	for _, protocol := range protocols {
		newProtocol := protocol
		if protocol.Type == 0 {
			go deviceAutoPublish(ctx, device, newProtocol)
		} else if protocol.Type == 1 {
			go deviceSubscribe(ctx, device, newProtocol)
		}
	}
}

// 读取协议并订阅相关topic
func deviceSubscribe(ctx context.Context, device *models.Device, protocol *models.Protocol) {
	defer utils.WrapRecover()
	key := fmt.Sprintf("device_%d", device.ID)
	var mqttClient mqtt2.Client
	value, ok := DeviceMqttClientMap.Load(key)
	if !ok {
		newMqttClient, err := mqtt.NewMqttClient(device.ServerIp, device.ServerPort, key)
		if err != nil {
			log.Errorf("创建mqttClient失败. Host:%s, IP:%s, error:%+v", device.ServerIp, device.ServerPort, err)
			data := map[string]interface{}{
				"state": "error",
			}
			_, err = services.Device.UpdateDevice(device.ID, data)
			if err != nil {
				log.Errorf("update device error:%+v", err)
				return
			}
			return
		}
		DeviceMqttClientMap.Store(key, newMqttClient)
		mqttClient = newMqttClient
	} else {
		mqttClient = value.(mqtt2.Client)
	}
	mqttClient.Subscribe(
		protocol.SubTopic,
		byte(protocol.Qos),
		mqtt.NewMessageHandlerWithProtocol(mqttClient.(mqtt2.Client), *protocol),
	)
	log.Infof("device_id:%d, protocol_id:%d [订阅]运行中",device.ID, protocol.ID)
	for {
		select {
		case <-ctx.Done():
			log.Infof("终止设备订阅任务, device_id:%d", device.ID)
			//if err:=mqttClient.Unsubscribe(protocol.SubTopic).Error(); err != nil {
			//	log.Errorf("unsubscibe 异常. error:%+v", err)
			//}
			return
		}
	}

}

// 定时发送协议
func deviceAutoPublish(ctx context.Context, device *models.Device, protocol *models.Protocol) {
	defer utils.WrapRecover()
	c := cron.New()
	pubFunc := func() {
		var mqttClient mqtt2.Client
		key := fmt.Sprintf("device_%d", device.ID)
		value, ok := DeviceMqttClientMap.Load(key)
		if !ok {
			newMqttClient, err := mqtt.NewMqttClient(device.ServerIp, device.ServerPort, key)
			if err != nil {
				log.Errorf("创建mqttClient失败. Host:%s, IP:%s", device.ServerIp, device.ServerPort)
				data := map[string]string{
					"state": "error",
				}
				_, err = services.Device.UpdateDevice(device.ID, data)
				if err != nil {
					log.Errorf("update device error:%+v", err)
					return
				}
				return
			}
			DeviceMqttClientMap.Store(key, newMqttClient)
			mqttClient = newMqttClient
		} else {
			mqttClient = value.(mqtt2.Client)
		}
		mqtt.Publish(mqttClient.(mqtt2.Client), protocol.PubTopic, protocol.Qos, protocol.Content)
	}

	var strategy map[string]interface{}
	err := json.Unmarshal(protocol.Strategy, &strategy)
	if err != nil {
		log.Errorf("handle strategy error. protocol_id:%d, error:%+v", protocol.ID, err)
	}

	intervalSpec := fmt.Sprintf("@every %ds", int(strategy["interval"].(float64)))
	entryID, err := c.AddFunc(intervalSpec, pubFunc)
	if err != nil {
		log.Errorf("定时任务创建失败. error:%s", err)
		return
	}
	c.Start()
	log.Infof("device_id:%d, protocol_id:%d, [发布]运行中", device.ID, protocol.ID)
	if strategy["duration"].(float64) == 0 {
		for {
			select {
			case <-ctx.Done():
				log.Infof("终止设备发布任务, device_id:%d", device.ID)
				c.Remove(entryID)
				return
			}
		}
	}else {
		for {
			select {
			case <-ctx.Done():
				log.Infof("终止设备发布任务, device_id:%d", device.ID)
				c.Remove(entryID)
				return
			case <-time.After(time.Duration(strategy["duration"].(float64)) * time.Second):
				log.Infof("定时发布任务完成, device_id:%d, strategy:%+v", device.ID, strategy)
				c.Remove(entryID)
				return
			}
		}
	}
}
