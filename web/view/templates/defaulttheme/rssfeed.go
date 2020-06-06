// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package defaulttheme

import (
	"github.com/elWyatt/allmark/web/view/templates/templatenames"
)

func init() {
	templates[templatenames.RSSFeed] = rssFeedTemplate
}

var rssFeedTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
<channel>

<title><![CDATA[ {{.Title}} ]]></title>
<description><![CDATA[ {{.Description}} ]]></description>
<link>{{.Link}}</link>
<pubDate>{{.PubDate}}</pubDate>
<ttl>1800</ttl>

{{ range .Items }}
<item>
	<title><![CDATA[ {{.Title}} ]]></title>
	<description><![CDATA[ {{.Description}} ]]></description>
	<link>{{.Link}}</link>
	<pubDate>{{.PubDate}}</pubDate>
</item>
{{ end}}

</channel>
</rss>`
