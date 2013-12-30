// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package itemhandler

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/services/conversion"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/handlerutil"
	"github.com/andreaskoch/allmark2/ui/web/server/index"
	"github.com/andreaskoch/allmark2/ui/web/view/templates"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
	"io"
	"net/http"
)

func New(logger logger.Logger, index *index.Index, patherFactory paths.PatherFactory, converter conversion.Converter) *ItemHandler {
	return &ItemHandler{
		logger:        logger,
		index:         index,
		patherFactory: patherFactory,
		converter:     converter,
	}
}

type ItemHandler struct {
	logger        logger.Logger
	index         *index.Index
	patherFactory paths.PatherFactory
	converter     conversion.Converter
}

func (handler *ItemHandler) Func() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// get the request route
		requestPath := handlerutil.GetRequestedPathFromRequest(r)
		requestRoute, err := route.NewFromRequest(requestPath)
		if err != nil {
			fmt.Fprintln(w, "%s", err)
			return
		}

		// make sure the request body is closed
		defer r.Body.Close()

		// check if there is a item for the request
		item, found := handler.index.IsMatch(*requestRoute)
		if !found {
			fmt.Fprintln(w, "item not found")
			return
		}

		// Parent
		parent := handler.index.GetParent(item)
		if parent != nil {
			fmt.Fprintf(w, "Parent: %s\n", parent.Title)
		}

		// convert content
		pathProvider := handler.patherFactory.Relative()
		convertedContent, err := handler.converter.Convert(pathProvider, item)

		if err != nil {
			fmt.Fprintln(w, "Unable to convert content. Error: %s", err)
			return
		}

		// create a view model
		viewModel := viewmodel.Model{
			Type:        item.Type.String(),
			Title:       item.Title,
			Description: item.Description,
			Content:     convertedContent,
		}

		render(w, viewModel)

		// Childs
		childs := handler.index.GetChilds(item)
		for _, child := range childs {
			fmt.Fprintf(w, "Child: %s\n", child.Title)
		}
	}
}

func render(writer io.Writer, viewModel viewmodel.Model) {

	templateProvider := templates.NewProvider(".")

	// get a template
	if template, err := templateProvider.GetFullTemplate(viewModel.Type); err == nil {

		err := template.Execute(writer, viewModel)
		if err != nil {
			fmt.Println(err)
		}

	} else {

		fmt.Fprintf(writer, "No template for item of type %q.", viewModel.Type)

	}

}
