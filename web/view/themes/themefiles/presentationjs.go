// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themefiles

const PresentationJs = `
$(function() {

  var presentationSelector = 'article.presentation > .content';

  // abort if the presentation selector is not found
  if ($(presentationSelector).length === 0) {
    return;
  }

  /**
   * Split the document body into separate slides
   */
  var transformPresentationStructure = function() {
    var presentationContent = $(presentationSelector).html();
    var slides = presentationContent.split("<hr>")
    var newHtml = '<section class="slide">' + slides.join('</section><section class="slide">') + '</section>';
    $(presentationSelector).html(newHtml);
  };

  var originalWidth = "";
  var originalFontSize = "";


  /**
   * Toggle the page header elements
   */
  var togglePresentationMode = function() {
    $("body>nav.toplevel").toggle();
    $("body>nav.breadcrumb").toggle();
    $("body>nav.search").toggle();
    $("aside.sidebar").toggle();
    $(".publisher").toggle();
    $("article.presentation").toggleClass("presentation-mode");
    $("article.presentation>header").toggle();
    $("article.presentation>nav").toggle();
    $("article.presentation>.description").toggle();
    $("article.presentation>.aliases").toggle();
    $("article.presentation>.tags").toggle();
    $("aside.export").toggle();
    $("body>footer").toggle();
    $(".ribbon").toggle();
    $(".allmark-promo").toggle();

    // toggle width and font size
    if (originalWidth === "" && originalFontSize === "") {
      originalWidth = $("body").css("width");
      originalFontSize = $("body").css("font-size");

      $("body").css("width", "95%");
      $("body").css("font-size", "1.5em");
    } else {
      $("body").css("width", originalWidth);
      $("body").css("font-size", originalFontSize);

      originalWidth = "";
      originalFontSize = "";
    }
  };

  /**
   * Transform all slides into a presentation
   */
  var renderPresentation = function() {


    if ($(presentationSelector).length == 0) {
      // this document is not a presentation
      return;
    }

    // transform the content
    transformPresentationStructure();

    // render the presentation
    $.deck('.slide', {
      selectors: {
        container: presentationSelector
      },

      keys: {
        goto: 71 // 'g'
      }
    });

  };

  // handle keyboard shortcuts
  $(document).keydown(function(e) {

    /* <ctrl> + <shift> */
    if (e.ctrlKey && (e.which === 16) ) {
      console.log( "You pressed Ctrl + Shift" );
      togglePresentationMode();
    }

  });

    // load deck.js
    appendStyleSheet("/theme/deck.css");
    $.getScript("/theme/deck.js", function(){

    // render the presentaton
    renderPresentation();

      // register a on change listener
      if (typeof(autoupdate) === 'object' && typeof(autoupdate.onchange) === 'function') {
          autoupdate.onchange(
              "Render Presentation",
              function() {
                  renderPresentation();
              }
          );
      }

    });


});
`
