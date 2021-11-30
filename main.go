package main

import (
	"github.com/kneed/iot-device-simulator/db/models"
	"github.com/kneed/iot-device-simulator/pkg/logging"
	"github.com/kneed/iot-device-simulator/server"
	"github.com/kneed/iot-device-simulator/settings"
	"github.com/kneed/iot-device-simulator/simulator"
)

func main() {
	settings.Init()
	logging.Init()
	models.InitDB()
	simulator.StartSimulator()
	server.Run()
}
