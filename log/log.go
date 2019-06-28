package log

import (
	"os"

	"github.com/rs/zerolog"
)

var Entry zerolog.Logger

type Config struct {
	// file name
	Dir string
}

func Init(c *Config) {
	Entry = zerolog.New(os.Stdout).With().Logger().Level(zerolog.InfoLevel)
}
