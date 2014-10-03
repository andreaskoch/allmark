// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

import (
	"fmt"
)

var masterTemplate = fmt.Sprintf(`<!DOCTYPE HTML>
<html lang="{{.LanguageTag}}">
<head>
	<meta charset="utf-8">
	<base href="{{ .BaseUrl }}">

	<title>{{.PageTitle}}</title>

	<link rel="schema.DC" href="http://purl.org/dc/terms/">
	<link rel="search" type="application/opensearchdescription+xml" title="{{.RepositoryName}}" href="/opensearch.xml" />

	<meta name="DC.date" content="{{.CreationDate}}">

	{{if .GeoLocation }}
	{{if .GeoLocation.Coordinates}}
	<meta name="geo.position" content="{{.GeoLocation.Coordinates}}">
	{{end}}

	{{if .GeoLocation.PlaceName}}
	<meta name="geo.placename" content="{{.GeoLocation.PlaceName}}">
	{{end}}
	{{end}}

	<link rel="canonical" href="{{.BaseUrl}}">
	<link rel="alternate" type="application/rss+xml" title="RSS" href="/rss.xml">
	<link rel="shortcut icon" href="/theme/favicon.ico">

	<link rel="stylesheet" href="/theme/deck.css" media="screen">
	<link rel="stylesheet" href="/theme/screen.css" media="screen">
	<link rel="stylesheet" href="/theme/print.css" media="print">
	<link rel="stylesheet" href="/theme/codehighlighting/highlight.css" media="screen, print">

	<script src="/theme/modernizr.js"></script>
</head>
<body>

{{ if .ToplevelNavigation}}
<nav class="toplevel">
	<ul>
	{{range .ToplevelNavigation.Entries}}
	<li>
		<a href="{{.Path}}">{{.Title}}</a>
	</li>
	{{end}}
	</ul>
</nav>
{{end}}

<nav class="search">
	<form action="/search" method="GET">
		<input class="typeahead" type="text" name="q" placeholder="search" autocomplete="off">
		<input type="submit" style="visibility:hidden; position: fixed;"/>
	</form>
</nav>

{{ if .BreadcrumbNavigation}}
<nav class="breadcrumb">
	<ul>
	{{range .BreadcrumbNavigation.Entries}}
	<li>
		<a href="{{.Path}}">{{.Title}}</a>
	</li>
	{{end}}
	</ul>
</nav>
{{end}}

<article class="{{.Type}} level-{{.Level}}">
%s
</article>

<aside class="sidebar">

	{{if .ItemNavigation}}
	<section class="navigation">
		{{if .ItemNavigation.Parent}}
		<a href="{{.ItemNavigation.Parent.Path}}" title="{{.ItemNavigation.Parent.Title}}">↑ Parent</a>
		{{end}}
	</section>
	{{end}}

	{{ if .Childs }}
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

</aside>

<div class="cleaner"></div>

{{if or .PrintUrl .JsonUrl .RtfUrl}}
<aside class="export">
<ul>
	{{if .PrintUrl}}<li><a href="{{.PrintUrl}}">Print</a></li>{{end}}
	{{if .JsonUrl}}<li><a href="{{.JsonUrl}}">JSON</a></li>{{end}}
	{{if .RtfUrl}}<li><a href="{{.RtfUrl}}">Rich Text</a></li>{{end}}
</ul>
</aside>
{{end}}

<footer>
	<nav>
		<ul>
			<li><a href="/search">Search</a></li>
			<li><a href="/tags.html">Tags</a></li>
			<li><a href="/sitemap.html">Sitemap</a></li>
			<li><a href="/feed.rss">RSS Feed</a></li>
		</ul>
	</nav>
</footer>

<script src="/theme/jquery.js"></script>
<script src="/theme/jquery.tmpl.js"></script>
<script src="/theme/jquery.lazyload.js"></script>
<script src="/theme/latest.js"></script>
<script src="/theme/typeahead.js"></script>
<script src="/theme/search.js"></script>

<script src="/theme/autoupdate.js"></script>
<script src="/theme/pdf.js"></script>
<script src="/theme/pdf-preview.js"></script>
<script src="/theme/codehighlighting/highlight.js"></script>
<script src="/theme/deck.js"></script>
<script src="/theme/presentation.js"></script>
<script src="/theme/site.js"></script>

</body>
</html>`, ChildTemplatePlaceholder)
