// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/content"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/paths/webpaths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/conversion"
	"github.com/andreaskoch/allmark2/ui/web/server/handler"
	"github.com/gorilla/mux"
	"math"
	"net/http"
	"os"
	"path/filepath"
)

const (

	// Dynamic Routes
	ItemHandlerRoute       = "/{path:.*}"
	TagmapHandlerRoute     = "/tags.html"
	SitemapHandlerRoute    = "/sitemap.html"
	XmlSitemapHandlerRoute = "/sitemap.xml"
	RssHandlerRoute        = "/feed.rss"
	RobotsTxtHandlerRoute  = "/robots.txt"
	DebugHandlerRoute      = "/debug/index"
	WebSocketHandlerRoute  = "/ws"

	// Static Routes
	ThemeFolderRoute = "/theme/"
)

func New(logger logger.Logger, config *config.Config, converter conversion.Converter) (*Server, error) {

	itemIndex := index.CreateItemIndex(logger)
	fileIndex := index.CreateFileIndex(logger)
	patherFactory := webpaths.NewFactory(logger, itemIndex)

	return &Server{
		config:        config,
		logger:        logger,
		patherFactory: patherFactory,
		converter:     converter,
		itemIndex:     itemIndex,
		fileIndex:     fileIndex,
	}, nil
}

type Server struct {
	isRunning bool

	config        *config.Config
	logger        logger.Logger
	patherFactory paths.PatherFactory
	converter     conversion.Converter
	itemIndex     *index.ItemIndex
	fileIndex     *index.FileIndex
}

func (server *Server) ServeItem(item *model.Item) {
	server.logger.Debug("Serving item %q", item)
	server.itemIndex.Add(item)
}

func (server *Server) ServeFolder(baseFolder, folderPath string) {
	if !fsutil.DirectoryExists(folderPath) {
		return
	}

	parentRoute, err := route.NewFromRequest("")
	if err != nil {
		panic(err)
	}

	filepath.Walk(folderPath, func(folderEntryPath string, folderEntryInfo os.FileInfo, err error) error {

		if folderEntryInfo.IsDir() {
			return nil
		}

		fileRoute, err := route.NewFromFilePath(baseFolder, folderEntryPath)
		if err != nil {
			return err
		}

		file, err := model.NewFromPath(fileRoute, parentRoute, content.FileProvider(folderEntryPath, fileRoute))
		if err != nil {
			return err
		}

		server.fileIndex.Add(file)

		return nil
	})
}

func (server *Server) IsRunning() bool {
	return server.isRunning
}

func (server *Server) Start() chan error {
	result := make(chan error)

	go func() {
		server.isRunning = true

		// register requst routers
		requestRouter := mux.NewRouter()
		requestRouter.HandleFunc(RobotsTxtHandlerRoute, handler.NewRobotsTxtHandler(server.logger, server.config, server.itemIndex, server.patherFactory).Func())
		requestRouter.HandleFunc(XmlSitemapHandlerRoute, handler.NewXmlSitemapHandler(server.logger, server.config, server.itemIndex, server.patherFactory).Func())
		requestRouter.HandleFunc(SitemapHandlerRoute, handler.NewSitemapHandler(server.logger, server.config, server.itemIndex, server.patherFactory).Func())
		requestRouter.HandleFunc(DebugHandlerRoute, handler.NewDebugHandler(server.logger, server.itemIndex, server.fileIndex).Func())
		requestRouter.HandleFunc(RssHandlerRoute, handler.NewRssHandler(server.logger, server.config, server.itemIndex, server.fileIndex, server.patherFactory, server.converter).Func())
		requestRouter.HandleFunc(ItemHandlerRoute, handler.NewItemHandler(server.logger, server.config, server.itemIndex, server.fileIndex, server.patherFactory, server.converter).Func())

		// start http server: http
		httpBinding := server.getHttpBinding()
		server.logger.Info("Starting http server %q\n", httpBinding)

		if err := http.ListenAndServe(httpBinding, requestRouter); err != nil {
			result <- fmt.Errorf("Server failed with error: %v", err)
		} else {
			result <- nil
		}

		server.isRunning = false
	}()

	return result
}

func (server *Server) getHttpBinding() string {

	// validate the port
	port := server.config.Server.Http.Port
	if port < 1 || port > math.MaxUint16 {
		panic(fmt.Sprintf("%q is an invalid value for a port. Ports can only be in the range of %v to %v,", port, 1, math.MaxUint16))
	}

	return fmt.Sprintf(":%v", port)
}
