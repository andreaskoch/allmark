package handlers

import (
	"allmark.io/modules/common/util/fsutil"
	"allmark.io/modules/common/util/hashutil"
	"allmark.io/modules/web/header"
	"net/http"
	"path/filepath"
	"strings"
)

// Static creates a static http handler for the given directory.
func Static(directory, prefix string) http.Handler {
	return http.StripPrefix(prefix, http.FileServer(http.Dir(directory)))
}

// AddETAgToStaticFileHandler wraps the given staticFileHandler and adds an ETag header.
func AddETAgToStaticFileHandler(staticFileHandler http.Handler, headerWriter header.HeaderWriter, baseFolder, requestPrefixToStripFromRequestURI string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		headerWriter.Write(w, "")

		// determine the hash
		etag := ""

		// prepare the request uri
		requestURI := r.RequestURI
		if requestPrefixToStripFromRequestURI != "" {
			requestURI = stripPathFromRequest(r, requestPrefixToStripFromRequestURI)
		}

		// assemble the filepath on disc
		filePath := filepath.Join(baseFolder, requestURI)

		// read the the hash
		if file, err := fsutil.OpenFile(filePath); err == nil {
			defer file.Close()
			if fileHash, hashErr := hashutil.GetHash(file); hashErr == nil {
				etag = fileHash
			}
		}
		if etag != "" {
			header.ETag(w, etag)
		}

		staticFileHandler.ServeHTTP(w, r)
	})
}

func stripPathFromRequest(request *http.Request, path string) string {
	return strings.TrimPrefix(request.RequestURI, path)
}
