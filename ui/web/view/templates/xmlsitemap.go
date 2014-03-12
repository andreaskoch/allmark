// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

import (
	"fmt"
)

var xmlSitemapTemplate = fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
%s
</urlset>`, ChildTemplatePlaceholder)

var xmlSitemapContentTemplate = `<url>
<loc>{{ .Loc }}</loc>
<lastmod>{{ .LastModified }}</lastmod>
</url>`
