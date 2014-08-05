// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package printhandler

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/services/conversion"
	"github.com/andreaskoch/allmark2/ui/web/orchestrator"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/errorhandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/handlerutil"
	"github.com/andreaskoch/allmark2/ui/web/view/templates"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func New(logger logger.Logger, config *config.Config, itemIndex *index.Index, patherFactory paths.PatherFactory, converter conversion.Converter) *PrintHandler {

	// templates
	templateProvider := templates.NewProvider(config.TemplatesFolder())

	// error
	error404Handler := errorhandler.New(logger, config, itemIndex, patherFactory)

	// viewmodel
	conversionModelOrchestrator := orchestrator.NewConversionModelOrchestrator(itemIndex, converter)

	return &PrintHandler{
		logger:                      logger,
		itemIndex:                   itemIndex,
		config:                      config,
		patherFactory:               patherFactory,
		templateProvider:            templateProvider,
		error404Handler:             error404Handler,
		conversionModelOrchestrator: conversionModelOrchestrator,
	}
}

type PrintHandler struct {
	logger                      logger.Logger
	itemIndex                   *index.Index
	config                      *config.Config
	patherFactory               paths.PatherFactory
	templateProvider            *templates.Provider
	error404Handler             *errorhandler.ErrorHandler
	conversionModelOrchestrator orchestrator.ConversionModelOrchestrator
}

func (handler *PrintHandler) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// get the path from the request variables
		vars := mux.Vars(r)
		path := vars["path"]

		// get the request route
		requestRoute, err := route.NewFromRequest(path)
		if err != nil {
			handler.logger.Error("Unable to get route from request. Error: %s", err.Error())
			return
		}

		// make sure the request body is closed
		defer r.Body.Close()

		// check if there is a item for the request
		item, found := handler.itemIndex.IsMatch(*requestRoute)
		if !found {

			// display a 404 error page
			error404Handler := handler.error404Handler.Func()
			error404Handler(w, r)
			return
		}

		// prepare a path provider which includes the hostname
		hostname := handlerutil.GetHostnameFromRequest(r)
		addressPrefix := fmt.Sprintf("http://%s/", hostname)
		pathProvider := handler.patherFactory.Absolute(addressPrefix)

		// render the view model
		viewModel := handler.conversionModelOrchestrator.GetConversionModel(pathProvider, item)
		handler.render(w, viewModel)
	}
}

func (handler *PrintHandler) render(writer io.Writer, viewModel viewmodel.ConversionModel) {

	// get a template
	template, err := handler.templateProvider.GetSubTemplate(templates.ConversionTemplateName)
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
