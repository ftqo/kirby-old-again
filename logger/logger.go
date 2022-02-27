package logger

import (
	"io"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var L zerolog.Logger
var f *os.File

func Initialize() {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	logsPath := path.Join(d, "../logs")
	_, err := os.Open(logsPath)
	if err != nil {
		err = os.Mkdir(logsPath, 0o755)
		if err != nil {
			panic("Unable to create logs folder: " + err.Error())
		}
	}
	f, err = os.Create(logsPath + "/" + time.Now().Format("2006-01-02-15:04:05") + ".log")
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	L = log.Output(zerolog.ConsoleWriter{Out: io.MultiWriter(os.Stdout, f)})
	L.Info().Msg("Initialized logger")
}

func Close() {
	L.Info().Msg("Closing logger. THIS SHOULD BE THE FINAL LOG MESSAGE")
	err := f.Close()
	if err != nil {
		// UHHHHHHHHHHHHHHH someoneplstellmewhattodohere
		L.Error().Err(err).Msg("Failed to close logger")
	}
}
