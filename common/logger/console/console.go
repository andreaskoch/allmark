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

func (logger *ConsoleLogger) Debug(format string, v ...interface{}) {
	fmt.Fprintf(logger.output, fmt.Sprintf("%s: \t%s\n", LogLevelDebug, format), v)
}

func (logger *ConsoleLogger) Info(format string, v ...interface{}) {
	fmt.Fprintf(logger.output, fmt.Sprintf("%s: \t%s\n", LogLevelInfo, format), v)
}

func (logger *ConsoleLogger) Warn(format string, v ...interface{}) {
	fmt.Fprintf(logger.output, fmt.Sprintf("%s: \t%s\n", LogLevelWarn, format), v)
}

func (logger *ConsoleLogger) Error(format string, v ...interface{}) {
	fmt.Fprintf(logger.output, fmt.Sprintf("%s: \t%s\n", LogLevelError, format), v)
}

func (logger *ConsoleLogger) Fatal(format string, v ...interface{}) {
	fmt.Fprintf(logger.output, fmt.Sprintf("%s: \t%s\n", LogLevelFatal, format), v)
}
