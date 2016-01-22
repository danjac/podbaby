package commands

import (
	"github.com/Sirupsen/logrus"
)

func configureLogger() *logrus.Logger {
	logger := logrus.New()

	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	}

	return logger

}
