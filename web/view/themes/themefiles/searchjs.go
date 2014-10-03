// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themefiles

const SearchJs = `
// data source: titles
var titlesDataSource = new Bloodhound({
	datumTokenizer: Bloodhound.tokenizers.obj.whitespace('value'),
	queryTokenizer: Bloodhound.tokenizers.whitespace,
	limit: 10,
	prefetch: {
		url: '/titles.json',
	}
});

titlesDataSource.initialize();

// data source: search
var searchDataSource = new Bloodhound({
	datumTokenizer: Bloodhound.tokenizers.obj.whitespace('value'),
	queryTokenizer: Bloodhound.tokenizers.whitespace,
	remote: '/search.json?q=%QUERY'
});

searchDataSource.initialize();

$('.typeahead').typeahead(
	{
		minLength: 1,
		items: 10,
		highlight: true,
	},
	{
		name: 'item-titles',
		displayKey: 'value',
		source: titlesDataSource.ttAdapter(),
		templates: {
			header: '<h3>Items</h3>'
		}
	},
	{
		name: 'searchresults',
		displayKey: 'value',
		source: searchDataSource.ttAdapter(),
		templates: {
			header: '<h3>Search Results</h3>'
		}
	}
).on('typeahead:selected', function(event, datum) {
	window.location = datum.route;
});
`
