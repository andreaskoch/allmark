// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package console

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger/loglevel"
	"io"
	"os"
	"time"
)

const (
	LogLevelDebug      = "Debug"
	LogLevelInfo       = "Info"
	LogLevelStatistics = "Statistics"
	LogLevelWarn       = "Warn"
	LogLevelError      = "Error"
	LogLevelFatal      = "Fatal"
)

func Default() *ConsoleLogger {
	return New(loglevel.Info)
}

func New(level loglevel.LogLevel) *ConsoleLogger {
	return &ConsoleLogger{
		output: os.Stdout,
		level:  level,
	}
}

type ConsoleLogger struct {
	output io.Writer
	level  loglevel.LogLevel
}

func (logger *ConsoleLogger) SetOutput(w io.Writer) {
	logger.output = w
}

func (logger *ConsoleLogger) Debug(format string, v ...interface{}) {
	if logger.level > loglevel.Debug {
		return
	}

	logger.print(LogLevelDebug, fmt.Sprintf(format, v...))
}

func (logger *ConsoleLogger) Level() loglevel.LogLevel {
	return logger.level
}

func (logger *ConsoleLogger) Info(format string, v ...interface{}) {
	if logger.level > loglevel.Info {
		return
	}

	logger.print(LogLevelInfo, fmt.Sprintf(format, v...))
}

func (logger *ConsoleLogger) Statistics(format string, v ...interface{}) {
	if logger.level > loglevel.Statistics {
		return
	}

	logger.print(LogLevelStatistics, fmt.Sprintf(format, v...))
}

func (logger *ConsoleLogger) Warn(format string, v ...interface{}) {
	if logger.level > loglevel.Warn {
		return
	}

	logger.print(LogLevelWarn, fmt.Sprintf(format, v...))
}

func (logger *ConsoleLogger) Error(format string, v ...interface{}) {
	if logger.level > loglevel.Error {
		return
	}

	logger.print(LogLevelError, fmt.Sprintf(format, v...))
}

func (logger *ConsoleLogger) Fatal(format string, v ...interface{}) {
	logger.print(LogLevelFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (logger *ConsoleLogger) print(level, message string) {

	timestamp := time.Now().Format(time.RFC1123)
	fmt.Fprintln(logger.output, timestamp+" "+level+":  "+message)
}
