// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

import (
	"fmt"
)

var masterTemplate = fmt.Sprintf(`<!DOCTYPE HTML>
<html lang="{{.LanguageTag}}" itemscope itemtype="http://schema.org/WebPage" prefix="og: http://ogp.me/ns#" prefix="article: http://ogp.me/ns/article#">
<head>
	<base href="{{ .BaseURL }}">

	<title>{{.PageTitle}}</title>
	<meta name="description" content="{{.Description}}">

	<link rel="search" type="application/opensearchdescription+xml" title="{{.RepositoryName}}" href="/opensearch.xml" />

	{{if .Publisher.Name }}
	<meta name="publisher" content="{{.Publisher.Name}}">
	{{end}}

	{{if .GeoLocation }}
	{{if .GeoLocation.Coordinates}}
	<meta name="geo.position" content="{{.GeoLocation.Coordinates}}">
	{{end}}

	{{if .GeoLocation.PlaceName}}
	<meta name="geo.placename" content="{{.GeoLocation.PlaceName}}">
	{{end}}
	{{end}}

	<meta property="og:site_name" content="{{ .RepositoryName }}" />
	<meta property="og:type" content="article" />
	<meta property="og:title" content="{{.PageTitle}}" />
	<meta property="og:description" content="{{.Description}}" />
	<meta property="og:url" content="{{ .Route | absolute }}" />
	{{if .LanguageTag}}<meta property="og:locale" content="{{ replace .LanguageTag "-" "_" }}" />{{end}}
	{{if .Images}}{{range .Images}}
	<meta property="og:image" content="{{ .Route | absolute }}" />{{end}}{{end}}
	{{if .CreationDate}}<meta property="article:published_time" content="{{.CreationDate}}" />{{end}}
	{{if .LastModifiedDate}}<meta property="article:modified_time" content="{{.LastModifiedDate}}" />{{end}}
	{{if .Tags}}{{range .Tags}}
	<meta property="article:tag" content="{{ .Name }}" />{{end}}{{end}}

	<link rel="canonical" href="{{ .Route | absolute }}">
	<link rel="alternate" hreflang="{{.LanguageTag}}" href="{{.Route}}">
	<link rel="alternate" type="application/rss+xml" title="RSS" href="/feed.rss">
	<link rel="shortcut icon" href="/theme/favicon.ico">

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
<nav class="breadcrumb" itemprop="breadcrumb">
	{{range .BreadcrumbNavigation.Entries}}
		<a href="{{.Path}}">{{.Title}}</a>{{if not .IsLast}} » {{end}}
	{{end}}
</nav>
{{end}}

<article class="{{.Type}} level-{{.Level}}" itemprop="mainContentOfPage" itemscope itemtype=http://schema.org/BlogPosting>
%s
</article>

<aside class="sidebar">

	{{if .ItemNavigation}}
	<nav class="navigation">
		<div class="navelement parent">
			{{if .ItemNavigation.Parent}}
			<a href="{{.ItemNavigation.Parent.Path}}" title="{{.ItemNavigation.Parent.Title}}">↑ Parent</a>
			{{end}}
		</div>

		<div class="navelement previous">
			{{if .ItemNavigation.Previous}}
			<a class="previous" href="{{.ItemNavigation.Previous.Path}}" title="{{.ItemNavigation.Previous.Title}}">← Previous</a>
			{{end}}
		</div>

		<div class="navelement next">
			{{if .ItemNavigation.Next}}
			<a class="next" href="{{.ItemNavigation.Next.Path}}" title="{{.ItemNavigation.Next.Title}}">Next →</a>
			{{end}}
		</div>
	</nav>
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

{{if or .PrintURL .JSONURL .MarkdownURL .RTFURL}}
<aside class="export">
<ul>
	{{if .PrintURL}}<li><a href="{{.PrintURL}}">Print</a></li>{{end}}
	{{if .JSONURL}}<li><a href="{{.JSONURL}}">JSON</a></li>{{end}}
	{{if .MarkdownURL}}<li><a href="{{.MarkdownURL}}">Markdown</a></li>{{end}}
	{{if .RTFURL}}<li><a href="{{.RTFURL}}">Rich Text</a></li>{{end}}
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
			<li><a href="/!">Shortlinks</a></li>
		</ul>
	</nav>
</footer>

<script src="/theme/jquery.js"></script>
<script src="/theme/jquery.tmpl.js"></script>
<script src="/theme/lazysizes.js"></script>
<script src="/theme/site.js"></script>
<script src="/theme/typeahead.js"></script>
<script src="/theme/search.js"></script>

{{ if .IsRepositoryItem }}
{{ if .LiveReloadEnabled }}<script src="/theme/autoupdate.js"></script>{{ end }}
<script src="/theme/presentation.js"></script>
<script src="/theme/latest.js"></script>
<script src="/theme/codehighlighting/highlight.js"></script>
<script type="text/javascript">
$(function() {
	$('pre code').each(function(i, block) {
		hljs.highlightBlock(block);
	});

	// register a on change listener
	if (typeof(autoupdate) === 'object' && typeof(autoupdate.onchange) === 'function') {
		autoupdate.onchange(
			"Code Highlighting",
			function() {
				$('pre code').each(function(i, block) {
					hljs.highlightBlock(block);
				});
			}
		);
	}
});
</script>
{{ end }}

{{if .Analytics.Enabled}}
{{if .Analytics.GoogleAnalytics.Enabled}}
<script>
  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
  })(window,document,'script','//www.google-analytics.com/analytics.js','ga');

  ga('create', '{{.Analytics.GoogleAnalytics.TrackingID}}', 'auto');
  ga('send', 'pageview');
</script>
{{end}}
{{end}}

</body>
</html>`, ChildTemplatePlaceholder)
