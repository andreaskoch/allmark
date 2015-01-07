// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package console

import (
	"allmark.io/modules/common/logger/loglevel"
	"bytes"
	"strings"
	"testing"
)

func Test_Debug(t *testing.T) {
	// arrange
	buf := new(bytes.Buffer)

	message := "A test message"

	logger := New(loglevel.Debug)
	logger.SetOutput(buf)

	// act
	logger.Debug(message)

	// assert
	logOutput := buf.String()

	// test the log level prefix
	if !strings.Contains(logOutput, LogLevelDebug) {
		t.Errorf("The log message should contain the log level prefix %q", LogLevelDebug)
	}

	// test the log message content
	if !strings.Contains(logOutput, message) {
		t.Errorf("The function should have written %q to the log.", message)
	}
}

func Test_Info(t *testing.T) {
	// arrange
	buf := new(bytes.Buffer)

	message := "A test message"

	logger := New(loglevel.Debug)
	logger.SetOutput(buf)

	// act
	logger.Info(message)

	// assert
	logOutput := buf.String()

	// test the log level prefix
	if !strings.Contains(logOutput, LogLevelInfo) {
		t.Errorf("The log message should contain the log level prefix %q", LogLevelInfo)
	}

	// test the log message content
	if !strings.Contains(logOutput, message) {
		t.Errorf("The function should have written %q to the log.", message)
	}
}

func Test_Warn(t *testing.T) {
	// arrange
	buf := new(bytes.Buffer)

	message := "A test message"

	logger := New(loglevel.Debug)
	logger.SetOutput(buf)

	// act
	logger.Warn(message)

	// assert
	logOutput := buf.String()

	// test the log level prefix
	if !strings.Contains(logOutput, LogLevelWarn) {
		t.Errorf("The log message should contain the log level prefix %q", LogLevelWarn)
	}

	// test the log message content
	if !strings.Contains(logOutput, message) {
		t.Errorf("The function should have written %q to the log.", message)
	}
}

func Test_Error(t *testing.T) {
	// arrange
	buf := new(bytes.Buffer)

	message := "A test message"

	logger := New(loglevel.Debug)
	logger.SetOutput(buf)

	// act
	logger.Error(message)

	// assert
	logOutput := buf.String()

	// test the log level prefix
	if !strings.Contains(logOutput, LogLevelError) {
		t.Errorf("The log message should contain the log level prefix %q", LogLevelError)
	}

	// test the log message content
	if !strings.Contains(logOutput, message) {
		t.Errorf("The function should have written %q to the log.", message)
	}
}
