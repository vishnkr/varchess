package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
	"vc-server/vc-core/internal/config"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

type Logger struct{
	*zerolog.Logger
}

const LoggerKey int = 0

func New() Logger {
	env := os.Getenv(config.EnvKey)
	isDev := env == "development"

	var writer io.Writer
	if isDev {
		writer = zerolog.MultiLevelWriter(os.Stdout, &lumberjack.Logger{
			Filename:   fmt.Sprintf("logs/app-%s.log", time.Now().Format("2006-01-02")),
			MaxBackups: 3,
			MaxAge:     7,
			Compress:   true,
		})
	} else {
		writer = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	}

	logger := zerolog.New(writer).With().Timestamp().Logger()
	return Logger{&logger}
}


func FromContext(ctx context.Context) Logger{
	if l, ok := ctx.Value(LoggerKey).(Logger); ok {
		return l
	}
	return New()
}
