package models

import (
	"github.com/kneed/iot-device-simulator/db/migrate"
	"github.com/kneed/iot-device-simulator/pkg/logging"
	"github.com/kneed/iot-device-simulator/settings"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var db *gorm.DB

func InitDB() {

	var err error
	dsn := settings.DatabaseSetting.Dsn
	dbLogger := logger.New(
		logging.NewLogger(logging.NewLogFileWriter("sql", 7)),
		logger.Config{
			SlowThreshold: 0,
			Colorful:      true,
			LogLevel:      logger.Info,
		},
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction:                   true,
		NamingStrategy:                           schema.NamingStrategy{SingularTable: true},
		DisableAutomaticPing:                     true,
		Logger:                                   dbLogger,
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	if err != nil {
		log.Fatalf("连接到数据库失败. dsn:%s", dsn)
	}
	sqlDb, _ := db.DB()
	sqlDb.SetMaxOpenConns(settings.DatabaseSetting.MaxOpenConnection)
	sqlDb.SetMaxIdleConns(settings.DatabaseSetting.MaxIdleConnection)
	migrate.DbMigrate()
}

type Model struct {
	ID        int            `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Paginate 通用分页逻辑
func Paginate(pageNum int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNum == 0 {
			pageNum = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (pageNum - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func GetDB() *gorm.DB {
	return db
}
