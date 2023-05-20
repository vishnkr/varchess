package logger

import (
	"io"
	"os"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var once sync.Once
var log zerolog.Logger

func Get() zerolog.Logger{
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano

		logLevel,err:= strconv.Atoi(os.Getenv("LOG_LEVEL"))
		if err!=nil{
			logLevel = int(zerolog.InfoLevel)
		}
		var output io.Writer = zerolog.ConsoleWriter{
			Out:os.Stdout,
			TimeFormat: time.RFC3339,
			FieldsExclude: []string{
				"user_agent",
				"git_revision",
				"go_version",
			  },
		}

		/*if os.Getenv("APP_ENV")!="development"{
			fileLogger:=&lumberjack.Logger{
				Filename: "test.log",
				MaxSize: 5,
				MaxBackups: 10,
				MaxAge: 14,
				Compress: true,
			}
			output = zerolog.MultiLevelWriter(os.Stderr, fileLogger)
		}*/
		var gitRevision string
		buildInfo,ok:= debug.ReadBuildInfo()
		if ok{
			for _,v:= range buildInfo.Settings{
				if v.Key == "vcs.revision"{
					gitRevision = v.Value
					break
				}
			}
		}

		log = zerolog.New(output).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Str("git_revision",gitRevision).
			Str("go_version",buildInfo.GoVersion).
			Logger()
	})
	return log
}

