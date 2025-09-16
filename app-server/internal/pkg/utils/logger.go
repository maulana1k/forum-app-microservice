package utils

import (
	"os"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()

	Logger.SetFormatter(&nested.Formatter{
		TimestampFormat: "Jan 02 15:04:05.000",
		HideKeys:        true,
		FieldsOrder:     []string{"component", "method", "path", "request_id"},
	})

	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.InfoLevel)
}

func GetLogger() *logrus.Logger {
	return Logger
}
