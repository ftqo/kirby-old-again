package logger

import (
	"compress/gzip"
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
var fp string

func Initialize() {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	logsPath := path.Join(d, "../logs")
	_, err := os.Open(logsPath)
	if err != nil {
		err = os.Mkdir(logsPath, 0o755)
		if err != nil {
			panic("Failed to create logs folder: " + err.Error())
		}
	}
	fp = logsPath + "/" + time.Now().Format("2006-01-02--15-04-05") + ".log"
	f, err = os.Create(fp)
	if err != nil {
		panic("Failed to create logs file: " + err.Error())
	}
	L = log.Output(zerolog.ConsoleWriter{Out: io.MultiWriter(os.Stdout, f)})
	L.Info().Msg("Initialized logger")
}

func Close() {
	L.Info().Msg("Closing logger and compressing logs; THIS SHOULD BE THE FINAL LOG MESSAGE")
	cf, err := os.Create(fp + ".gz")
	if err != nil {
		L.Error().Err(err).Msg("Failed to create file for compressed logs; keeping non-compressed logs")
		f.Close()
		return
	}
	gz := gzip.NewWriter(cf)
	f.Close()
	f, _ = os.Open(fp)
	_, err = io.Copy(gz, f)
	gz.Close()
	if err != nil {
		L.Error().Err(err).Msg("Failed to compress logs; keeping non-compressed logs")
		f.Close()
	} else {
		f.Close()
		os.Remove(fp)
	}
}

func NoDebug() {
	zerolog.SetGlobalLevel(1)
}
