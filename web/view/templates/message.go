// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

const messageTemplate = `
<section class="description">
{{.Description}}
</section>

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
