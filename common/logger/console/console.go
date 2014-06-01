// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package console

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger/loglevel"
	"io"
	"os"
)

const (
	LogLevelDebug = "Debug"
	LogLevelInfo  = "Info"
	LogLevelWarn  = "Warn"
	LogLevelError = "Error"
	LogLevelFatal = "Fatal"
)

func New(loglevel loglevel.LogLevel) *ConsoleLogger {
	return &ConsoleLogger{
		output: os.Stdout,
		level:  loglevel,
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

	logger.print(LogLevelDebug, format, v)
}

func (logger *ConsoleLogger) Info(format string, v ...interface{}) {
	if logger.level > loglevel.Info {
		return
	}

	logger.print(LogLevelInfo, format, v)
}

func (logger *ConsoleLogger) Warn(format string, v ...interface{}) {
	if logger.level > loglevel.Warn {
		return
	}

	logger.print(LogLevelWarn, format, v)
}

func (logger *ConsoleLogger) Error(format string, v ...interface{}) {
	if logger.level > loglevel.Error {
		return
	}

	logger.print(LogLevelError, format, v)
}

func (logger *ConsoleLogger) Fatal(format string, v ...interface{}) {
	logger.print(LogLevelFatal, format, v)
	os.Exit(1)
}

func (logger *ConsoleLogger) print(level, format string, v []interface{}) {
	if len(v) > 0 {
		fmt.Fprintf(logger.output, level+": \t"+format+"\n", v)
	} else {
		fmt.Fprintf(logger.output, level+": \t"+format+"\n")
	}
}
