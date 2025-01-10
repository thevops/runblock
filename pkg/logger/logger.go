package logger

import (
	"os"

	"github.com/charmbracelet/log"
)

var Log *log.Logger

func init() {
	Log = log.New(os.Stdout)

	logLevel := os.Getenv("LOG_LEVEL")
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		Log.SetLevel(log.InfoLevel) // default to InfoLevel if parsing fails
	} else {
		Log.SetLevel(level)
	}
}
