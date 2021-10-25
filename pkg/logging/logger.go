package logging

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/kneed/iot-device-simulator/settings"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const DefaultTimeFormat = "2006-01-02 15:04:05"

var logLevelMap = map[string]logrus.Level{
	"Info":   logrus.InfoLevel,
	"Debug":  logrus.DebugLevel,
	"Waring": logrus.WarnLevel,
}

func Init() {
	logFileName := settings.AppSetting.LogFileName
	logMaxDays := settings.AppSetting.LogMaxDays
	setLogger(NewLogFileWriter(logFileName, logMaxDays))
}

func NewLogger(output io.Writer) *logrus.Logger {
	logLevel := logLevelMap[settings.AppSetting.LogLevel]
	logger := logrus.New()
	logger.SetLevel(logLevel)
	logger.SetReportCaller(true) // 调用的文件名,行号等信息
	logger.SetOutput(output)     // 设置writer, 日志将写入到哪里
	logger.SetFormatter(new(logFormater))
	return logger

}

func setLogger(output io.Writer) {
	logLevel := logLevelMap[settings.AppSetting.LogLevel]
	logrus.SetLevel(logLevel)
	logrus.SetReportCaller(true) // 调用的文件名,行号等信息
	logrus.SetOutput(output)     // 设置writer, 日志将写入到哪里
	logrus.SetFormatter(new(logFormater))
}


func NewLogFileWriter(name string, days int) io.Writer {
	pwd, _ := os.Getwd()

	logDir := path.Join(pwd, "log")
	filePath := path.Join(logDir, name)
	rotateLogs, err := rotatelogs.New(
		filePath+".%Y-%m-%d"+".log",
		rotatelogs.WithLinkName(filePath+".log"),
		rotatelogs.WithMaxAge(time.Duration(days)*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		panic(err)
	}
	switch settings.AppSetting.LogLevel {
	case "Debug":
		return io.MultiWriter(rotateLogs, os.Stdout)
	default:
		return rotateLogs
	}

}

type logFormater struct {
}

// Format 自定义日志的格式
func (s *logFormater) Format(entry *logrus.Entry) ([]byte, error) {
	file := filepath.Base(entry.Caller.File)
	text := fmt.Sprintf("[%s][%v][%v %v]:%v",
		strings.ToUpper(entry.Level.String()),
		entry.Time.Format(DefaultTimeFormat),
		file,
		entry.Caller.Line,
		entry.Message,
	)
	if len(entry.Data) > 0 {
		encoder := jsoniter.NewEncoder(entry.Buffer)
		_ = encoder.Encode(entry.Data)
		text += "\t" + entry.Buffer.String()
	} else {
		text += "\n"
	}

	return []byte(text), nil
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		requestBody, _ := ioutil.ReadAll(tee)
		c.Request.Body = ioutil.NopCloser(&buf)
		//处理请求
		c.Next()
		//请求方式
		reqMethod := c.Request.Method
		//请求路由
		reqUrl := c.Request.RequestURI
		//状态码
		statusCode := c.Writer.Status()
		logrus.Debug(fmt.Sprintf("%s:%s\nBody:%s\nresp_code:%d",
			reqMethod, reqUrl, requestBody, statusCode))
	}
}