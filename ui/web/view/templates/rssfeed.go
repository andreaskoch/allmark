// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

import (
	"fmt"
)

var rssFeedTemplate = fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
<channel>

<title><![CDATA[ {{.Title}} ]]></title>
<description><![CDATA[ {{.Description}} ]]></description>
<link>{{.Link}}</link>
<pubDate>{{.PubDate}}</pubDate>
<ttl>1800</ttl>

%s

</channel>
</rss>`, ChildTemplatePlaceholder)

var rssFeedContentTemplate = `<item>
	<title><![CDATA[ {{.Title}} ]]></title>
	<description><![CDATA[ {{.Description}} ]]></description>
	<link>{{.Link}}</link>
	<pubDate>{{.PubDate}}</pubDate>
</item>`
