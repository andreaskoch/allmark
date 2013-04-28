// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

import "fmt"

var masterTemplate = fmt.Sprintf(`<!DOCTYPE HTML>
<html lang="{{.LanguageTag}}">
<meta charset="utf-8">
<head>
	<title>{{.Title}}</title>

	<link rel="schema.DC" href="http://purl.org/dc/terms/">
	<meta name="DC.date" content="{{.Date}}">

	<link rel="stylesheet" type="text/css" href="/theme/screen.css">
</head>
<body>

<article>
%s
</article>

</body>
</html>`, ChildTemplatePlaceholder)
