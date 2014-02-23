// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package itemhandler

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
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

func New(logger logger.Logger, config *config.Config, itemIndex *index.ItemIndex, fileIndex *index.FileIndex, patherFactory paths.PatherFactory, converter conversion.Converter) *ItemHandler {
	return &ItemHandler{
		logger:        logger,
		itemIndex:     itemIndex,
		fileIndex:     fileIndex,
		config:        config,
		patherFactory: patherFactory,
		converter:     converter,
	}
}

type ItemHandler struct {
	logger        logger.Logger
	itemIndex     *index.ItemIndex
	fileIndex     *index.FileIndex
	config        *config.Config
	patherFactory paths.PatherFactory
	converter     conversion.Converter
}

func (handler *ItemHandler) Func() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// get the request route
		requestRoute, err := handlerutil.GetRouteFromRequest(r)
		if err != nil {
			fmt.Fprintln(w, "%s", err)
			return
		}

		// make sure the request body is closed
		defer r.Body.Close()

		// stage 1: check for theme files
		themeRoute, err := route.NewFromRequest("theme")
		if err != nil {
			fmt.Fprintln(w, "%s", err)
			return
		}

		if isThemeFile := requestRoute.IsChildOf(themeRoute); isThemeFile {

			if file, found := handler.fileIndex.IsMatch(*requestRoute); found {
				fileContentProvider := file.ContentProvider()
				data, err := fileContentProvider.Data()
				if err != nil {
					return
				}

				fmt.Fprintf(w, "%s", data)
				return
			}
		}

		// stage 2: check if there is a item for the request
		if item, found := handler.itemIndex.IsMatch(*requestRoute); found {

			// create the view model
			pathProvider := handler.patherFactory.Relative(item.Route())
			viewModel := getViewModel(handler.itemIndex, pathProvider, handler.converter, item)

			// render the view model
			render(w, viewModel)
			return
		}

		// stage 3: check if there is a file for the request
		if file, found := handler.itemIndex.IsFileMatch(*requestRoute); found {
			contentProvider := file.ContentProvider()

			// read the file data
			data, err := contentProvider.Data()
			if err != nil {
				return
			}

			fmt.Fprintf(w, "%s", data)
			return
		}

		fmt.Fprintln(w, fmt.Sprintf("item %q not found.", requestRoute))
		return
	}
}

func getViewModel(itemIndex *index.ItemIndex, pathProvider paths.Pather, converter conversion.Converter, item *model.Item) viewmodel.Model {

	// convert content
	convertedContent, err := converter.Convert(pathProvider, item)
	if err != nil {
		return viewmodel.Model{}
	}

	// create a view model
	viewModel := viewmodel.Model{
		Base: getBaseModel(item),

		Content: convertedContent,

		Childs:               getChildModels(itemIndex, item),
		ToplevelNavigation:   getToplevelNavigation(itemIndex),
		BreadcrumbNavigation: getBreadcrumbNavigation(itemIndex, item),
		TopDocs:              getTopDocuments(itemIndex, pathProvider, converter, item),
	}

	return viewModel
}

func getTopDocuments(itemIndex *index.ItemIndex, pathProvider paths.Pather, converter conversion.Converter, item *model.Item) []*viewmodel.Model {
	childModels := make([]*viewmodel.Model, 0)

	childItems := itemIndex.GetChilds(item.Route())

	for _, childItem := range childItems {

		// skip virtual items
		if childItem.IsVirtual() {
			continue
		}

		// todo: choose the right level, don't just take all childs
		childModel := getViewModel(itemIndex, pathProvider, converter, childItem)
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

func getChildModels(itemIndex *index.ItemIndex, item *model.Item) []*viewmodel.Base {
	childModels := make([]*viewmodel.Base, 0)

	childItems := itemIndex.GetChilds(item.Route())
	for _, childItem := range childItems {
		baseModel := getBaseModel(childItem)
		childModels = append(childModels, &baseModel)
	}

	return childModels
}

func getToplevelNavigation(itemIndex *index.ItemIndex) *viewmodel.ToplevelNavigation {
	root, err := route.NewFromRequest("")
	if err != nil {
		return nil
	}

	toplevelEntries := make([]*viewmodel.ToplevelEntry, 0)
	for _, child := range itemIndex.GetChilds(root) {

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

func getBreadcrumbNavigation(itemIndex *index.ItemIndex, item *model.Item) *viewmodel.BreadcrumbNavigation {

	// create a new bread crumb navigation
	navigation := &viewmodel.BreadcrumbNavigation{
		Entries: make([]*viewmodel.Breadcrumb, 0),
	}

	// abort if item or model is nil
	if item == nil {
		return navigation
	}

	// recurse if there is a parent
	if parent := itemIndex.GetParent(item.Route()); parent != nil {
		navigation.Entries = append(navigation.Entries, getBreadcrumbNavigation(itemIndex, parent).Entries...)
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
