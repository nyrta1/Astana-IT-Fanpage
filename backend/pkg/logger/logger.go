package logger

import "github.com/sirupsen/logrus"

var log = logrus.New()

func InitLogger() {
	log.SetFormatter(&logrus.TextFormatter{})
	log.SetLevel(logrus.InfoLevel)
}

func LogWithFields(fields logrus.Fields) *logrus.Entry {
	return log.WithFields(fields)
}

func GetLogger() *logrus.Logger {
	return log
}
