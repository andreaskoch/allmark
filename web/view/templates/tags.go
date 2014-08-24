// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

var tagmapContentTemplate = `
<li class="tag">
	<a name="{{.Name}}" href={{.Route}}>{{.Name}}</a>
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

{{if .TagCloud}}
<section class="tagcloud">
	<div class="tags">
	{{range .TagCloud}}
	<span class="level-{{.Level}}">
		<a href="{{.Route}}">{{.Name}}</a>
	</span>
	{{end}}
	</div>
</section>
{{end}}

<ul class="tags">
{{.Content}}
</ul>

</section>
`
