// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/converter"
)

func New(logger logger.Logger, converter converter.Converter) (*Server, error) {
	return &Server{
		logger: logger,
	}, nil
}

type Server struct {
	logger logger.Logger
}

func (server *Server) Serve(item *model.Item) {
	server.logger.Debug("Serving item %q", item)
}
