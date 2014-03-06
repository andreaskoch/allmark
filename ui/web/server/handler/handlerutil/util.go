// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlerutil

import (
	"github.com/andreaskoch/allmark2/common/route"
	"io"
	"net/http"
	"text/template"
)

func GetRouteFromRequest(r *http.Request) (*route.Route, error) {
	return route.NewFromRequest(r.URL.Path)
}

func GetHostnameFromRequest(r *http.Request) string {
	return r.Host
}

func RenderTemplate(model interface{}, template *template.Template, writer io.Writer) error {
	return template.Execute(writer, model)
}
