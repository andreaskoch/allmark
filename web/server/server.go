// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/dataaccess"
	"github.com/andreaskoch/allmark2/services/converter"
	"github.com/andreaskoch/allmark2/services/parser"
	"github.com/andreaskoch/allmark2/web/orchestrator"
	"github.com/andreaskoch/allmark2/web/server/handler"
	"github.com/andreaskoch/allmark2/web/webpaths"
	"github.com/gorilla/mux"
	"math"
	"net/http"
)

var (
	BasePath      = "/"
	TagPathPrefix = fmt.Sprintf("%stags.html#", BasePath)

	// Dynamic Routes
	PrintHandlerRoute  = `/{path:.+\.print$|print$}`
	JsonHandlerRoute   = `/{path:.+\.json$|json$}`
	RtfHandlerRoute    = `/{path:.+\.rtf$|rtf$}`
	UpdateHandlerRoute = `/{path:.+\.ws$|ws$}`

	ItemHandlerRoute = "/{path:.*$}"

	TagmapHandlerRoute                = "/tags.html"
	SitemapHandlerRoute               = "/sitemap.html"
	XmlSitemapHandlerRoute            = "/sitemap.xml"
	RssHandlerRoute                   = "/feed.rss"
	RobotsTxtHandlerRoute             = "/robots.txt"
	SearchHandlerRoute                = "/search"
	OpenSearchDescriptionHandlerRoute = "/opensearch.xml"

	TypeAheadSearchHandlerRoute = "/search.json"
	TypeAheadTitlesHandlerRoute = "/titles.json"

	// Static Routes
	ThemeFolderRoute = "/theme"
)

func New(logger logger.Logger, config config.Config, repository dataaccess.Repository, parser parser.Parser, converter converter.Converter) (*Server, error) {

	// paths
	patherFactory := webpaths.NewFactory(logger, repository)
	itemPathProvider := patherFactory.Absolute(BasePath)
	tagPathProvider := patherFactory.Absolute(TagPathPrefix)
	webPathProvider := webpaths.NewWebPathProvider(patherFactory, itemPathProvider, tagPathProvider)

	// orchestrator
	orchestratorFactory := orchestrator.NewFactory(logger, repository, parser, converter, webPathProvider)

	// handlers
	handlerFactory := handler.NewFactory(logger, config, *orchestratorFactory)

	return &Server{
		logger: logger,
		config: config,

		handlerFactory: handlerFactory,
	}, nil

}

type Server struct {
	isRunning bool

	logger logger.Logger
	config config.Config

	handlerFactory *handler.Factory
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

		// websocket update handler
		// updateHub := update.NewHub(server.logger, server.updateHub)
		// go updateHub.Run()

		updateHandler := server.handlerFactory.NewUpdateHandler()
		requestRouter.Handle(UpdateHandlerRoute, websocket.Handler(updateHandler.Func()))

		// serve auxiliary dynamic files
		requestRouter.HandleFunc(RobotsTxtHandlerRoute, gzipReponse(server.handlerFactory.NewRobotsTxtHandler().Func()))
		requestRouter.HandleFunc(XmlSitemapHandlerRoute, gzipReponse(server.handlerFactory.NewXmlSitemapHandler().Func()))
		requestRouter.HandleFunc(TagmapHandlerRoute, gzipReponse(server.handlerFactory.NewTagsHandler().Func()))
		requestRouter.HandleFunc(SitemapHandlerRoute, gzipReponse(server.handlerFactory.NewSitemapHandler().Func()))
		requestRouter.HandleFunc(RssHandlerRoute, gzipReponse(server.handlerFactory.NewRssHandler().Func()))
		requestRouter.HandleFunc(PrintHandlerRoute, gzipReponse(server.handlerFactory.NewPrintHandler().Func()))
		requestRouter.HandleFunc(SearchHandlerRoute, gzipReponse(server.handlerFactory.NewSearchHandler().Func()))
		requestRouter.HandleFunc(OpenSearchDescriptionHandlerRoute, gzipReponse(server.handlerFactory.NewOpenSearchDescriptionHandler().Func()))
		requestRouter.HandleFunc(TypeAheadSearchHandlerRoute, gzipReponse(server.handlerFactory.NewTypeAheadSearchHandler().Func()))
		requestRouter.HandleFunc(TypeAheadTitlesHandlerRoute, gzipReponse(server.handlerFactory.NewTypeAheadTitlesHandler().Func()))

		// serve static files
		if themeFolder := server.config.ThemeFolder(); fsutil.DirectoryExists(themeFolder) {
			s := http.StripPrefix(ThemeFolderRoute, http.FileServer(http.Dir(themeFolder)))
			requestRouter.PathPrefix(ThemeFolderRoute).Handler(s)
		}

		// serve items
		requestRouter.HandleFunc(RtfHandlerRoute, server.handlerFactory.NewRtfHandler().Func())
		requestRouter.HandleFunc(JsonHandlerRoute, gzipReponse(server.handlerFactory.NewJsonHandler().Func()))
		requestRouter.HandleFunc(ItemHandlerRoute, gzipReponse(server.handlerFactory.NewItemHandler().Func()))

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
