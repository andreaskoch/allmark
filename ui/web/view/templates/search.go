// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

var searchContentTemplate = `
<nav class="search">
	<form action="search" method="GET">
		<input type="text" name="q" placeholder="search" value="{{.Query}}">
		<input type="submit" value="Search">
	</form>
</nav>

<ol>
	{{ range .Results }}
	<li>
			<a href="{{.Route}}">{{.Title}}</a>
			<p>{{.Description}}</p>
	</li>
	{{ end }}
</ol>`

const searchTemplate = `
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
