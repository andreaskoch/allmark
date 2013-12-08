// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/model"
)

func New(logger logger.Logger) (*Server, error) {
	return &Server{
		logger: logger,
	}, nil
}

type Server struct {
	logger logger.Logger
}

func (server *Server) Serve(item *model.Item) {
	server.logger.Info("Serving item %q", item)
}
