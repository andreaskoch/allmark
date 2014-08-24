// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

const documentTemplate = `
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

{{ if .Locations }}
<section class="locations">
	<header>
		Locations:
	</header>

	<ol class="list">
	{{range .Locations}}
	<li class="location">
		<a href="{{.Route}}">{{.Title}}</a>
		{{ if .Description }}
		<p>{{.Description}}</p>
		{{end}}

		{{ if .GeoLocation }}

		{{ if .GeoLocation.Address }}
		<p class="address">{{ .GeoLocation.Address }}</p>
		{{end}}

		{{ if .GeoLocation.Coordinates }}
		<p class="geo">
			<span class="latitude">{{ .GeoLocation.Latitude }}</span>;
			<span class="longitude">{{ .GeoLocation.Longitude }}</span>
		</p>
		{{end}}

		{{ end }}
	</li>
	{{end}}
	</ol>
</section>
{{end}}

{{ if .Tags }}
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
