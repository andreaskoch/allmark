// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package defaulttheme

import (
	"allmark.io/modules/web/view/templates/templatenames"
)

func init() {
	templates[templatenames.AliasIndex] = aliasIndexTemplate
}

const aliasIndexTemplate = `
<header>
<h1 class="title">
{{.Title}}
</h1>
</header>

<section class="description">
{{.Description}}
</section>

<section class="content">

<ol class="shortlinks">

{{ if eq (len .Aliases) 0 }}
-- There are currently not items with aliases in this repository --
{{ else }}
{{ range .Aliases }}
<li class="shortlink">
	<a href="{{.Route}}" title="► {{.TargetRoute}}">{{.Name}}</a>
</li>
{{ end }}
{{ end }}

</ol>

</section>
`
