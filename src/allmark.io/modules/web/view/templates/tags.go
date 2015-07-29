// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

var tagmapContentTemplate = `
<li class="tag">
	<a name="{{.Anchor}}" href={{.Route}}>{{.Name}}</a>
	<ol class="childs">
		{{range .Childs}}
		<li class="child">
			<a href="{{.Route}}" class="child-title child-link">{{.Title}}</a>
			<p class="child-description">{{.Description}}</p>
		</li>
		{{end}}
	</ol>
</li>
`

const tagmapTemplate = `
<header>
<h1 class="title">
{{.Title}}
</h1>
</header>

<section class="description">
{{.Description}}
</section>

<section class="content">

<ul class="tags">
{{.Content}}
</ul>

</section>
`

const tagsSnippet = `{{define "tags-snippet"}}
{{ if .Tags }}
<div class="cleaner"></div>
<section class="tags">
	<header>
		Tags:
	</header>

	<ul>
	{{range .Tags}}
	<li>
		<a href="{{.Route}}" rel="tag">{{.Name}}</a>
	</li>
	{{end}}
	</ul>
</section>
{{end}}
{{end}}`
