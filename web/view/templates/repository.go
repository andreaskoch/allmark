// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

const repositoryTemplate = `
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

<div class="cleaner"></div>

{{ if .Childs }}
<section class="childs">
<ol class="list">
{{range .Childs}}
<li class="child">
	<a href="{{.Route}}" class="child-title child-link">{{.Title}}</a>
	<p class="child-description">{{.Description}}</p>
</li>
{{end}}
</ol>
</section>
{{end}}

{{if .TagCloud}}
<section class="tagcloud">
	<h1>Tag Cloud</h1>
	<div class="tags">
	{{range .TagCloud}}
	<span class="level-{{.Level}}">
		<a href="{{.Route}}">{{.Name}}</a>
	</span>
	{{end}}
	</div>
</section>
{{end}}
`
