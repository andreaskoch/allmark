// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themes

const pdfPreviewJs = `
(function() {

    PDFJS.disableWorker = true;

    var createPdfPreview = function() {

        $(".pdf").each(function() {
            
            var pdfDoc = null;

            // file link
            var pdfLinkElement = $(this).find("a:first").first();
            if (pdfLinkElement.length !== 1) {
                // abort: pdf link not found
                return;
            }

            // pdf file path
            var pdfFilePath = $(pdfLinkElement).attr("href");
            if (typeof(pdfFilePath) !== 'string' || pdfFilePath.length == 0) {
                // abort: file path is not set
                return;
            }

            // create meta data elements
            var pagerContainer = $('<div>', { 'class': 'entry pager'}).html('<div class="display"><span class="label">Page:</span> <span class="pager-current-page"></span> <span class="separator">of</span> <span class="pager-total-pagecount"></span></div> <div class="controls"><button title="Show previous page" class="prev">&larr;</button><button title="Show next page" class="next">&rarr;</button></div>');
            var zoomControlContainer = $('<div>', { 'class': 'entry zoomlevel'}).html('<div class="display"><span class="label">Zoom:</span> <span class="zoomlevel-current"></span></div> <div class="controls"><button title="Click to zoom out. Or double-click the document area while holden the <ctl>-key." class="zoom-out">-</button><button title="Click to zoom in. Or double-click the document area." class="zoom-in">+</button></div>');
            var snapshotControlContainer = $('<div>', { 'class': 'entry snapshot'}).html('<a href="" class="open-new-window" title="Click to create a snapshot of the current view">Snapshot</a>');
            var downloadContainer = $('<div>', { 'class': 'entry download'}).html('<a href="' + pdfFilePath +  ' " title="Click to download the pdf: ' + pdfFilePath + '" target="_blank"><span>Download</span></a>');
            
            var metaDataContainer = $('<nav>', { 'class': 'metadata'});
            metaDataContainer.append(pagerContainer);
            metaDataContainer.append(zoomControlContainer);
            metaDataContainer.append(snapshotControlContainer);
            metaDataContainer.append(downloadContainer);

            $(this).append(metaDataContainer);

            // pager: display elements
            var currentPage = $(pagerContainer).find(".display .pager-current-page:first").first();
            var totalPages = $(pagerContainer).find(".display .pager-total-pagecount:first").first();
            if (currentPage.length !== 1 || totalPages.length !== 1) {

                // restore previous state
                $(metaDataContainer).remove();

                // abort: pager elements missing
                return;
            }

            // pager: controls
            var previousPageButton = $(pagerContainer).find(".controls .prev:first").first();
            var nextPageButton = $(pagerContainer).find(".controls .next:first").first();
            if (previousPageButton.length !== 1 || nextPageButton.length !== 1) {

                // restore previous state
                $(metaDataContainer).remove();

                // abort: pager elements missing
                return;
            }

            // zoom: display elements
            var zoomLevel = $(zoomControlContainer).find(".zoomlevel-current:first").first();
            if (zoomLevel.length !== 1) {

                // restore previous state
                $(metaDataContainer).remove();

                // abort: zoom elements missing
                return;
            }

            // pager: controls
            var zoomInButton = $(zoomControlContainer).find(".controls .zoom-in:first").first();
            var zoomOutButton = $(zoomControlContainer).find(".controls .zoom-out:first").first();
            if (zoomInButton.length !== 1 || zoomOutButton.length !== 1) {

                // restore previous state
                $(metaDataContainer).remove();

                // abort: pager elements missing
                return;
            }

            // snapshow: controls
            var snapshotButton = $(snapshotControlContainer).find(".open-new-window:first").first();
            if (snapshotButton.length !== 1) {

                // restore previous state
                $(metaDataContainer).remove();

                // abort: snapshot controls missing
                return;
            }

            // create preview area
            var canvasElement = $('<canvas>');
            var previewContainer = $('<div>', { 'class': 'previewarea' });
            previewContainer.append(canvasElement);

            $(this).append(previewContainer);

            // get the canvas
            var canvas = $(this).find(".previewarea canvas:first").first()[0];
            if (typeof(canvas) !== 'object') {

                // restore previous state
                $(metaDataContainer).remove();
                $(previewContainer).remove();

                // abort: canvas element missing
                return;
            }

            // remove the link
            $(pdfLinkElement).hide();

            /**
             * Render the specified page number
             * @param {integer} pageNumberToDisplay The number of the page to display
             */
            var renderPage = function(pageNumberToDisplay) {

                var scale = getScale();

                // Using promise to fetch the page
                pdfDoc.getPage(pageNumberToDisplay).then(function(page) {

                    // calculate the new viewport according to the given scale
                    var previewContainerWidthAfterScale = $(previewContainer).width() * scale;
                    var currentViewPortWidth = page.getViewport(1.0).width;
                    var scaleRelativeToPreviewContainer = previewContainerWidthAfterScale / currentViewPortWidth;

                    var viewport = page.getViewport(scaleRelativeToPreviewContainer);

                    // adapt canvas size according to new viewport width and heigth
                    canvas.height = viewport.height;
                    canvas.width = viewport.width;

                    // get the canvas context
                    var context = canvas.getContext('2d');

                    // Render PDF page into canvas context
                    var renderContext = {
                        canvasContext: context,
                        viewport: viewport
                    };

                    page.render(renderContext);
                });

                // Update UI Elements
                updateCurrentPageControl(pageNumberToDisplay);
                updateTotalPagesControl(pdfDoc.numPages);
                updateScaleControl(scale);
            };

            /**
             * Get the current zoom level/scale
             * @return {float} The current zoom level / scale (0.0 - n)
             */
            var getScale = function() {
                if (typeof(this.scale) !== "number") {
                    this.scale = 1.00;
                }

                return this.scale;
            };

            /**
             * Set the a new scale
             * @param {float} newScale The new scale level
             */
            var updateScaleControl = function(newScale) {
                var percent = Math.floor(newScale * 100);
                $(zoomLevel).text(percent + "%");
            };

            var zoomIn = function() {
                var currentScale = getScale();
                var newScale = currentScale + 0.25;

                if (newScale <= 2.5) {
                    this.scale = newScale;
                    renderPage(getCurrentPage());
                }
            };

            var zoomOut = function() {
                var currentScale = getScale();
                var newScale = currentScale - 0.25;

                if (newScale >= 0.25) {
                    this.scale = newScale;
                    renderPage(getCurrentPage());
                }
            };

            /**
             * Get the current page number (default: 1)
             * @return {integer} The current page number
             */
            var getCurrentPage = function() {
                var val = $(currentPage).text();
                if (val !== "") {
                    return parseInt(val);
                }

                return 1;
            };

            /**
             * Set the current page numer
             * @param {integer} newPageNumber The new page number
             */
            var updateCurrentPageControl = function(newPageNumber) {
                $(currentPage).text(newPageNumber);
            };

            /**
             * Get the total number of pages in the current pdf
             * @return {integer} The total number of pages in the current pdf document
             */
            var getTotalPageCount = function(pageCount) {
                var total = $(totalPages).text();
                if (total !== "") {
                    return parseInt(total);
                }

                return 1;
            };
            
            /**
             * Set the total number of pages in the current pdf
             * @param {integer} newTotal The new total number of pages
             */
            var updateTotalPagesControl = function(newTotal) {
                $(totalPages).text(newTotal);
            };


            /**
             * Go to the page with the specified page number.
             * @param {integer} pageNumber The page number to display
             */
            var gotoPage = function(pageNumber) {

                var pageNumberToRender = pageNumber;
                var total = getTotalPageCount();

                if (pageNumber > total) {
                    pageNumberToRender = pageNumber % total;
                } else if (pageNumber < 1) {
                    pageNumberToRender = total - pageNumber;
                }

                renderPage(pageNumberToRender);
            };

            /**
             * Go to the next page
             */
            var gotoNextPage = function() {
                return gotoPage(getCurrentPage() + 1);
            };

            /**
             * Go to the previous page
             */
            var gotoPreviousPage = function() {
                return gotoPage(getCurrentPage() - 1);
            };

            var waitForFinalEvent = (function () {

                var timers = {};

                return function (callback, ms, uniqueId) {
                    if (!uniqueId) {
                      uniqueId = "Don't call this twice without a uniqueId";
                    }

                    if (timers[uniqueId]) {
                        clearTimeout (timers[uniqueId]);
                    }

                    timers[uniqueId] = setTimeout(callback, ms);
                };
            })();

            // render on resize
            $(window).resize(function(){
                waitForFinalEvent(function() {
                    renderPage(getCurrentPage());   
                }, 500, "Render after resize");
            });
            
            // goto the the previous page when the previous-page button is clicked
            $(previousPageButton).click(function() {
                gotoPreviousPage();
            });     

            // goto the the next page when the next-page button is clicked
            $(nextPageButton).click(function() {
                gotoNextPage();
            });

            // ctrl-key event listener
            var controlKeyPressed = false;
            (function() {
                /**
                 * A event listener which will set the global controlKeyPressed variable to true if the control key is being pressed
                 * @param {event} evt The key event
                 */
                var listenForControlKey = function(evt) {
                    controlKeyPressed = evt.ctrlKey;            
                };

                // attach the control key listener to the keydown event (-> set controlKeyPressed to 'true')
                $(document).keydown(listenForControlKey);

                // attach the control key listener to the keyup event (-> set controlKeyPressed back to 'false')
                $(document).keyup(listenForControlKey);

            })();

            // zoom-in on double-click
            $(canvas).dblclick(function(event){
                if (controlKeyPressed) {
                    zoomOut();
                } else {
                    zoomIn();
                }
            });

            // attach zoom-in event
            $(zoomInButton).click(function() {
                zoomIn();
            });

            // attach zoom-in event
            $(zoomOutButton).click(function() {
                zoomOut();
            });

            $(snapshotButton).click(function() {
                window.open(canvas.toDataURL('image/png'));
            });

            // Asynchronously download PDF as an ArrayBuffer
            PDFJS.getDocument(pdfFilePath).then(function(doc) {
                pdfDoc = doc;
                renderPage(getCurrentPage());
            });
        });

    }

    // render all pdf documents
    createPdfPreview();

    // register a on change listener
    if (typeof(autoupdate) === 'object' && typeof(autoupdate.onchange) === 'function') {
        autoupdate.onchange(
            "Render PDF Preview",
            function() {
                createPdfPreview();
            }
        );
    }

})();
`
