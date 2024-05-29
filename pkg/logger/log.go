package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/zer0day88/tinder/config"
	"github.com/zer0day88/tinder/pkg/environment"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func New() zerolog.Logger {

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	var output io.Writer

	logMinimumLevel := config.Key.LogLevel

	if config.Key.Environment == environment.Production {
		logMinimumLevel = zerolog.WarnLevel
		output = os.Stdout
	} else {
		output = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		}
	}

	log := zerolog.New(output).With().
		Timestamp().Caller().Logger()

	log.Level(logMinimumLevel)

	return log

}
