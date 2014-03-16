// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlerutil

import (
	"io"
	"text/template"
)

func RenderTemplate(model interface{}, template *template.Template, writer io.Writer) error {
	return template.Execute(writer, model)
}
