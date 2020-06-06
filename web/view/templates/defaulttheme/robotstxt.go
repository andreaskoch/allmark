// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package defaulttheme

import (
	"github.com/elWyatt/allmark/web/view/templates/templatenames"
)

func init() {
	templates[templatenames.RobotsTxt] = robotsTxtTemplate
}

const robotsTxtTemplate = `{{ range .Disallows }}User-agent: {{.UserAgent}}
{{ range .Paths }}Disallow: {{.}}
{{ end }}{{ end }}
Sitemap: {{.SitemapURL}}
`
