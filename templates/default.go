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
			if (typeof(message) !== 'object' || typeof(message.route) !== 'string' || typeof(message.model) !== 'object') {
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
			$('title').text(model.title);
			$('.title').text(model.title);

			// update the description
			$('.description').html(model.description);

			// update the content
			$('.content').html(model.content);

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
<ul class="childs">
{{range .Childs}}
<li>
	<a href="{{.RelativeRoute}}">{{.Title}}</a>
	<p>{{.Description}}</p>
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
<ol class="childs">
{{range .Childs}}
<li>
	<h2><a href="{{.RelativeRoute}}" title="{{.Description}}">{{.Title}}</a></h2>
	<p>{{.Description}}</p>
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
