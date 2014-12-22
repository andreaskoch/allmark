// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/web/server/header"
	"github.com/andreaskoch/allmark2/web/view/themes"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
	"path/filepath"
)

type Theme struct {
	logger logger.Logger

	error404Handler Handler
}

func (handler *Theme) Func() func(w http.ResponseWriter, r *http.Request) {

	defaultTheme := themes.GetTheme()

	return func(w http.ResponseWriter, r *http.Request) {

		// get the path from the request variables
		vars := mux.Vars(r)
		path := vars["path"]

		themeFile := defaultTheme.Get(path)
		if themeFile == nil {

			// display a 404 error page
			error404Handler := handler.error404Handler.Func()
			error404Handler(w, r)

		}

		// detect the mime type
		data := themeFile.Data()
		mimeType := getMimeType(path, data)

		// set headers
		header.ContentType(w, r, mimeType)
		header.Cache(w, r, header.STATICCONTENT_CACHEDURATION_SECONDS)
		fmt.Fprintf(w, `%s`, data)
	}
}

func getMimeType(uri string, data []byte) string {
	extention := filepath.Ext(uri)
	mimeType := mime.TypeByExtension(extention)

	// fallback mimetype detection
	if mimeType == "" {
		mimeType = http.DetectContentType(data)

	}

	return mimeType
}
