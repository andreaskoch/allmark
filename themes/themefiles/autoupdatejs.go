// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themefiles

const AutoupdateJs = `

var Autoupdate = (function () {
    function Autoupdate() { }

    var self = this;

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

                // execute the on change callbacks
                executeOnChangeCallbacks();

                return;
            }

            /**
             * Update an existing item list entry
             * @param Element entryToUpdate The node which shall be updated
             * @param Object model The model containing values to use for the update
             */
            var updateEntry = function(entryToUpdate, model) {
                // update the title text
                $(entryToUpdate).find(".subentry-link:first").html(model.title);

                // update the title link
                $(entryToUpdate).find(".subentry-link:first").attr("href", model.relativeRoute);

                // update the description
                $(entryToUpdate).find(".subentry-description:first").html(model.description);                
            };

            var entries = model.subEntries;
            var existingEntries = $(".subentries>.subentry");
            var numberOfExistingEntries = existingEntries.length;
            var numberOfNewEntries = entries.length;

            var entryTemplate = "<li class=\"subentry\"><a href=\"#\" class=\"subentry-title subentry-link\"></a><p class=\"subentry-description\"></p></li>";

            if (numberOfExistingEntries > numberOfNewEntries) {

                for (var x = (numberOfNewEntries-1); x < numberOfNewEntries; x++) {
                    $(existingEntries[x]).remove();
                }

            }

            // update or add
            for (var i = 0; i < numberOfNewEntries; i++) {
                var index = i + 1;
                var newEntry = entries[i];

                // get the item to update
                if (index <= numberOfExistingEntries) {

                    // update the existing entry
                    updateEntry(existingEntries[i], newEntry);

                } else {

                    // append and update a new entry
                    updateEntry($(".subentries").append(entryTemplate).find(".subentry:last"), newEntry);
                }
            }

            // execute the on change callbacks
            executeOnChangeCallbacks();

        };
    };    

    Autoupdate.prototype.start = function () {

        // check if websockets are supported
        if(!window["WebSocket"]) {
            console.log("Your browser does not support WebSockets.");
            return;
        }

        // establish the connection
        connect(getWebSocketUrl());
    };

    Autoupdate.prototype.onchange = function(name, callback) {
        if (typeof(self.changeCallbacks) !== 'object') {
            self.changeCallbacks = {};
        }

        self.changeCallbacks[name] = callback;
    };

    return Autoupdate;
})();

var autoupdate = new Autoupdate();
autoupdate.start();
`
