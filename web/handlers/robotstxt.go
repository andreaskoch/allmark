// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"github.com/andreaskoch/allmark/web/header"
	"github.com/andreaskoch/allmark/web/view/templates"
	"github.com/andreaskoch/allmark/web/view/viewmodel"
	"fmt"
	"net/http"
)

// RobotsTxt creates a http handler for serving the robots.txt.
func RobotsTxt(headerWriter header.HeaderWriter, templateProvider templates.Provider) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// template
		baseURL := getBaseURLFromRequest(r)
		robotsTxtTemplate, err := templateProvider.GetRobotsTxtTemplate(baseURL)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		// view model
		disallowPaths := []string{
			"/thumbnails",
			"/docx$",
			"/json$",
			"/markdown$",
			"/print$",
			"/ws$",
			"/*.docx$",
			"/*.json$",
			"/*.markdown$",
			"/*.print$",
			"/*.ws$",
		}
		disallow := viewmodel.RobotsTxtDisallow{
			UserAgent: "*",
			Paths:     disallowPaths,
		}

		sitemapURL := fmt.Sprintf("%s%s", baseURL, XMLSitemapHandlerRoute)
		model := viewmodel.RobotsTxt{
			Disallows:  []viewmodel.RobotsTxtDisallow{disallow},
			SitemapURL: sitemapURL,
		}

		// write
		headerWriter.Write(w, header.CONTENTTYPE_TEXT)
		renderTemplate(robotsTxtTemplate, model, w)
	})

}
