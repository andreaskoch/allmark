// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"fmt"
	"github.com/andreaskoch/allmark/templates"
	"github.com/andreaskoch/allmark/view"
	"io"
	"os"
)

func (renderer *Renderer) Error404(writer io.Writer) {

	// get the 404 page template
	templateType := templates.ErrorTemplateName
	template, err := renderer.templateProvider.GetFullTemplate(templateType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "No template of type %s found.", templateType)
		return
	}

	// create a error view model
	title := "Not found"
	content := fmt.Sprintf("The requested item was not found.")
	errorModel := view.Error(title, content, renderer.root.RelativePath, renderer.root.AbsolutePath)

	// attach the toplevel navigation
	errorModel.ToplevelNavigation = renderer.root.ToplevelNavigation

	// attach the bread crumb navigation
	errorModel.BreadcrumbNavigation = renderer.root.BreadcrumbNavigation

	// render the template
	writeTemplate(errorModel, template, writer)
}
