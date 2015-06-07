// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"allmark.io/modules/common/util/hashutil"
	"allmark.io/modules/web/header"
	"allmark.io/modules/web/view/themes"
	"fmt"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
	"path/filepath"
)

// InMemoryTheme creates a theme-handler that serves the theme-files from memory.
func InMemoryTheme(headerWriter header.HeaderWriter, error404Handler http.Handler) http.Handler {

	defaultTheme := themes.GetTheme()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get the path from the request variables
		vars := mux.Vars(r)
		path := vars["path"]

		themeFile := defaultTheme.Get(path)
		if themeFile == nil {

			// display a 404 error page
			error404Handler.ServeHTTP(w, r)

		}

		// detect the mime type
		data := themeFile.Data()
		mimeType := getMimeType(path, data)

		// etag
		etag := hashutil.FromBytes(data)

		// set headers
		headerWriter.Write(w, mimeType)
		if etag != "" {
			header.ETag(w, etag)
		}

		fmt.Fprintf(w, `%s`, data)
	})
}

// getMimeType derives the mime-type from the given URI and data.
func getMimeType(uri string, data []byte) string {
	extention := filepath.Ext(uri)
	mimeType := mime.TypeByExtension(extention)

	// fallback mimetype detection
	if mimeType == "" {
		mimeType = http.DetectContentType(data)

	}

	return mimeType
}
