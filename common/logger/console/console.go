// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package console

import (
	"fmt"
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

func New() *ConsoleLogger {
	return &ConsoleLogger{
		output: os.Stdout,
	}
}

type ConsoleLogger struct {
	output io.Writer
}

func (logger *ConsoleLogger) SetOutput(w io.Writer) {
	logger.output = w
}

func (logger *ConsoleLogger) Debug(v ...interface{}) {
	fmt.Fprintf(logger.output, LogLevelDebug, v, "\n")
}

func (logger *ConsoleLogger) Info(v ...interface{}) {
	fmt.Fprintf(logger.output, LogLevelInfo, v, "\n")
}

func (logger *ConsoleLogger) Warn(v ...interface{}) {
	fmt.Fprintf(logger.output, LogLevelWarn, v, "\n")
}

func (logger *ConsoleLogger) Error(v ...interface{}) {
	fmt.Fprintf(logger.output, LogLevelError, v, "\n")
}

func (logger *ConsoleLogger) Fatal(v ...interface{}) {
	fmt.Fprintf(logger.output, LogLevelFatal, v, "\n")
}
