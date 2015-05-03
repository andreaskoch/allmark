// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/util/hashutil"
	"allmark.io/modules/web/server/header"
	"allmark.io/modules/web/view/themes"
	"fmt"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
	"path/filepath"
)

type Theme struct {
	logger          logger.Logger
	headerWriter    header.HeaderWriter
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

		// etag
		etag := hashutil.FromBytes(data)

		// set headers
		handler.headerWriter.Write(w, mimeType)
		if etag != "" {
			header.ETag(w, etag)
		}

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
