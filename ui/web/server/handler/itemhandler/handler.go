// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package itemhandler

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/conversion"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/handlerutil"
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

		// stage 1: check if there is a item for the request
		if item, found := handler.index.IsMatch(*requestRoute); found {

			// create the view model
			pathProvider := handler.patherFactory.Relative(item.Route())
			viewModel := getViewModel(handler.index, pathProvider, handler.converter, item)

			// render the view model
			render(w, viewModel)

		}

		// stage 2: check if there is a file for the request
		if file, found := handler.index.IsFileMatch(*requestRoute); found {
			contentProvider := file.ContentProvider()

			// read the file data
			data, err := contentProvider.Data()
			if err != nil {
				return
			}

			fmt.Fprintf(w, "%s", data)
			return
		}

		fmt.Fprintln(w, "item not found.")
		return
	}
}

func getViewModel(index *index.Index, pathProvider paths.Pather, converter conversion.Converter, item *model.Item) viewmodel.Model {

	// convert content
	convertedContent, err := converter.Convert(pathProvider, item)
	if err != nil {
		return viewmodel.Model{}
	}

	// create a view model
	viewModel := viewmodel.Model{
		Base: getBaseModel(item),

		Content: convertedContent,

		Childs:               getChildModels(index, item),
		ToplevelNavigation:   getToplevelNavigation(index),
		BreadcrumbNavigation: getBreadcrumbNavigation(index, item),
		TopDocs:              getTopDocuments(index, pathProvider, converter, item),
	}

	return viewModel
}

func getTopDocuments(index *index.Index, pathProvider paths.Pather, converter conversion.Converter, item *model.Item) []*viewmodel.Model {
	childModels := make([]*viewmodel.Model, 0)

	childItems := index.GetChilds(item.Route())

	for _, childItem := range childItems {

		// skip virtual items
		if childItem.IsVirtual() {
			continue
		}

		// todo: choose the right level, don't just take all childs
		childModel := getViewModel(index, pathProvider, converter, childItem)
		childModels = append(childModels, &childModel)
	}

	return childModels
}

func getBaseModel(item *model.Item) viewmodel.Base {
	return viewmodel.Base{
		Type:  item.Type.String(),
		Route: item.Route().Value(),
		Level: item.Route().Level(),

		Title:       item.Title,
		Description: item.Description,
	}
}

func getChildModels(index *index.Index, item *model.Item) []*viewmodel.Base {
	childModels := make([]*viewmodel.Base, 0)

	childItems := index.GetChilds(item.Route())
	for _, childItem := range childItems {
		baseModel := getBaseModel(childItem)
		childModels = append(childModels, &baseModel)
	}

	return childModels
}

func getToplevelNavigation(index *index.Index) *viewmodel.ToplevelNavigation {
	root, err := route.NewFromRequest("")
	if err != nil {
		return nil
	}

	toplevelEntries := make([]*viewmodel.ToplevelEntry, 0)
	for _, child := range index.GetChilds(root) {

		// skip all childs which are not first level
		if child.Route().Level() != 1 {
			continue
		}

		toplevelEntries = append(toplevelEntries, &viewmodel.ToplevelEntry{
			Title: child.Title,
			Path:  child.Route().Value(),
		})

	}

	return &viewmodel.ToplevelNavigation{
		Entries: toplevelEntries,
	}
}

func getBreadcrumbNavigation(index *index.Index, item *model.Item) *viewmodel.BreadcrumbNavigation {

	// create a new bread crumb navigation
	navigation := &viewmodel.BreadcrumbNavigation{
		Entries: make([]*viewmodel.Breadcrumb, 0),
	}

	// abort if item or model is nil
	if item == nil {
		return navigation
	}

	// recurse if there is a parent
	if parent := index.GetParent(item.Route()); parent != nil {
		navigation.Entries = append(navigation.Entries, getBreadcrumbNavigation(index, parent).Entries...)
	}

	// append a new navigation entry and return it
	navigation.Entries = append(navigation.Entries, &viewmodel.Breadcrumb{
		Title: item.Title,
		Level: item.Route().Level(),
		Path:  item.Route().Value(),
	})

	return navigation
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
