// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

const messageTemplate = `
<section class="description">
{{.Description}}
</section>

<section class="content">
{{.Content}}
</section>

<div class="cleaner"></div>

{{ if .Childs }}
<section class="preview">
	<ul>
	</ul>
</section>

<section class="childs">
<h1>Childs</h1>

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

{{ if .Tags }}
<div class="cleaner"></div>

<section class="tags">
	<header>
		Tags:
	</header>

	<ul class="tags">
	{{range .Tags}}
	<li class="tag">
		<a href="{{.Route}}">{{.Name}}</a>
	</li>
	{{end}}
	</ul>
</section>
{{end}}
`
