package logger

import (
	"context"
	"fmt"
	"os"
	"time"
	"vc-core/internal/config"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

type Logger struct{
	*zerolog.Logger
}

const LoggerKey int = 0

func New() Logger{
	logFolder := "logs"
	env := os.Getenv(config.EnvKey)
	createLogs := env == "development"
	if createLogs {
		if _, err := os.Stat(logFolder); os.IsNotExist(err) {
			err := os.Mkdir(logFolder, os.ModePerm)
			if err != nil {
				fmt.Printf("Error creating logs folder: %v\n", err)
			}
		}
	}
	logFile := ""
	if createLogs {
		logFile = fmt.Sprintf("logs/app-%s.log", time.Now().Format("2006-01-02"))
	}
	logRotator := &lumberjack.Logger{
		Filename:   logFile,
		MaxBackups: 3, 
		MaxAge:     7,
		Compress:   true,
	}
	log := zerolog.New(logRotator).With().Timestamp().Logger()
	return Logger{&log}
}

func FromContext(ctx context.Context) Logger{
	if l, ok := ctx.Value(LoggerKey).(Logger); ok {
		return l
	}
	return New()
}
