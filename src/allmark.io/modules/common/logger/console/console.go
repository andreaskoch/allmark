// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package console

import (
	"allmark.io/modules/common/logger/loglevel"
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
	log.SetOutput(w)
}

func (logger *ConsoleLogger) Debug(format string, v ...interface{}) {
	if logger.level > loglevel.Debug {
		return
	}

	log.Println(fmt.Sprintf("%13s", LogLevelDebug) + fmt.Sprintf("%4s", "") + fmt.Sprintf(format, v...))
}

func (logger *ConsoleLogger) Level() loglevel.LogLevel {
	return logger.level
}

func (logger *ConsoleLogger) Info(format string, v ...interface{}) {
	if logger.level > loglevel.Info {
		return
	}

	log.Println(fmt.Sprintf("%13s", LogLevelInfo) + fmt.Sprintf("%4s", "") + fmt.Sprintf(format, v...))
}

func (logger *ConsoleLogger) Statistics(format string, v ...interface{}) {
	if logger.level > loglevel.Statistics {
		return
	}

	log.Println(fmt.Sprintf("%13s", LogLevelStatistics) + fmt.Sprintf("%4s", "") + fmt.Sprintf(format, v...))
}

func (logger *ConsoleLogger) Warn(format string, v ...interface{}) {
	if logger.level > loglevel.Warn {
		return
	}

	log.Println(fmt.Sprintf("%13s", LogLevelWarn) + fmt.Sprintf("%4s", "") + fmt.Sprintf(format, v...))
}

func (logger *ConsoleLogger) Error(format string, v ...interface{}) {
	if logger.level > loglevel.Error {
		return
	}

	log.Println(fmt.Sprintf("%13s", LogLevelError) + fmt.Sprintf("%4s", "") + fmt.Sprintf(format, v...))
}

func (logger *ConsoleLogger) Fatal(format string, v ...interface{}) {
	log.Fatalln(LogLevelFatal + fmt.Sprintf("%4s", "") + fmt.Sprintf(format, v...))
	os.Exit(1)
}
