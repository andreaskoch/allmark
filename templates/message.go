// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

const messageTemplate = `<!DOCTYPE HTML>
<html lang="{{.LanguageTag}}">
<head>
	<title>{{.Title}}</title>
</head>

<body>

<article>

<header>
<h1>
{{.Title}}
</h1>
</header>

<section>
{{.Content}}
</section>

</article>

</body>
</html>`