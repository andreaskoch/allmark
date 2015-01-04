// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

const presentationTemplate = `
<header>
<h1 class="title" itemprop="name">
{{.Title}}
</h1>
</header>

<section class="description" itemprop="description">
{{.Description}}
</section>

<section class="publisher">
created {{if .Author}}by <span class="author" itemprop="author" rel="author"><a href="{{ .Author.Url }}" title="{{ .Author.Name }}" target="_blank">{{ .Author.Name }}</a></span>{{end}}{{ if .CreationDate }} on <span class="creationdate" itemprop="dateCreated">{{ .CreationDate }}</span>{{ end }}
</section>

<nav>
	<div class="nav-element pager deck-status">
		<span class="deck-status-current"></span> /	<span class="deck-status-total"></span>
	</div>

	<div class="nav-element controls">
		<button class="deck-prev-link" title="Previous">&#8592;</button>
		<button href="#" class="deck-next-link" title="Next">&#8594;</button>
	</div>

	<div class="nav-element jumper">
		<form action="." method="get" class="goto-form">
			<label for="goto-slide">Go to slide:</label>
			<input type="text" name="slidenum" id="goto-slide" list="goto-datalist">
			<datalist id="goto-datalist"></datalist>
			<input type="submit" value="Go">
		</form>
	</div>
</nav>

<section class="content" itemprop="articleBody">
{{.Content}}
</section>

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
