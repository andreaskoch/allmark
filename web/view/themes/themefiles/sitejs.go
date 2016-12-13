// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themefiles

const SiteJs = `
/**
 * appendStyleSheet adds the style sheet with the given path to the page
 * @param {string} path The style-sheet file path
 */
function appendStyleSheet(path) {
	$('<link/>', {
		rel: 'stylesheet',
		type: 'text/css',
		href: path,
	}).appendTo('head');
}

/**
 * getAnchorNameFromText returnes a normalized anchor name for the given text
 * @param {string} text
 * @return {string}
 */
function getAnchorNameFromText(text) {
	var anchorName = text.replace(/[\s]/g, "-");
	anchorName = anchorName.replace(/-{2,}/g, "-");
	return anchorName.replace(/[^\w\d-]/g, "");
}

/**
 * addDeepLinksToElements adds anchor links to the elements with the given selector
 * @param {string} cssSelector A CSS-selector
 */
function addDeepLinksToElements(cssSelector) {
	$(cssSelector).each(function() {
		var headlineText = $(this).text();

		var tagName = this.tagName;
		var headlineLevel = tagName.replace(/[^\d]/g, "");

		var anchorText = headlineLevel + "-" + getAnchorNameFromText(headlineText);

		$(this).before('<a class="deeplink" name="' + anchorText + '">' + headlineText + '</a>');
		$(this).wrap('<a href="#' +anchorText + '"></a>')
	});
}
`
