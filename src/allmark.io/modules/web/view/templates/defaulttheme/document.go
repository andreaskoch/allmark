// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package defaulttheme

import (
	"allmark.io/modules/web/view/templates/templatenames"
)

func init() {
	templates[templatenames.Document] = documentTemplate
}

const documentTemplate = `
<header>
<h1 class="title" itemprop="name">
{{.Title}}
</h1>
</header>

<section class="description" itemprop="description">
{{.Description}}
</section>

{{template "publisher-snippet" .}}

<section class="content" itemprop="articleBody">
{{.Content}}
</section>

<div class="cleaner"></div>

{{ if .Childs }}
<section class="preview">
	<ul>
	</ul>
</section>
{{end}}

{{template "aliases-snippet" .}}
{{template "tags-snippet" .}}
`
