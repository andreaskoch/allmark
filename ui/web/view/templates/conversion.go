// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

const conversionTemplate = `
<html>
<head>
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
