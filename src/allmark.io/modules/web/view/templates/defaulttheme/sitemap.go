// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package defaulttheme

import (
	"allmark.io/modules/web/view/templates/templatenames"
	"fmt"
)

func init() {
	templates[templatenames.Sitemap] = sitemapTemplate
	templates[templatenames.SitemapEntry] = sitemapContentTemplate
}

const SitemapChildPlaceholder = `@@CHILDS@@`

const sitemapTemplate = `
<header>
<h1 class="title">
{{.Title}}
</h1>
</header>

<section class="description">
{{.Description}}
</section>

<section class="content">
<ul class="tree">
{{ .Tree }}
</ul>
</section>
`

var sitemapContentTemplate = fmt.Sprintf(`<li>
	<a href="{{.Path}}" {{ if .Description }}title="{{.Description}}"{{ end }}>{{.Title}}</a>

	{{ if .Childs }}
	<ul>
		%s
	</ul>
	{{ end }}
</li>`, SitemapChildPlaceholder)
