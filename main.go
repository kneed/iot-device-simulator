package main

import (
	"github.com/kneed/iot-device-simulator/db/models"
	"github.com/kneed/iot-device-simulator/pkg/logging"
	"github.com/kneed/iot-device-simulator/settings"
)

func main() {
	settings.Init()
	logging.Init()
	models.InitDB()
}
