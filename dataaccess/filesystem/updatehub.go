// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
)

func newUpdateHub(logger logger.Logger) *UpdateHub {
	return &UpdateHub{logger}
}

type UpdateHub struct {
	logger logger.Logger
}

func (hub *UpdateHub) StartWatching(route route.Route) {
	hub.logger.Debug("Starting watcher for %q", route)
}

func (hub *UpdateHub) StopWatching(route route.Route) {
	hub.logger.Debug("Stopping watcher for %q", route)

}
