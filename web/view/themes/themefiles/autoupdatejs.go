// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themefiles

const AutoupdateJs = `var Autoupdate = (function () {
    function Autoupdate() { }

    var self = this;

    /**
     * Get the currently opened web route
     * @return string The currently opened web route (e.g. "documents/Sample-Document")
     */
    var getCurrentRoute = function() {
        var url = document.location.pathname;

        // remove leading slash
        var leadingSlash = /^\//;
        url = url.replace(leadingSlash, "");

        return url;
    };

    /**
     * Get the URL for the web socket connection
     * @return string The url for the web socket connection (e.g. "ws://example.com:8080/documents/Sample-Document.ws")
     */
    var getWebSocketURL = function() {
        var routeParameter = getCurrentRoute();
        var host = document.location.host;
        var protocol = "ws";
        if (location.protocol === 'https:') {
            protocol = "wss";
        }

        if (routeParameter === "") {
            return protocol + "://" + host + "/" + "ws";
        }

        return protocol + "://" + host + "/" + routeParameter + ".ws";
    };

    /**
     * Execute all on change callbacks
     */
    var executeOnChangeCallbacks = function() {
        if (typeof(self.changeCallbacks) !== 'object') {
            return;
        }

        for (var callbackName in self.changeCallbacks) {
            console.log("Executing on change callback: " + callbackName);
            self.changeCallbacks[callbackName]();
        }
    };

    /**
     * cleanupSnippetCode removes newline characters from the given snippet code
     * @param string snippetCode Code of a snippet that might contain newline characters
     * @return string
     */
    var cleanupSnippetCode = function(snippetCode) {
      return snippetCode.replace('\\n', "");
    };

    /**
     * getCSSSelectorFromSnippetName returns the proper CSS selector for the snippet with the given name.
     * @param string snippetName The name of a snippet (e.g. "toplevelnavigation")
     * @return string A CSS selector that matches the snippet with the given snippetName.
     */
    var getCSSSelectorFromSnippetName = function(snippetName) {
      switch(snippetName) {
        case "aliases":
          return "body > article > section.aliases";

        case "tags":
          return "body > article > section.tags";

        case "publisher":
          return "body > article > section.publisher";

        case "toplevelnavigation":
          return "body>nav.toplevel";

        case "breadcrumbnavigation":
          return "body>nav.breadcrumb";

        case "itemnavigation":
          return "aside.sidebar>nav.navigation";

        case "children":
          return "aside.sidebar>.children";

        case "tagcloud":
          return "aside.sidebar>.tagcloud";

        default:
          return "#" + snippetName;
      }
    };

    /**
     * Connect to the server
     * @param string webSocketURL The url of the web-socket to connect to
     */
    var connect = function(webSocketURL) {
        var reconnectionTimeInSeconds = 3;
        var connection = new WebSocket(webSocketURL);

        connection.onclose = function(evt) {
            console.log("Connection closed. Trying to reconnect in " + reconnectionTimeInSeconds + " seconds.");

            setTimeout(function() {

                console.log("Reconnecting");
                connect(webSocketURL);

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

            // on-change handlers
            executeOnChangeCallbacks();

            // snippets
            if (typeof(model.snippets) !== 'object') {
              return;
            }

            for (var snippetName in model.snippets) {

              // the new snippet content
              var snippetContent = cleanupSnippetCode(model.snippets[snippetName]);

              // get the CSS selector of the snippet
              var cssSelector = getCSSSelectorFromSnippetName(snippetName);

              // check if the snippet exists
              var elementExists = $(cssSelector).length > 0;
              if (elementExists == false) {
                continue;
              }

              // replace the existing snippet
              $(cssSelector).replaceWith(snippetContent);
            }

        };
    };

    Autoupdate.prototype.start = function () {

        // check if websockets are supported
        if(!window["WebSocket"]) {
            console.log("Your browser does not support WebSockets.");
            return;
        }

        // establish the connection
        connect(getWebSocketURL());
    };

    Autoupdate.prototype.onchange = function(name, callback) {
        if (typeof(self.changeCallbacks) !== 'object') {
            self.changeCallbacks = {};
        }

        self.changeCallbacks[name] = callback;
    };

    return Autoupdate;
})();

autoupdate = new Autoupdate();
autoupdate.start();
`
