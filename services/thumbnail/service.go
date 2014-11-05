// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package thumbnail

import (
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/dataaccess"
)

func NewConversionService(logger logger.Logger, config config.Config, repository dataaccess.Repository) *ConversionService {
	return &ConversionService{
		logger:     logger,
		config:     config,
		repository: repository,
	}
}

type ConversionService struct {
	logger     logger.Logger
	config     config.Config
	repository dataaccess.Repository
}
