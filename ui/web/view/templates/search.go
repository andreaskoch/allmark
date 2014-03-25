// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

var searchContentTemplate = `
<nav>
	<form action="/search" method="GET">
		<input type="text" name="q" placeholder="search" value="{{.Query}}">
		<input type="submit" value="Search">
	</form>
</nav>

{{if .ResultCount}}
<header>
	Displaying {{.ResultCount}} of {{.TotalResultCount}} search results for "{{.Query}}":
</header>

<ol>
	{{ range .Results }}
	<li>
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
{{end}}`

const searchTemplate = `
<header>
<h1 class="title">
{{.Title}}
</h1>
</header>

<section class="description">
{{.Description}}
</section>

<section class="content">
{{.Content}}
</section>
`
