// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webpaths

import (
	"github.com/andreaskoch/allmark/common/route"
	"strings"
)

// Create a new absolute web path provider
func newAbsoluteWebPathProvider(prefix string) *AbsoluteWebPathProvider {
	return &AbsoluteWebPathProvider{
		prefix: prefix,
	}
}

type AbsoluteWebPathProvider struct {
	prefix string
}

// Get the absolute path for the supplied item
func (webPathProvider *AbsoluteWebPathProvider) Path(itemPath string) string {

	// return the supplied item path if it is already absolute
	if IsAbsoluteURI(itemPath) {
		return itemPath
	}

	// don't do it twice
	if strings.HasPrefix(itemPath, webPathProvider.prefix) {
		return itemPath
	}

	return webPathProvider.prefix + itemPath
}

func (webPathProvider *AbsoluteWebPathProvider) Base() route.Route {
	return route.New()
}
