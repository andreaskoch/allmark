// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themes

var defaultTheme *Theme

func init() {

	const screenCss = `
html {
    font-size: 100%;
    overflow-y: scroll;
    -webkit-text-size-adjust: 100%;
    -ms-text-size-adjust: 100%;
}

body {
    color: #444;
    font-family: Georgia, Palatino, 'Palatino Linotype', Times, 'Times New Roman', "Hiragino Sans GB", "STXihei", "微软雅黑", serif;
    font-size: 12px;
    line-height: 1.5em;
    background: #fefefe;
    width: 75%;
    margin: 10px auto;
    padding: 1em;
    outline: 1300px solid #FAFAFA;
}

a {
    color: #0645ad;
    text-decoration: none;
}

a:visited {
    color: #0b0080;
}

a:hover {
    color: #06e;
}

a:active {
    color: #faa700;
}

a:focus {
    outline: thin dotted;
}

a:hover, a:active {
    outline: 0;
}

span.backtick {
    border: 1px solid #EAEAEA;
    border-radius: 3px;
    background: #F8F8F8;
    padding: 0 3px 0 3px;
}

::-moz-selection {
    background: rgba(255,255,0,0.3);
    color: #000;
}

::selection {
    background: rgba(255,255,0,0.3);
    color: #000;
}

a::-moz-selection {
    background: rgba(255,255,0,0.3);
    color: #0645ad;
}

a::selection {
    background: rgba(255,255,0,0.3);
    color: #0645ad;
}

p {
    margin: 1em 0;
}

img {
    max-width: 100%;
}

h1,h2,h3,h4,h5,h6 {
    font-weight: normal;
    color: #111;
    line-height: 1em;
}

h4,h5,h6 {
    font-weight: bold;
}

h1 {
    font-size: 2.5em;
}

h2 {
    font-size: 2em;
    border-bottom: 1px solid silver;
    padding-bottom: 5px;
}

h3 {
    font-size: 1.5em;
}

h4 {
    font-size: 1.2em;
}

h5 {
    font-size: 1em;
}

h6 {
    font-size: 0.9em;
}

blockquote {
    color: #666666;
    margin: 0;
    padding-left: 3em;
    border-left: 0.5em #EEE solid;
}

hr {
    display: block;
    height: 2px;
    border: 0;
    border-top: 1px solid #aaa;
    border-bottom: 1px solid #eee;
    margin: 1em 0;
    padding: 0;
}

pre , code, kbd, samp {
    color: #000;
    font-family: monospace;
    font-size: 0.88em;
    border-radius: 3px;
    background-color: #F8F8F8;
    border: 1px solid #CCC;
}

pre {
    white-space: pre;
    white-space: pre-wrap;
    word-wrap: break-word;
    padding: 5px 12px;
}

pre code {
    border: 0px !important;
    padding: 0;
}

code {
    padding: 0 3px 0 3px;
}

b, strong {
    font-weight: bold;
}

dfn {
    font-style: italic;
}

ins {
    background: #ff9;
    color: #000;
    text-decoration: none;
}

mark {
    background: #ff0;
    color: #000;
    font-style: italic;
    font-weight: bold;
}

sub, sup {
    font-size: 75%;
    line-height: 0;
    position: relative;
    vertical-align: baseline;
}

sup {
    top: -0.5em;
}

sub {
    bottom: -0.25em;
}

ul, ol {
    margin: 1em 0;
    padding: 0 0 0 2em;
}

li p:last-child {
    margin: 0;
}

dd {
    margin: 0 0 0 2em;
}

img {
    border: 0;
    -ms-interpolation-mode: bicubic;
    vertical-align: middle;
}

table {
    border-collapse: collapse;
    border-spacing: 0;
}

td {
    vertical-align: top;
}

.imagegallery>h1 {
    font-size: 1.2em;
}

.imagegallery ol {
    list-style: none;
    margin-left: 0;
}

.collection>h1 {
    font-size: 1.2em;
}

.collection>ol {
    list-style: none;
    padding: 0;
    margin: 0 0 0 15px;
}

.collection>ol>li {
    display: block;
    margin: 0 0 15px 0;
}

.collection>ol>li>h2 {
    font-size: 1.0em;
    font-weight: bold;
    margin: 0;
}

.csv
{
    font-family: "Lucida Sans Unicode", "Lucida Grande", Sans-Serif;
    font-size: 1.0em;
    margin: 45px;
    text-align: left;
    border-collapse: collapse;
    border: 1px solid #69c;
}
.csv thead
{
    padding: 12px 17px 12px 17px;
    font-weight: normal;
    font-size: 1.2em;
    color: #039;
    border-bottom: 1px dashed #69c;
}
.csv td
{
    padding: 7px 17px 7px 17px;
    color: #669;
}
.csv tbody tr:hover td
{
    color: #339;
    background: #d0dafd;
}

@media only screen and (min-width: 480px) {
    body {
        font-size: 14px;
        width: 95%
    };
}

@media only screen and (min-width: 768px) {
    body {
        font-size: 16px;
        width: 95%
    };
}

@media only screen and (min-width: 1024px) {
    body {
        font-size: 16px;
        width: 75%
    };
}

@media print {
    * {
        background: transparent !important;
        color: black !important;
        filter: none !important;
        -ms-filter: none !important;
    }

    body {
        font-size: 12pt;
        max-width: 100%;
        width: 100%
        outline: none;
    }

    a, a:visited {
        text-decoration: underline;
    }

    hr {
        height: 1px;
        border: 0;
        border-bottom: 1px solid black;
    }

    a[href]:after {
        content: " (" attr(href) ")";
    }

    abbr[title]:after {
        content: " (" attr(title) ")";
    }

    .ir a:after, a[href^="javascript:"]:after, a[href^="#"]:after {
        content: "";
    }

    pre, blockquote {
        border: 1px solid #999;
        padding-right: 1em;
        page-break-inside: avoid;
    }

    tr, img {
        page-break-inside: avoid;
    }

    img {
        max-width: 95% !important;
    }

    @page :left {
        margin: 15mm 20mm 15mm 10mm;
    }

    @page :right {
        margin: 15mm 10mm 15mm 20mm;
    }

    p, h2, h3 {
        orphans: 3;
        widows: 3;
    }

    h2, h3 {
        page-break-after: avoid;
    };
}`

	const autoupdateJs = `
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
        return decodeURI(url.replace(/^\/+/, ""));
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

    /**
     * Connect to the server
     * @param string webSocketUrl The url of the web-socket to connect to
     */
    var connect = function(webSocketUrl) {
        var reconnectionTimeInSeconds = 3;
        var connection = new WebSocket(webSocketUrl);

        connection.onclose = function(evt) {
            console.log("Connection closed. Trying to reconnect in " + reconnectionTimeInSeconds + " seconds.");

            setTimeout(function() {

                console.log("Reconnecting");
                connect(webSocketUrl);

            }, (reconnectionTimeInSeconds * 1000));
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
                console.log("no match", message);
                return;
            }

            console.log("match", message);

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
    };

    // establish the connection
    connect(getWebSocketUrl());
});`

	defaultTheme = &Theme{
		Name: "default",
		Files: []*ThemeFile{
			&ThemeFile{
				Filename: "screen.css",
				Content:  screenCss,
			},
			&ThemeFile{
				Filename: "autoupdate.js",
				Content:  autoupdateJs,
			},
		},
	}
}
