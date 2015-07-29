// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

// aliasesSnippet defines the templates for the aliases section of repository items.
const aliasesSnippet = `{{define "aliases-snippet"}}
{{ if .Aliases }}
<div class="cleaner"></div>
<section class="aliases">

{{ if gt (len .Aliases) 1 }}
	<header title="Direct links to this document">
		Shortlinks:
	</header>
{{else}}
	<header title="A direct link to this document">
		Shortlink:
	</header>
{{end}}

<ul>
{{range .Aliases}}
<li>
	<input type="text" value="{{.Route | absolute}}" title="Redirects to {{.TargetRoute | absolute}}" readonly="readonly" />
</li>
{{end}}
</ul>

</section>
{{end}}
{{end}}`
