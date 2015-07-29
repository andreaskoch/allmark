// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

var aliasIndexContentTemplate = `
<li class="shortlink">
	<a href="{{.Route}}" title="► {{.TargetRoute}}">{{.Name}}</a>
</li>
`

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
{{.Content}}
</ol>

</section>
`
