// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

const repositoryTemplate = `
<header>
<h1 class="title" itemprop="name">
{{.Title}}
</h1>
</header>

<section class="description" itemprop="description">
{{.Description}}
</section>

{{if .Author.Name}}
<section class="publisher">
by <span class="author" itemprop="author" rel="author"><a href="{{ .Author.Url }}" title="{{ .Author.Name }}" target="_blank">{{ .Author.Name }}</a></span>
</section>
{{end}}

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
`
