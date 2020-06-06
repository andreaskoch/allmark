// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

import (
	"github.com/elWyatt/allmark/common/logger/loglevel"
)

type Logger interface {
	// Level returns the current log level.
	Level() loglevel.LogLevel

	// Debug formats according to a format specifier and writes a debug log message.
	Debug(format string, v ...interface{})

	// Info formats according to a format specifier and writes an info log message.
	Info(format string, v ...interface{})

	// Statistics formats according to a format specifier and writes a statistics log message.
	Statistics(format string, v ...interface{})

	// Warn formats according to a format specifier and writes a warns log message.
	Warn(format string, v ...interface{})

	// Errror formats according to a format specifier and writes an error log message.
	Error(format string, v ...interface{})

	// Fatal formats according to a format specifier and writes a fatal log message.
	Fatal(format string, v ...interface{})
}
