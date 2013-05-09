// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

import (
	"fmt"
)

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

<script src="//ajax.googleapis.com/ajax/libs/jquery/2.0.0/jquery.min.js"></script>
<script type="text/Javascript">
	$(function() {

		// check if websockets are supported
		if (!window["WebSocket"]) {
			console.log("Your browser does not support WebSockets.");
			return;
		}

		/**
		 * Get the currently opened web route
		 * @return string The currently opened web route (e.g. "/documents/Sample Document/index.html")
		 */
		var getCurrentRoute = function() {
			var url = document.location.pathname;
			return url.replace(/^\/+/, "");
		};

		/**
		 * Get the Url for the web socket connection
		 * @return string The url for the web socket connection (e.g. "ws://example.com:8080/ws")
		 */
		var getWebSocketUrl = function() {
			routeParameter = "route=" + getCurrentRoute();
			host = document.location.host;
			webSocketHandler = "/ws";
			websocketUrl = "ws://" + host + webSocketHandler + "?" + routeParameter;
			return websocketUrl;			
		};

		var connection = new WebSocket(getWebSocketUrl());

		connection.onclose = function(evt) {
			console.log("Connection closed.");
		};

		connection.onopen = function() {
			console.log("Connection established.")
		};

		connection.onmessage = function(evt) {

			// validate event data
			if (typeof(evt) !== 'object' || typeof(evt.data) !== 'string') {
				console.log("Invalid data from server.");
				return;
			}

			// unwrap the message
			message = JSON.parse(evt.data);

			// check if all required fields are present
			if (message === null || typeof(message) !== 'object' || typeof(message.route) !== 'string' || message.model === null || typeof(message.model) !== 'object') {
				console.log("Invalid response format.", message);
				return;
			}

			// check if update is applicable for the current route
			if (message.route !== getCurrentRoute()) {
				return;
			}

			// check the model structure
			var model = message.model;
			if (typeof(model.content) !== 'string' || typeof(model.description) !== 'string' || typeof(model.title) !== 'string') {
				console.log("Cannot update the view with the given model object. Missing some required fields.", model);
				return;
			}

			// update the title
			$('title').html(model.title);
			$('.title').html(model.title);

			// update the description
			$('.description').html(model.description);

			// update the content
			$('.content').html(model.content);

			// update sub entries (if available)
			if (model.subEntries === null || typeof(model.subEntries) !== 'object') {
				return;
			}

			var entries = model.subEntries;
			var existingEntries = $(".subentries>.subentry");
			var numberOfExistingEntries = existingEntries.length;
			var numberOfNewEntries = entries.length;

			var fallbackEntryTemplate = "<li><a href=\"#\" class=\"subentry-title subentry-link\"></a><p class=\"subentry-description\"></p></li>";

			for (var i = 0; i < numberOfNewEntries; i++) {
				var index = i + 1;
				var newEntry = entries[i];

				// get the item to update
				var entryToUpdate;
				if (index <= numberOfExistingEntries) {

					// use an existing item
					entryToUpdate = existingEntries[i];

				} else {

					// create a new item
					var entryTemplate = fallbackEntryTemplate;
					if (numberOfExistingEntries > 0) {

						// use the first entry for the template
						entryTemplate = existingEntries[0].html();

					}

					// append the template
					entryToUpdate = $(entryTemplate);
					$(".subentries").append(entryToUpdate);
				}

				// update the title text
				$(entryToUpdate).find(".subentry-link:first").html(newEntry.title);

				// update the title link
				$(entryToUpdate).find(".subentry-link:first").attr("href", newEntry.relativeRoute);

				// update the description
				$(entryToUpdate).find(".subentry-description:first").html(newEntry.description);
			}

		};
	});
</script>

</body>
</html>`, ChildTemplatePlaceholder)

const repositoryTemplate = `
<header>
<h1 class="title">
{{.Title}}
</h1>
</header>

<section class="description">
{{.Description}}
</section>

<section class="content">
{{.Content}}
</section>

<section>
<ul class="subentries">
{{range .Childs}}
<li class="subentry">
	<a href="{{.RelativeRoute}}" class="subentry-title subentry-link">{{.Title}}</a>
	<p class="subentry-description">{{.Description}}</p>
</li>
{{end}}
</ul>
</section>
`

const collectionTemplate = `
<header>
<h1 class="title">
{{.Title}}
</h1>
</header>

<section class="description">
{{.Description}}
</section>

<section class="content">
{{.Content}}
</section>

<section class="collection">
<h1>Documents</h2>
<ol class="subentries">
{{range .Childs}}
<li class="subentry">
	<a href="{{.RelativeRoute}}" class="subentry-title subentry-link">{{.Title}}</a>
	<p class="subentry-description">{{.Description}}</p>
</li>
{{end}}
</ol>
</section>
`

const documentTemplate = `
<header>
<h1 class="title">
{{.Title}}
</h1>
</header>

<section class="description">
{{.Description}}
</section>

<section class="content">
{{.Content}}
</section>
`

const messageTemplate = `
<section class="content">
{{.Content}}
</section>

<section class="description">
{{.Description}}
</section>
`

const errorTemplate = `
<header>
<h1 class="title">
{{.Title}}
</h1>
</header>

<section class="description">
{{.Description}}
</section>

<section class="content">
{{.Content}}
</section>
`
