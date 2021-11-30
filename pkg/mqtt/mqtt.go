package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kneed/iot-device-simulator/db/models"
	log "github.com/sirupsen/logrus"
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
	log.Info("Mqtt Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Warnf("Connect lost: %v", err)
}

// NewMqttClient todo 相同的server应该被复用
func NewMqttClient(broker string, port string, clientId string) (mqtt.Client, error) {
	var client mqtt.Client
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", broker, port))
	opts.SetClientID(fmt.Sprintf("simulator_%s", clientId))
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return client, nil
}

func Publish(client mqtt.Client, topic string, qos int, data []byte) {
	log.Infof("MQTT publish topic:%s, data:%+v", topic, string(data))
	// 失败重试
	for i := 0; i < retryTimes; i++ {
		token := client.Publish(topic, byte(qos), false, data)
		if err := token.Error(); err != nil {
			log.Error("MQTT publish e:", err, " current_retry_times:", i, " total_retry_times:", retryTimes)
			time.Sleep(time.Second * 3)
			continue
		}
		token.Wait()
		break
	}
}

// NewMessageHandlerWithProtocol 构建一个消息处理函数
func NewMessageHandlerWithProtocol(mqttClient mqtt.Client, protocol models.Protocol) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		log.Infof("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
		Publish(mqttClient, protocol.PubTopic, protocol.Qos, protocol.Content)
	}
}
