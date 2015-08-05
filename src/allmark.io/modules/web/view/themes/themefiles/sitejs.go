// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themefiles

const SiteJs = `
function appendStyleSheet(path) {
	$('<link/>', {
		rel: 'stylesheet',
		type: 'text/css',
		href: path,
	}).appendTo('head');
}
`
