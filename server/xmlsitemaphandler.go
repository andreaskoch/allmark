// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"github.com/andreaskoch/allmark/util"
	"net/http"
)

var xmlSitemapHandler = func(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Fprintln(w, `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)

	for _, item := range items {
		route := item.AbsoluteRoute
		location := fmt.Sprintf(`http://%s/%s`, r.Host, util.EncodeUrl(route))
		lastMod := item.Date

		fmt.Fprintln(w, `<url>`)
		fmt.Fprintln(w, fmt.Sprintf(`<loc>%s</loc>`, location))
		fmt.Fprintln(w, fmt.Sprintf(`<lastmod>%s</lastmod>`, lastMod))
		fmt.Fprintln(w, `</url>`)
	}

	fmt.Fprintln(w, `</urlset>`)
}
