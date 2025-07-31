package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger(logFile string, logPath string) (*logrus.Logger, error) {
	Log = logrus.New()
	Log.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	// Construct the log file path
	if logPath != "" {
		if err := os.MkdirAll(logPath, 0755); err != nil {
			return nil, err
		}
		logFile = logPath + "/" + logFile
	}

	// Output to both file and console
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	Log.SetOutput(io.MultiWriter(os.Stdout, file))

	return Log, nil
}
