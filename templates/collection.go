// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

const collectionTemplate = `
<header>
<h1>
{{.Title}}
</h1>
</header>

<section>
{{.Description}}
</section>

<section>
{{.Content}}
</section>

<section class="collection">
<h1>Documents</h2>
<ol>
{{range .Entries}}
<li>
	<h2><a href="{{.Path}}" title="{{.Description}}">{{.Title}}</a></h2>
	<p>{{.Description}}</p>
</li>
{{end}}
</ol>
</section>
`
