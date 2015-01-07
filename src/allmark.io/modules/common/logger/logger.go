// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

import (
	"allmark.io/modules/common/logger/loglevel"
)

type Logger interface {
	Level() loglevel.LogLevel

	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Statistics(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
}
