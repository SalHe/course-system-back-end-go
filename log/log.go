package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	C "github.com/se2022-qiaqia/course-system/config"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

var L *zerolog.Logger
var logLevel = map[string]zerolog.Level{
	"debug": zerolog.DebugLevel,
	"info":  zerolog.InfoLevel,
	"warn":  zerolog.WarnLevel,
	"error": zerolog.ErrorLevel,
	"fatal": zerolog.FatalLevel,
	"panic": zerolog.PanicLevel,
}

func Init() {
	level := logLevel[C.Config.Log.Level]

	writers := make([]io.Writer, 1)
	writers[0] = createLumberjackLogger()
	if C.Config.Log.Console {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}
	logger := log.
		Output(io.MultiWriter(writers...)).
		Level(level)
	L = &logger
}

func createLumberjackLogger() *lumberjack.Logger {
	return C.Config.Log.Logger
}
