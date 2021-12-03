package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kneed/iot-device-simulator/db/models"
	"github.com/kneed/iot-device-simulator/services"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

const (
	qos        byte = 0
	retryTimes int  = 3
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Info("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	reader := client.OptionsReader()
	log.Infof("Mqtt Connected, %s", reader.ClientID())
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	reader := client.OptionsReader()
	log.Warnf("connect lost: %+v, server:%s", err, reader.Servers())
	clientID := reader.ClientID()
	deviceID, err := strconv.Atoi(strings.Split(clientID, "_")[2])
	if err != nil {
		log.Errorf("获取device_id失败. error:%+v", err)
	}
	_, err = services.Device.UpdateDevice(deviceID, map[string]interface{}{"state": "error"})
	if err != nil {
		log.Errorf("处理模拟器与mqtt broker连接丢失时发生错误. error:%+v", err)
	}
}

// NewMqttClient 创建新的连接
func NewMqttClient(broker string, port string, clientId string) (mqtt.Client, error) {
	var client mqtt.Client
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", broker, port))
	opts.SetClientID(fmt.Sprintf("simulator_%s", clientId))
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.SetAutoReconnect(false)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return client, nil
}

func Publish(client mqtt.Client, topic string, qos int, data []byte) {
	reader := client.OptionsReader()
	log.Infof("MQTT publish topic:%s, data:%+v", topic, string(data))
	var success bool
	// 失败重试
	for i := 0; i < retryTimes; i++ {
		token := client.Publish(topic, byte(qos), false, data)
		if err := token.Error(); err != nil {

			log.Errorf("MQTT publish error:%+v, current_retry_times:%d total_retry_times:%d, cliend_id:%s, server:%s",
				err, i, retryTimes, reader.ClientID(), reader.Servers())
			time.Sleep(time.Second * 3)
			continue
		}
		success = token.Wait()
	}
	if !success && !client.IsConnected(){
		log.Errorf("设备与mqtt服务器之间连接异常,停止模拟器运行")
		clientID := reader.ClientID()
		deviceID, err := strconv.Atoi(strings.Split(clientID, "_")[2])
		if err != nil {
			log.Errorf("获取device_id失败. error:%+v", err)
		}
		_, err = services.Device.UpdateDevice(deviceID, map[string]interface{}{"state": "error"})
		if err != nil {
			log.Errorf("处理模拟器与mqtt broker连接丢失时发生错误. error:%+v", err)
		}
	}
}

// NewMessageHandlerWithProtocol 构建一个消息处理函数
func NewMessageHandlerWithProtocol(mqttClient mqtt.Client, protocol models.Protocol) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		log.Infof("Received message: %s from topic: %s", msg.Payload(), msg.Topic())
		Publish(mqttClient, protocol.PubTopic, protocol.Qos, protocol.Content)
	}
}
