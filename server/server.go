package server

import (
	"github.com/kneed/iot-device-simulator/settings"
	log "github.com/sirupsen/logrus"
)

func Init() {
	serverPort := settings.ServerSetting.Port
	r := NewRouter()
	err := r.Run(":" + serverPort)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}
}

