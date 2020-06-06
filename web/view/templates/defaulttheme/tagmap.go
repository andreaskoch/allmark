// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package defaulttheme

import (
	"github.com/elWyatt/allmark/web/view/templates/templatenames"
)

func init() {
	templates[templatenames.TagMap] = tagmapTemplate
}

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
{{ if .Tags }}
{{ range .Tags }}
<ul class="tags">
<li class="tag">
	<a name="{{.Anchor}}" href={{.Route}}>{{.Name}}</a>
	{{ if .Children }}
	<ol class="children">
		{{range .Children}}
		<li class="child">
			<a href="{{.Route}}" class="child-title child-link">{{.Title}}</a>
			<p class="child-description">{{.Description}}</p>
		</li>
		{{end}}
	</ol>
	{{ end }}
</li>
</ul>
{{ end }}
{{ else}}
-- There are currently not tagged items --
{{ end }}

</section>
`
