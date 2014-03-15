// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

import (
	"fmt"
)

var sitemapContentTemplate = fmt.Sprintf(`
<li>
	<a href="{{.Path}}" {{ if .Description }}title="{{.Description}}"{{ end }}>{{.Title}}</a>

	{{ if .Childs }}	
	<ol>
	%s
	</ol>
	{{ end }}
</li>`, ChildTemplatePlaceholder)

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
<ol>
{{.Content}}
</ol>
</section>
`
