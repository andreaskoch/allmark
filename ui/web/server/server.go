// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/conversion"
	"math"
	"net/http"
)

func New(logger logger.Logger, config *config.Config, converter conversion.Converter) (*Server, error) {
	return &Server{
		config:    config,
		logger:    logger,
		converter: converter,
		index:     newIndex(logger),
	}, nil
}

type Server struct {
	isRunning bool

	config    *config.Config
	logger    logger.Logger
	converter conversion.Converter
	index     *Index
}

func (server *Server) Serve(item *model.Item) {

	// start the server if it is not running
	if !server.IsRunning() {
		server.start()
	}

	server.logger.Debug("Serving item %q", item)
	server.index.Add(item)
}

func (server *Server) IsRunning() bool {
	return server.isRunning
}

func (server *Server) start() {

	go func() {
		server.isRunning = true

		// start http server: http
		httpBinding := server.getHttpBinding()
		server.logger.Info("Starting http server %q\n", httpBinding)

		if err := http.ListenAndServe(httpBinding, nil); err != nil {
			server.logger.Fatal("Server failed with error: %v", err)
		}

		server.isRunning = false
	}()

}

func (server *Server) getHttpBinding() string {

	// validate the port
	port := server.config.Server.Http.Port
	if port < 1 || port > math.MaxUint16 {
		panic(fmt.Sprintf("%q is an invalid value for a port. Ports can only be in the range of %v to %v,", port, 1, math.MaxUint16))
	}

	return fmt.Sprintf(":%v", port)
}
