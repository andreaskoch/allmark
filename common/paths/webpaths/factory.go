// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webpaths

import (
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
)

func NewFactory(logger logger.Logger, index *index.Index) *PatherFactory {
	return &PatherFactory{
		logger: logger,
		index:  index,
	}
}

type PatherFactory struct {
	logger logger.Logger
	index  *index.Index
}

func (factory *PatherFactory) Absolute(prefix string) paths.Pather {
	return newAbsoluteWebPathProvider(factory.logger, factory.index, prefix)
}

func (factory *PatherFactory) Relative(baseRoute *route.Route) paths.Pather {
	return newRelativeWebPathProvider(factory.logger, factory.index, baseRoute)
}
