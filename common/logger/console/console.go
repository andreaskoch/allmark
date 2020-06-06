// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package console provides a (console) logger that implements the
// allmark.io/modules/common/logger.Logger interface.
package console

import (
	"github.com/elWyatt/allmark/common/logger/loglevel"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	LogLevelDebug      = "Debug"
	LogLevelInfo       = "Info"
	LogLevelStatistics = "Statistics"
	LogLevelWarn       = "Warn"
	LogLevelError      = "Error"
	LogLevelFatal      = "Fatal"
)

// Default creates a default ConsoleLogger with the Info log level and os.Stdout as the output target.
func Default() *ConsoleLogger {
	return New(loglevel.Info)
}

// New creates a new instance of the ConsoleLogger with os.Stdout as the output target.
func New(level loglevel.LogLevel) *ConsoleLogger {
	return &ConsoleLogger{
		output: os.Stdout,
		level:  level,
	}
}

// ConsoleLogger implements the allmark.io/modules/common/logger.Logger interface
// and provides the ability to write log messages to a given output writer.
type ConsoleLogger struct {
	output io.Writer
	level  loglevel.LogLevel
}

// SetOutput sets the output of this logger to the supplied io.Writer.
func (logger *ConsoleLogger) SetOutput(w io.Writer) {
	log.SetOutput(w)
}

// Level returns the current log level.
func (logger *ConsoleLogger) Level() loglevel.LogLevel {
	return logger.level
}

// Debug formats according to a format specifier and writes a debug log message to standard output.
func (logger *ConsoleLogger) Debug(format string, v ...interface{}) {
	if logger.level > loglevel.Debug {
		return
	}

	log.Println(fmt.Sprintf("%13s", LogLevelDebug) + fmt.Sprintf("%4s", "") + fmt.Sprintf(format, v...))
}

// Info formats according to a format specifier and writes an info log message to standard output.
func (logger *ConsoleLogger) Info(format string, v ...interface{}) {
	if logger.level > loglevel.Info {
		return
	}

	log.Println(fmt.Sprintf("%13s", LogLevelInfo) + fmt.Sprintf("%4s", "") + fmt.Sprintf(format, v...))
}

// Statistics formats according to a format specifier and writes a statistics log message to standard output.
func (logger *ConsoleLogger) Statistics(format string, v ...interface{}) {
	if logger.level > loglevel.Statistics {
		return
	}

	log.Println(fmt.Sprintf("%13s", LogLevelStatistics) + fmt.Sprintf("%4s", "") + fmt.Sprintf(format, v...))
}

// Warn formats according to a format specifier and writes a warn log message to standard output.
func (logger *ConsoleLogger) Warn(format string, v ...interface{}) {
	if logger.level > loglevel.Warn {
		return
	}

	log.Println(fmt.Sprintf("%13s", LogLevelWarn) + fmt.Sprintf("%4s", "") + fmt.Sprintf(format, v...))
}

// Error formats according to a format specifier and writes an error log message to standard output.
func (logger *ConsoleLogger) Error(format string, v ...interface{}) {
	if logger.level > loglevel.Error {
		return
	}

	log.Println(fmt.Sprintf("%13s", LogLevelError) + fmt.Sprintf("%4s", "") + fmt.Sprintf(format, v...))
}

// Fatal formats according to a format specifier and writes a fatal log message to standard output and exits the application.
func (logger *ConsoleLogger) Fatal(format string, v ...interface{}) {
	log.Fatalln(LogLevelFatal + fmt.Sprintf("%4s", "") + fmt.Sprintf(format, v...))
	os.Exit(1)
}
