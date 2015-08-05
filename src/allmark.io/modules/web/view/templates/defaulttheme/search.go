// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package defaulttheme

import (
	"allmark.io/modules/web/view/templates/templatenames"
)

func init() {
	templates[templatenames.Search] = searchTemplate
}

const searchTemplate = `
<header>
<h1 class="title">
{{.Title}}
</h1>
</header>

<section class="description">
{{.Description}}
</section>

{{ if .Results }}
{{ with .Results}}

<section class="content">
<nav>
	<form action="/search" method="GET">
		<input type="text" name="q" placeholder="search" value="{{.Query}}" autocomplete="off">
		<input type="submit" value="Search">
	</form>
</nav>

{{if .ResultCount}}
<header>
	Displaying {{.ResultCount}} of {{.TotalResultCount}} search results for "{{.Query}}":
</header>

<ol start="{{.StartIndex}}">
	{{ range .Results }}
	<li data-index="{{.Index}}">
			<a class="title" href="{{.Route}}">{{.Title}}</a>
			<p class="description">{{.Description}}</p>
			<span class="path">{{.Path}}</span>
	</li>
	{{ end }}
</ol>
{{else}}
	{{if .Query}}
	No results found for "{{.Query}}".
	{{end}}
{{end}}
</section>

{{ end }}
{{ end }}
`
