// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webpaths

import (
	"github.com/elWyatt/allmark/common/logger"
	"github.com/elWyatt/allmark/common/paths"
	"github.com/elWyatt/allmark/common/route"
	"github.com/elWyatt/allmark/dataaccess"
)

func NewFactory(logger logger.Logger, repository dataaccess.Repository) *PatherFactory {
	return &PatherFactory{
		logger:     logger,
		repository: repository,
	}
}

type PatherFactory struct {
	logger     logger.Logger
	repository dataaccess.Repository
}

func (factory *PatherFactory) Absolute(prefix string) paths.Pather {
	return newAbsoluteWebPathProvider(prefix)
}

func (factory *PatherFactory) Relative(baseRoute route.Route) paths.Pather {
	return newRelativeWebPathProvider(factory.repository, baseRoute)
}
