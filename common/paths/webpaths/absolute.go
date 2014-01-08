// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webpaths

import (
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
)

// Create a new absolute web path provider
func newAbsoluteWebPathProvider(logger logger.Logger, index *index.Index, prefix string) *AbsoluteWebPathProvider {
	return &AbsoluteWebPathProvider{
		prefix: prefix,
		logger: logger,
		index:  index,
	}
}

type AbsoluteWebPathProvider struct {
	prefix string
	logger logger.Logger
	index  *index.Index
}

// Get the absolute path for the supplied item
func (webPathProvider *AbsoluteWebPathProvider) Path(itemPath string) string {
	return webPathProvider.prefix + "/" + itemPath
}
