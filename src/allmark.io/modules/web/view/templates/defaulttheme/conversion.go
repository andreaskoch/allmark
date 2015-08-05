// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package defaulttheme

import (
	"allmark.io/modules/web/view/templates/templatenames"
)

func init() {
	templates[templatenames.Conversion] = converterTemplate
}

const converterTemplate = `
<html>
<head>
	<meta charset="utf-8">
	<meta name="robots" content="noindex,nofollow">
	<link rel="canonical" href="{{ .Route | absolute }}">
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
