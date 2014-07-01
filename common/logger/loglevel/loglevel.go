// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package loglevel

import (
	"strings"
)

type LogLevel byte

const (
	Off = iota
	Debug
	Info
	Warn
	Error
	Fatal
)

func (logLevel LogLevel) String() string {
	switch logLevel {

	case Debug:
		return "Debug"

	case Info:
		return "Info"

	case Warn:
		return "Warn"

	case Error:
		return "Error"

	case Fatal:
		return "Fatal"

	case Off:
		return "Off"

	default:
		return "Off"

	}

	panic("Unreachable")
}

func FromString(logLevel string) LogLevel {

	logLevel = strings.ToLower(strings.TrimSpace(logLevel))
	switch logLevel {

	case "debug":
		return Debug

	case "info":
		return Info

	case "warn":
		return Warn

	case "error":
		return Error

	case "fatal":
		return Fatal

	case "off":
		return Off

	default:
		return Info

	}

	panic("Unreachable")
}
