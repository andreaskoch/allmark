// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"io"
	"text/template"
)

func renderTemplate(model interface{}, template *template.Template, writer io.Writer) error {
	return template.Execute(writer, model)
}
