// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"net/http"
)

var xmlSitemapHandler = func(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Fprint(w, `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)

	for route, item := range items {
		location := fmt.Sprintf(`http://%s/%s`, r.Host, route)
		lastMod := item.Date

		fmt.Fprint(w, `<url>`)
		fmt.Fprintf(w, `<loc>%s</loc>`, location)
		fmt.Fprintf(w, `<lastmod>%s</lastmod>`, lastMod)
		fmt.Fprint(w, `</url>`)
	}

	fmt.Fprint(w, `</urlset>`)
}
