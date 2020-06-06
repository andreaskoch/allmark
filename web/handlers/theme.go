// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"github.com/elWyatt/allmark/common/util/hashutil"
	"github.com/elWyatt/allmark/web/header"
	"github.com/elWyatt/allmark/web/view/themes"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

// InMemoryTheme creates a theme-handler that serves the theme-files from memory.
func InMemoryTheme(themeFolderPath string, headerWriter header.HeaderWriter, error404Handler http.Handler) http.Handler {

	defaultTheme := themes.GetTheme()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path
		path = strings.TrimPrefix(path, themeFolderPath)

		themeFile := defaultTheme.Get(path)
		if themeFile == nil {

			// display a 404 error page
			error404Handler.ServeHTTP(w, r)
			return

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
