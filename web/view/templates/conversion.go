// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

const converterTemplate = `
<html>
<head>
	<meta charset="utf-8">
	<meta name="robots" content="noindex,nofollow">

	<title>{{.Title}}</title>

	<link rel="canonical" href="{{.BaseUrl}}">
	<link rel="stylesheet" href="/theme/print.css">
</head>
<body>
<h1>
{{.Title}}
</h1>

<p>
{{.Description}}
</p>

{{.Content}}
</body>
</html>
`
