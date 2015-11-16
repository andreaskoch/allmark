// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themefiles

const PrintCss = `
* {
    background: #FFFFFF !important;
    color: black !important;
    filter: none !important;
    -ms-filter: none !important;
}

body {
    font: normal 18px/30px "Georgia Pro",georgia,serif;
    text-rendering: optimizeLegibility;
    cursor: default;
    max-width: 100%;
    width: 100%
    outline: none;
}

a, a:visited {
    font-weight: normal;
    text-decoration: none;
}

hr {
    height: 1px;
    border: 0;
    border-bottom: 1px solid black;
}

a[href]:after {
    content: " <" attr(href) ">";
}

a[text^="http"]:after {
    content: "";
}

abbr[title]:after {
    content: " <" attr(title) ">";
}

.ir a:after, a[href^="javascript:"]:after, a[href^="#"]:after {
    content: "";
}

pre, blockquote {
    border: 1px solid #999;
    padding: 1em;
    page-break-inside: avoid;
    white-space: pre-wrap;
}

ol, ul, tr, img {
    page-break-inside: avoid;
}

img {
    max-width: 95% !important;
}

p, h2, h3 {
    orphans: 3;
    widows: 3;
}

h2, h3, h4, h5, h6 {
    font-family: "Franklin ITC Pro Bold",sans-serif !important;
    page-break-after: avoid;
}

a.deeplink {
    display: none;
}

body>nav.toplevel {
    display: none;
}

body>nav.search {
    display: none;
}

body>nav.breadcrumb {
    display: none;
}

body>aside {
    display: none;
}

aside.sidebar {
    display: none;
}

body>footer {
    display: none;
}

.ribbon {
  display: none;
}

.allmark-promo {
  display: none;
}

.video video {
    display: none;
}

.video iframe {
    display: none;
}

.audio audio {
    display: none;
}

.presentation nav {
    display: none;
}`
