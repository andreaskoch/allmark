// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package header

import (
	"fmt"
	"net/http"
)

const (
	CONTENTTYPE_HTML = "text/html; charset=utf-8"
	CONTENTTYPE_TEXT = "text/plain; charset=utf-8"
	CONTENTTYPE_XML  = "text/xml; charset=utf-8"
	CONTENTTYPE_JSON = "application/json; charset=utf-8"
	CONTENTTYPE_DOCX = "application/vnd.openxmlformats-officedocument.wordprocessingml.document; charset=utf-8"
)

func Cache(w http.ResponseWriter, seconds int) {
	w.Header().Add("Cache-Control", fmt.Sprintf("public, max-age=%d", seconds))
}

func ETag(w http.ResponseWriter, hash string) {
	if hash == "" {
		return
	}

	w.Header().Add("ETag", hash)
}

func NoCache(w http.ResponseWriter) {
	w.Header().Add("Cache-Control", "no-cache")
}

func ContentType(w http.ResponseWriter, contentType string) {
	if contentType == "" {
		return
	}

	w.Header().Set("Content-Type", contentType)
}

func VaryAcceptEncoding(w http.ResponseWriter) {
	w.Header().Set("Vary", "Accept-Encoding")
}

// configurable header writer
type configurableHeaderWriter struct {
	cacheDuration int
}

func (headerWriter configurableHeaderWriter) Write(w http.ResponseWriter, contentType string) {
	Cache(w, headerWriter.cacheDuration)
	ContentType(w, contentType)
	VaryAcceptEncoding(w)
}

// no-cache header writer
type noCacheHeaderWriter struct {
}

func (headerWriter noCacheHeaderWriter) Write(w http.ResponseWriter, contentType string) {
	NoCache(w)
	ContentType(w, contentType)
	VaryAcceptEncoding(w)
}

type HeaderWriter interface {
	Write(w http.ResponseWriter, contentType string)
}

// HeaderWriter factory
func NewHeaderWriterFactory(reindexIntervalInSeconds int) WriterFactory {

	// default cache durations
	cacheDurationDynamic := 86400   // 1 day
	cacheDurationStatic := 31536000 // 1 year

	reindexingIsEnabled := reindexIntervalInSeconds > 0
	if reindexingIsEnabled {

		// cache durations based on the reindex interval
		cacheDurationDynamic = reindexIntervalInSeconds / 2
		cacheDurationStatic = reindexIntervalInSeconds / 2

	}

	// create the header writers
	static := configurableHeaderWriter{cacheDurationStatic}
	dynamic := configurableHeaderWriter{cacheDurationDynamic}
	noCache := noCacheHeaderWriter{}

	// create the factory with the given parameters
	return WriterFactory{
		static,
		dynamic,
		noCache,
	}

}

type WriterFactory struct {
	static  HeaderWriter
	dynamic HeaderWriter
	noCache HeaderWriter
}

func (writerFactory *WriterFactory) Static() HeaderWriter {
	return writerFactory.static
}

func (writerFactory *WriterFactory) Dynamic() HeaderWriter {
	return writerFactory.dynamic
}

func (writerFactory *WriterFactory) NoCache() HeaderWriter {
	return writerFactory.noCache
}
