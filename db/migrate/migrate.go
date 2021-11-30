package migrate

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kneed/iot-device-simulator/settings"
	log "github.com/sirupsen/logrus"
)

func DbMigrate() {
	dbUrl := settings.DatabaseSetting.Url
	log.Debugf("---------database migrate start---------")
	m, err := migrate.New(
		"file:./migrations",
		dbUrl)

	if err != nil {
		log.Fatalf("new migrate error: %v", err.Error())
	}

	if err = m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Fatalf("%v", err.Error())
		}
		log.Debug(err.Error())
	}
	log.Debug("---------database migrate finished---------")
}
