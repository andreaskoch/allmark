// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watcher

import (
	"fmt"
	"strings"
)

const (
	UNKNOWN EventType = iota
	MODIFIED
	CREATED
	DELETED
	RENAMED
)

type EventType int

func EventTypeFromText(eventText string) EventType {
	switch strings.ToLower(strings.TrimSpace(eventText)) {
	case "modified":
		return MODIFIED
	case "created":
		return CREATED
	case "deleted":
		return DELETED
	case "renamed":
		return RENAMED
	case "unknown":
		return UNKNOWN
	default:
		return UNKNOWN
	}

	panic("Unreachable")
}

func (watchEventType EventType) String() string {
	return fmt.Sprintf("%s", getEventName(watchEventType))
}

func getEventName(watchEventType EventType) string {
	switch watchEventType {
	case MODIFIED:
		return "modified"
	case CREATED:
		return "created"
	case DELETED:
		return "deleted"
	case RENAMED:
		return "renamed"
	case UNKNOWN:
		return "unknown"
	default:
		return "unknown"
	}

	panic("Unreachable")
}
