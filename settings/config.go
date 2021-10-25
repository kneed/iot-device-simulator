package settings

import (
	"fmt"
	"github.com/spf13/viper"
	"net/url"
	"strconv"
)

type Database struct {
	Dsn               string
	Url               string
	MaxOpenConnection int
	MaxIdleConnection int
}

type App struct {
	LogLevel    string
	LogFileName string
	LogMaxDays  int
}

type Server struct {
	Port string
}

var (
	Config          *viper.Viper
	DatabaseSetting Database
	AppSetting      App
	ServerSetting   Server
)

func Init() {
	Config = viper.New()
	// DB
	_ = Config.BindEnv("DB_HOST")
	_ = Config.BindEnv("DB_USER")
	_ = Config.BindEnv("DB_PASSWORD")
	_ = Config.BindEnv("DB_NAME")
	_ = Config.BindEnv("DB_PORT")
	_ = Config.BindEnv("MAX_OPEN_CONNECTION")
	_ = Config.BindEnv("MAX_IDLE_CONNECTION")

	// LOG
	_ = Config.BindEnv("LOG_LEVEL")
	_ = Config.BindEnv("LOG_FILENAME")
	_ = Config.BindEnv("LOG_MAXDAYS")

	// Server
	_ = Config.BindEnv("SERVER_PORT")

	dbHost := Config.GetString("DB_HOST")
	dbUser := Config.GetString("DB_USER")
	dbPort := Config.GetString("DB_PORT")
	dbPassword := Config.GetString("DB_PASSWORD")
	dbPasswordUrlFormat := url.QueryEscape(dbPassword)
	dbName := Config.GetString("DB_NAME")
	var (
		dbDsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			dbHost, dbUser, dbPassword, dbName, dbPort)
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPasswordUrlFormat, dbHost, dbPort, dbName)
	)
	DatabaseSetting = Database{
		Dsn:               dbDsn,
		Url:               dbURL,
		MaxOpenConnection: Config.GetInt("MAX_OPEN_CONNECTION"),
		MaxIdleConnection: Config.GetInt("MAX_IDLE_CONNECTION"),
	}
	logMaxDays, _ := strconv.Atoi(Config.GetString("LOG_MAXDAYS"))
	AppSetting = App{
		LogLevel:    Config.GetString("LOG_LEVEL"),
		LogFileName: Config.GetString("LOG_FILENAME"),
		LogMaxDays:  logMaxDays,
	}

	ServerSetting = Server{
		Port: Config.GetString("SERVER_PORT"),
	}
}
