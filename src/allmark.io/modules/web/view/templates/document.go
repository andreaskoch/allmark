// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

const documentTemplate = `
<header>
<h1 class="title" itemprop="name">
{{.Title}}
</h1>
</header>

<section class="description" itemprop="description">
{{.Description}}
</section>

{{if or .Author.Name .CreationDate}}
<section class="publisher">
{{if and .Author.Name .Author.URL}}

	created by <span class="author" itemprop="author" rel="author">
	<a href="{{ .Author.URL }}" title="{{ .Author.Name }}" target="_blank">
	{{ .Author.Name }}
	</a>
	</span>

{{else if .Author.Name}}

	created by <span class="author" itemprop="author" rel="author">{{ .Author.Name }}</span>

{{end}}
{{if .CreationDate}}

	{{if not .Author.Name}}created{{end}} on <span class="creationdate" itemprop="dateCreated">{{ .CreationDate }}</span>

{{end}}
</section>
{{end}}

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

{{ if .Locations }}
<div class="cleaner"></div>

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
<div class="cleaner"></div>

<section class="tags">
	<header>
		Tags:
	</header>

	<ul class="tags">
	{{range .Tags}}
	<li class="tag">
		<a href="{{.Route}}" rel="tag">{{.Name}}</a>
	</li>
	{{end}}
	</ul>
</section>
{{end}}
`
