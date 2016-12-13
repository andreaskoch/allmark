// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themefiles

const LatestJs = `
$(function() {

	/**
	 * Get the currently opened web route
	 * @return string The currently opened web route (e.g. "documents/Sample-Document")
	 */
	var getURL = function() {
	    var url = document.location.pathname;

	    // remove leading slash
	    var leadingSlash = /^\//;
	    url = url.replace(leadingSlash, "");

	    if (url === "") {
	    	return "/latest"
	    }

	    return "/" + url + ".latest";
	};

	var markup = '<li><h1><a href="${route}">${title}</a></h1><p><a href="${route}">${description}</a></p><section>{{html content}}</section></li>';

	$.template( "itemTemplate", markup );

	$.ajax({
		url: getURL(),
		success: function(items) {
			$.each(items, function(index, item) {
				$.tmpl( "itemTemplate", item).appendTo( "article>.preview>ul" );
			});
		},
		dataType: "json"
	});

});
`
