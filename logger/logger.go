package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var L zerolog.Logger

func init() {
	L = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}
