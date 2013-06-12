// Copyright 2013 Andreas Koch. All rights reserved.
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
	<title>{{.Title}}</title>

	<link rel="schema.DC" href="http://purl.org/dc/terms/">
	<meta name="DC.date" content="{{.Date}}">

	<link rel="stylesheet" type="text/css" href="/theme/screen.css" media="screen">
	<link rel="stylesheet" type="text/css" href="/theme/print.css" media="print">
</head>
<body class="level-{{.Level}}">

<article>
%s
</article>

<script type="text/javascript" src="/theme/jquery.js"></script>
<script type="text/javascript" src="/theme/pdf.js"></script>
<script type="text/javascript" src="/theme/autoupdate.js"></script>
<script type="text/javascript" src="/theme/pdf-preview.js"></script>

</body>
</html>`, ChildTemplatePlaceholder)

const repositoryTemplate = `
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

<section>
<ul class="subentries">
{{range .Childs}}
<li class="subentry">
	<a href="{{.RelativeRoute}}" class="subentry-title subentry-link">{{.Title}}</a>
	<p class="subentry-description">{{.Description}}</p>
</li>
{{end}}
</ul>
</section>
`

const collectionTemplate = `
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

<section class="collection">
<h1>Documents</h2>
<ol class="subentries">
{{range .Childs}}
<li class="subentry">
	<a href="{{.RelativeRoute}}" class="subentry-title subentry-link">{{.Title}}</a>
	<p class="subentry-description">{{.Description}}</p>
</li>
{{end}}
</ol>
</section>
`

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
`

const messageTemplate = `
<section class="content">
{{.Content}}
</section>

<section class="description">
{{.Description}}
</section>
`

const errorTemplate = `
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
`
