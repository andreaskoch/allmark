// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"allmark.io/modules/web/header"
	"allmark.io/modules/web/view/templates"
	"allmark.io/modules/web/view/viewmodel"
	"fmt"
	"net/http"
)

// RobotsTxt creates a http handler for serving the robots.txt.
func RobotsTxt(headerWriter header.HeaderWriter, templateProvider templates.Provider) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// template
		baseUrl := getBaseUrlFromRequest(r)
		robotsTxtTemplate, err := templateProvider.GetSubTemplate(baseUrl, templates.RobotsTxtTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		// view model
		disallowPaths := []string{
			"/thumbnails",
			"/rtf$",
			"/json$",
			"/print$",
			"/ws$",
			"/*.rtf$",
			"/*.json$",
			"/*.print$",
			"/*.ws$",
		}
		disallow := viewmodel.RobotsTxtDisallow{
			UserAgent: "*",
			Paths:     disallowPaths,
		}

		sitemapUrl := fmt.Sprintf("%s%s", baseUrl, XMLSitemapHandlerRoute)
		model := viewmodel.RobotsTxt{
			Disallows:  []viewmodel.RobotsTxtDisallow{disallow},
			SitemapUrl: sitemapUrl,
		}

		// write
		headerWriter.Write(w, header.CONTENTTYPE_TEXT)
		renderTemplate(robotsTxtTemplate, model, w)
	})

}
