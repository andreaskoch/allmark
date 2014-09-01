// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themefiles

const AutoupdateJs = `

var Autoupdate = (function () {
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
     * Get the Url for the web socket connection
     * @return string The url for the web socket connection (e.g. "ws://example.com:8080/documents/Sample-Document.ws")
     */
    var getWebSocketUrl = function() {
        routeParameter = getCurrentRoute();
        host = document.location.host;
        webSocketHandler = "/ws";

        if (routeParameter === "") {
            return "ws://" + host + "/" + "ws";
        }

        return "ws://" + host + "/" + routeParameter + ".ws";
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

            // stop here if there are no childs
            if (model.childs === null || typeof(model.childs) !== 'object' || model.childs.length == 0) {

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
                $(entryToUpdate).find(".child-link:first").html(model.title);

                // update the title link
                $(entryToUpdate).find(".child-link:first").attr("href", model.route);

                // update the description
                $(entryToUpdate).find(".child-description:first").html(model.description);
            };

            var entries = model.childs;
            var numberOfNewEntries = entries.length;

            // check if the childs container exists
            if ($(".childs").length === 0) {

                // This indicates that the document type has changed.
                // Reload the page if the container does not exists yet
                document.location.reload();
                return;

            }

            var existingEntries = $(".childs>.list>.child");
            var numberOfExistingEntries = existingEntries.length;

            var entryTemplate = "<li class=\"child\"><a href=\"#\" class=\"child-title child-link\"></a><p class=\"child-description\"></p></li>";

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
                    updateEntry($(".childs>.list").append(entryTemplate).find(".child:last"), newEntry);
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
