// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package defaulttheme

import (
	"allmark.io/modules/web/view/templates/templatenames"
)

func init() {
	templates[templatenames.OpenSearchDescription] = openSearchDescriptionTemplate
}

const openSearchDescriptionTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<OpenSearchDescription xmlns="http://a9.com/-/spec/opensearch/1.1/">
  <ShortName>{{.Title}}</ShortName>
  <Description>{{.Description}}</Description>
  <Tags>{{.Tags}}</Tags>
  <Contact />
  <Image height="16" width="16" type="image/x-icon">{{.FavIconURL}}</Image>
  <URL type="text/html" template="{{.SearchURL}}" />
  <OutputEncoding>UTF-8</OutputEncoding>
  <InputEncoding>UTF-8</InputEncoding>
</OpenSearchDescription>
`
