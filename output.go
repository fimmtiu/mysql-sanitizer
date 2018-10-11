package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	logLevelNormal  = 0
	logLevelVerbose = 1
	logLevelDebug   = 2
	logLevelDump    = 3
)

// Output is just a wrapper around log.Logger.
type Output struct {
	Logger *log.Logger
	Level  int
}

// NewOutput returns a new Output object.
func NewOutput(config Config) (out Output) {
	var fileHandle io.Writer

	if config.LogFile != "-" {
		file, err := os.OpenFile(config.LogFile, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Fatalf("Can't open logfile %s: %s", config.LogFile, err)
		}
		fileHandle = file
	} else {
		fileHandle = os.Stdout
	}

	out.Logger = log.New(fileHandle, "", log.Ldate|log.Lmicroseconds|log.LUTC)
	out.Level = config.LogLevel
	return out
}

// Dump dumps a hexadecimal version of the given chunk of memory to the log.
func (out Output) Dump(slice []byte, format string, args ...interface{}) {
	if out.Level >= logLevelDump {
		str := fmt.Sprintf(format, args...)
		rowCount := len(slice) / 16
		if len(slice)%16 > 0 {
			rowCount++
		}
		for i := 0; i < rowCount; i++ {
			str += "      "
			for j := 0; j < 16; j++ {
				if len(slice)-i*16 <= j {
					str += "   "
				} else {
					str += fmt.Sprintf("%02x ", slice[i*16+j])
				}
			}

			str += "      "
			for j := 0; j < 16; j++ {
				if len(slice)-i*16 <= j {
					str += "  "
				} else if slice[i*16+j] >= 0x20 && slice[i*16+j] <= 0x7e {
					str += fmt.Sprintf("%c ", slice[i*16+j])
				} else {
					str += "* "
				}
			}

			str += "\n"
		}
		out.Logger.Print(str)
	}
}

// Debug prints internal debugging messages to the log.
func (out Output) Debug(format string, args ...interface{}) {
	if out.Level >= logLevelDebug {
		out.Logger.Printf(format, args...)
	}
}

// Verbose prints low-priority messages to the log.
func (out Output) Verbose(format string, args ...interface{}) {
	if out.Level >= logLevelVerbose {
		out.Logger.Printf(format, args...)
	}
}

// Log prints ordinary messages to the log.
func (out Output) Log(format string, args ...interface{}) {
	out.Logger.Printf(format, args...)
}
