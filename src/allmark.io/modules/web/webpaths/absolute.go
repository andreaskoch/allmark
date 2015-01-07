// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webpaths

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/route"
)

// Create a new absolute web path provider
func newAbsoluteWebPathProvider(logger logger.Logger, prefix string) *AbsoluteWebPathProvider {
	return &AbsoluteWebPathProvider{
		prefix: prefix,
		logger: logger,
	}
}

type AbsoluteWebPathProvider struct {
	prefix string
	logger logger.Logger
}

// Get the absolute path for the supplied item
func (webPathProvider *AbsoluteWebPathProvider) Path(itemPath string) string {
	return webPathProvider.prefix + itemPath
}

func (webPathProvider *AbsoluteWebPathProvider) Base() route.Route {
	return route.New()
}
