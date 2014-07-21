// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtfhandler

import (
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/services/conversion"
	"github.com/andreaskoch/allmark2/ui/web/orchestrator"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/errorhandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/handlerutil"
	"github.com/andreaskoch/allmark2/ui/web/view/templates"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
	"io"
	"net/http"
)

func New(logger logger.Logger, config *config.Config, itemIndex *index.Index, patherFactory paths.PatherFactory, converter conversion.Converter) *RtfHandler {

	// templates
	templateProvider := templates.NewProvider(config.TemplatesFolder())

	// error
	error404Handler := errorhandler.New(logger, config, itemIndex, patherFactory)

	// viewmodel
	conversionModelOrchestrator := orchestrator.NewConversionModelOrchestrator(itemIndex, converter)

	return &RtfHandler{
		logger:                      logger,
		itemIndex:                   itemIndex,
		config:                      config,
		patherFactory:               patherFactory,
		templateProvider:            templateProvider,
		error404Handler:             error404Handler,
		conversionModelOrchestrator: conversionModelOrchestrator,
	}
}

type RtfHandler struct {
	logger                      logger.Logger
	itemIndex                   *index.Index
	config                      *config.Config
	patherFactory               paths.PatherFactory
	templateProvider            *templates.Provider
	error404Handler             *errorhandler.ErrorHandler
	conversionModelOrchestrator orchestrator.ConversionModelOrchestrator
}

func (handler *RtfHandler) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// get the request route
		requestRoute, err := handlerutil.GetRouteFromRequest(r)
		if err != nil {
			handler.logger.Error("Unable to get route from request. Error: %s", err)
			return
		}

		// make sure the request body is closed
		defer r.Body.Close()

		// stage 1: check if there is a item for the request
		if item, found := handler.itemIndex.IsMatch(*requestRoute); found {

			// render the view model
			pathProvider := handler.patherFactory.Relative(item.Route())
			viewModel := handler.conversionModelOrchestrator.GetConversionModel(pathProvider, item)
			handler.render(w, viewModel)
			return
		}

		// display a 404 error page
		error404Handler := handler.error404Handler.Func()
		error404Handler(w, r)
	}
}

func (handler *RtfHandler) render(writer io.Writer, viewModel viewmodel.ConversionModel) {

	// get a template
	template, err := handler.templateProvider.GetFullTemplate(viewModel.Type)
	if err != nil {
		handler.logger.Error("No template for item of type %q.", viewModel.Type)
		return
	}

	// render template
	if err := handlerutil.RenderTemplate(viewModel, template, writer); err != nil {
		handler.logger.Error("%s", err)
		return
	}

}
