package app

import (
	"github.com/sirupsen/logrus"
	"os"
)

func SetLogrus(level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(lvl)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05",
	})

	logrus.SetOutput(os.Stdout)
}
