// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themefiles

const LatestJs = `
/*!
Deck JS - deck.core
Copyright (c) 2011 Caleb Troughton
Dual licensed under the MIT license and GPL license.
https://github.com/imakewebthings/deck.js/blob/master/MIT-license.txt
https://github.com/imakewebthings/deck.js/blob/master/GPL-license.txt
*/

/*
The deck.core module provides all the basic functionality for creating and
moving through a deck.  It does so by applying classes to indicate the state of
the deck and its slides, allowing CSS to take care of the visual representation
of each state.  It also provides methods for navigating the deck and inspecting
its state, as well as basic key bindings for going to the next and previous
slides.  More functionality is provided by wholly separate extension modules
that use the API provided by core.
*/
(function($, deck, document, undefined) {
	var slides, // Array of all the uh, slides...
	current, // Array index of the current slide
	$container, // Keeping this cached
	
	events = {
		/*
		This event fires whenever the current slide changes, whether by way of
		next, prev, or go. The callback function is passed two parameters, from
		and to, equal to the indices of the old slide and the new slide
		respectively. If preventDefault is called on the event within this handler
		the slide change does not occur.
		
		$(document).bind('deck.change', function(event, from, to) {
		   alert('Moving from slide ' + from + ' to ' + to);
		});
		*/
		change: 'deck.change',
		
		/*
		This event fires at the beginning of deck initialization, after the options
		are set but before the slides array is created.  This event makes a good hook
		for preprocessing extensions looking to modify the deck.
		*/
		beforeInitialize: 'deck.beforeInit',
		
		/*
		This event fires at the end of deck initialization. Extensions should
		implement any code that relies on user extensible options (key bindings,
		element selectors, classes) within a handler for this event. Native
		events associated with Deck JS should be scoped under a .deck event
		namespace, as with the example below:
		
		var $d = $(document);
		$.deck.defaults.keys.myExtensionKeycode = 70; // 'h'
		$d.bind('deck.init', function() {
		   $d.bind('keydown.deck', function(event) {
		      if (event.which === $.deck.getOptions().keys.myExtensionKeycode) {
		         // Rock out
		      }
		   });
		});
		*/
		initialize: 'deck.init' 
	},
	
	options = {},
	$d = $(document),
	
	/*
	Internal function. Updates slide and container classes based on which
	slide is the current slide.
	*/
	updateStates = function() {
		var oc = options.classes,
		osc = options.selectors.container,
		old = $container.data('onSlide'),
		$all = $();
		
		// Container state
		$container.removeClass(oc.onPrefix + old)
			.addClass(oc.onPrefix + current)
			.data('onSlide', current);
		
		// Remove and re-add child-current classes for nesting
		$('.' + oc.current).parentsUntil(osc).removeClass(oc.childCurrent);
		slides[current].parentsUntil(osc).addClass(oc.childCurrent);
		
		// Remove previous states
		$.each(slides, function(i, el) {
			$all = $all.add(el);
		});
		$all.removeClass([
			oc.before,
			oc.previous,
			oc.current,
			oc.next,
			oc.after
		].join(" "));
		
		// Add new states back in
		slides[current].addClass(oc.current);
		if (current > 0) {
			slides[current-1].addClass(oc.previous);
		}
		if (current + 1 < slides.length) {
			slides[current+1].addClass(oc.next);
		}
		if (current > 1) {
			$.each(slides.slice(0, current - 1), function(i, el) {
				el.addClass(oc.before);
			});
		}
		if (current + 2 < slides.length) {
			$.each(slides.slice(current+2), function(i, el) {
				el.addClass(oc.after);
			});
		}
	},
	
	/* Methods exposed in the jQuery.deck namespace */
	methods = {
		
		/*
		jQuery.deck(selector, options)
		
		selector: string | jQuery | array
		options: object, optional
				
		Initializes the deck, using each element matched by selector as a slide.
		May also be passed an array of string selectors or jQuery objects, in
		which case each selector in the array is considered a slide. The second
		parameter is an optional options object which will extend the default
		values.
		
		$.deck('.slide');
		
		or
		
		$.deck([
		   '#first-slide',
		   '#second-slide',
		   '#etc'
		]);
		*/	
		init: function(elements, opts) {
			var startTouch,
			tolerance,
			esp = function(e) {
				e.stopPropagation();
			};
			
			options = $.extend(true, {}, $[deck].defaults, opts);
			slides = [];
			current = 0;
			$container = $(options.selectors.container);
			tolerance = options.touch.swipeTolerance;
			
			// Pre init event for preprocessing hooks
			$d.trigger(events.beforeInitialize);
			
			// Hide the deck while states are being applied to kill transitions
			$container.addClass(options.classes.loading);
			
			// Fill slides array depending on parameter type
			if ($.isArray(elements)) {
				$.each(elements, function(i, e) {
					slides.push($(e));
				});
			}
			else {
				$(elements).each(function(i, e) {
					slides.push($(e));
				});
			}
			
			/* Remove any previous bindings, and rebind key events */
			$d.unbind('keydown.deck').bind('keydown.deck', function(e) {
				if (e.which === options.keys.next || $.inArray(e.which, options.keys.next) > -1) {
					methods.next();
					e.preventDefault();
				}
				else if (e.which === options.keys.previous || $.inArray(e.which, options.keys.previous) > -1) {
					methods.prev();
					e.preventDefault();
				}
			})
			/* Stop propagation of key events within editable elements */
			.undelegate('input, textarea, select, button, meter, progress, [contentEditable]', 'keydown', esp)
			.delegate('input, textarea, select, button, meter, progress, [contentEditable]', 'keydown', esp);
			
			/* Bind touch events for swiping between slides on touch devices */
			$container.unbind('touchstart.deck').bind('touchstart.deck', function(e) {
				if (!startTouch) {
					startTouch = $.extend({}, e.originalEvent.targetTouches[0]);
				}
			})
			.unbind('touchmove.deck').bind('touchmove.deck', function(e) {
				$.each(e.originalEvent.changedTouches, function(i, t) {
					if (startTouch && t.identifier === startTouch.identifier) {
						if (t.screenX - startTouch.screenX > tolerance || t.screenY - startTouch.screenY > tolerance) {
							$[deck]('prev');
							startTouch = undefined;
						}
						else if (t.screenX - startTouch.screenX < -1 * tolerance || t.screenY - startTouch.screenY < -1 * tolerance) {
							$[deck]('next');
							startTouch = undefined;
						}
						return false;
					}
				});
				e.preventDefault();
			})
			.unbind('touchend.deck').bind('touchend.deck', function(t) {
				$.each(t.originalEvent.changedTouches, function(i, t) {
					if (startTouch && t.identifier === startTouch.identifier) {
						startTouch = undefined;
					}
				});
			})
			.scrollLeft(0).scrollTop(0);
			
			/*
			Kick iframe videos, which dont like to redraw w/ transforms.
			Remove this if Webkit ever fixes it.
			 */
			$.each(slides, function(i, $el) {
				$el.unbind('webkitTransitionEnd.deck').bind('webkitTransitionEnd.deck',
				function(event) {
					if ($el.hasClass($[deck]('getOptions').classes.current)) {
						var embeds = $(this).find('iframe').css('opacity', 0);
						window.setTimeout(function() {
							embeds.css('opacity', 1);
						}, 100);
					}
				});
			});
			
			if (slides.length) {
				updateStates();
			}
			
			// Show deck again now that slides are in place
			$container.removeClass(options.classes.loading);
			$d.trigger(events.initialize);
		},
		
		/*
		jQuery.deck('go', index)
		
		index: integer | string
		
		Moves to the slide at the specified index if index is a number. Index is
		0-based, so $.deck('go', 0); will move to the first slide. If index is a
		string this will move to the slide with the specified id. If index is out
		of bounds or doesn't match a slide id the call is ignored.
		*/
		go: function(index) {
			var e = $.Event(events.change),
			ndx;
			
			/* Number index, easy. */
			if (typeof index === 'number' && index >= 0 && index < slides.length) {
				ndx = index;
			}
			/* Id string index, search for it and set integer index */
			else if (typeof index === 'string') {
				$.each(slides, function(i, $slide) {
					if ($slide.attr('id') === index) {
						ndx = i;
						return false;
					}
				});
			};
			
			/* Out of bounds, id doesn't exist, illegal input, eject */
			if (typeof ndx === 'undefined') return;
			
			$d.trigger(e, [current, ndx]);
			if (e.isDefaultPrevented()) {
				/* Trigger the event again and undo the damage done by extensions. */
				$d.trigger(events.change, [ndx, current]);
			}
			else {
				current = ndx;
				updateStates();
			}
		},
		
		/*
		jQuery.deck('next')
		
		Moves to the next slide. If the last slide is already active, the call
		is ignored.
		*/
		next: function() {
			methods.go(current+1);
		},
		
		/*
		jQuery.deck('prev')
		
		Moves to the previous slide. If the first slide is already active, the
		call is ignored.
		*/
		prev: function() {
			methods.go(current-1);
		},
		
		/*
		jQuery.deck('getSlide', index)
		
		index: integer, optional
		
		Returns a jQuery object containing the slide at index. If index is not
		specified, the current slide is returned.
		*/
		getSlide: function(index) {
			var i = typeof index !== 'undefined' ? index : current;
			if (typeof i != 'number' || i < 0 || i >= slides.length) return null;
			return slides[i];
		},
		
		/*
		jQuery.deck('getSlides')
		
		Returns all slides as an array of jQuery objects.
		*/
		getSlides: function() {
			return slides;
		},
		
		/*
		jQuery.deck('getContainer')
		
		Returns a jQuery object containing the deck container as defined by the
		container option.
		*/
		getContainer: function() {
			return $container;
		},
		
		/*
		jQuery.deck('getOptions')
		
		Returns the options object for the deck, including any overrides that
		were defined at initialization.
		*/
		getOptions: function() {
			return options;
		},
		
		/*
		jQuery.deck('extend', name, method)
		
		name: string
		method: function
		
		Adds method to the deck namespace with the key of name. This doesn’t
		give access to any private member data — public methods must still be
		used within method — but lets extension authors piggyback on the deck
		namespace rather than pollute jQuery.
		
		$.deck('extend', 'alert', function(msg) {
		   alert(msg);
		});

		// Alerts 'boom'
		$.deck('alert', 'boom');
		*/
		extend: function(name, method) {
			methods[name] = method;
		}
	};
	
	/* jQuery extension */
	$[deck] = function(method, arg) {
		if (methods[method]) {
			return methods[method].apply(this, Array.prototype.slice.call(arguments, 1));
		}
		else {
			return methods.init(method, arg);
		}
	};
	
	/*
	The default settings object for a deck. All deck extensions should extend
	this object to add defaults for any of their options.
	
	options.classes.after
		This class is added to all slides that appear after the 'next' slide.
	
	options.classes.before
		This class is added to all slides that appear before the 'previous'
		slide.
		
	options.classes.childCurrent
		This class is added to all elements in the DOM tree between the
		'current' slide and the deck container. For standard slides, this is
		mostly seen and used for nested slides.
		
	options.classes.current
		This class is added to the current slide.
		
	options.classes.loading
		This class is applied to the deck container during loading phases and is
		primarily used as a way to short circuit transitions between states
		where such transitions are distracting or unwanted.  For example, this
		class is applied during deck initialization and then removed to prevent
		all the slides from appearing stacked and transitioning into place
		on load.
		
	options.classes.next
		This class is added to the slide immediately following the 'current'
		slide.
		
	options.classes.onPrefix
		This prefix, concatenated with the current slide index, is added to the
		deck container as you change slides.
		
	options.classes.previous
		This class is added to the slide immediately preceding the 'current'
		slide.
		
	options.selectors.container
		Elements matched by this CSS selector will be considered the deck
		container. The deck container is used to scope certain states of the
		deck, as with the onPrefix option, or with extensions such as deck.goto
		and deck.menu.
		
	options.keys.next
		The numeric keycode used to go to the next slide.
		
	options.keys.previous
		The numeric keycode used to go to the previous slide.
		
	options.touch.swipeTolerance
		The number of pixels the users finger must travel to produce a swipe
		gesture.
	*/
	$[deck].defaults = {
		classes: {
			after: 'deck-after',
			before: 'deck-before',
			childCurrent: 'deck-child-current',
			current: 'deck-current',
			loading: 'deck-loading',
			next: 'deck-next',
			onPrefix: 'on-slide-',
			previous: 'deck-previous'
		},
		
		selectors: {
			container: '.deck-container'
		},
		
		keys: {
			// enter, space, page down, right arrow, down arrow,
			next: [13, 32, 34, 39, 40],
			// backspace, page up, left arrow, up arrow
			previous: [8, 33, 37, 38]
		},
		
		touch: {
			swipeTolerance: 60
		}
	};
	
	$d.ready(function() {
		$('html').addClass('ready');
	});
	
	/*
	FF + Transforms + Flash video don't get along...
	Firefox will reload and start playing certain videos after a
	transform.  Blanking the src when a previously shown slide goes out
	of view prevents this.
	*/
	$d.bind('deck.change', function(e, from, to) {
		var oldFrames = $[deck]('getSlide', from).find('iframe'),
		newFrames = $[deck]('getSlide', to).find('iframe');
		
		oldFrames.each(function() {
	    	var $this = $(this),
	    	curSrc = $this.attr('src');
            
            if(curSrc) {
            	$this.data('deck-src', curSrc).attr('src', '');
            }
		});
		
		newFrames.each(function() {
			var $this = $(this),
			originalSrc = $this.data('deck-src');
			
			if (originalSrc) {
				$this.attr('src', originalSrc);
			}
		});
	});
})(jQuery, 'deck', document);

/*!
Deck JS - deck.goto
Copyright (c) 2011 Caleb Troughton
Dual licensed under the MIT license and GPL license.
https://github.com/imakewebthings/deck.js/blob/master/MIT-license.txt
https://github.com/imakewebthings/deck.js/blob/master/GPL-license.txt
*/

/*
This module adds the necessary methods and key bindings to show and hide a form
for jumping to any slide number/id in the deck (and processes that form
accordingly). The form-showing state is indicated by the presence of a class on
the deck container.
*/
(function($, deck, undefined) {
	var $d = $(document);
	
	/*
	Extends defaults/options.
	
	options.classes.goto
		This class is added to the deck container when showing the Go To Slide
		form.
		
	options.selectors.gotoDatalist
		The element that matches this selector is the datalist element that will
		be populated with options for each of the slide ids.  In browsers that
		support the datalist element, this provides a drop list of slide ids to
		aid the user in selecting a slide.
		
	options.selectors.gotoForm
		The element that matches this selector is the form that is submitted
		when a user hits enter after typing a slide number/id in the gotoInput
		element.
	
	options.selectors.gotoInput
		The element that matches this selector is the text input field for
		entering a slide number/id in the Go To Slide form.
		
	options.keys.goto
		The numeric keycode used to show the Go To Slide form.
		
	options.countNested
		If false, only top level slides will be counted when entering a
		slide number.
	*/
	$.extend(true, $[deck].defaults, {
		classes: {
			goto: 'deck-goto'
		},
		
		selectors: {
			gotoDatalist: '#goto-datalist',
			gotoForm: '.goto-form',
			gotoInput: '#goto-slide'
		},
		
		keys: {
			goto: 71 // g
		},
		
		countNested: true
	});

	/*
	jQuery.deck('showGoTo')
	
	Shows the Go To Slide form by adding the class specified by the goto class
	option to the deck container.
	*/
	$[deck]('extend', 'showGoTo', function() {
		$[deck]('getContainer').addClass($[deck]('getOptions').classes.goto);
		$($[deck]('getOptions').selectors.gotoInput).focus();
	});

	/*
	jQuery.deck('hideGoTo')
	
	Hides the Go To Slide form by removing the class specified by the goto class
	option from the deck container.
	*/
	$[deck]('extend', 'hideGoTo', function() {
		$($[deck]('getOptions').selectors.gotoInput).blur();
		$[deck]('getContainer').removeClass($[deck]('getOptions').classes.goto);
	});

	/*
	jQuery.deck('toggleGoTo')
	
	Toggles between showing and hiding the Go To Slide form.
	*/
	$[deck]('extend', 'toggleGoTo', function() {
		$[deck]($[deck]('getContainer').hasClass($[deck]('getOptions').classes.goto) ? 'hideGoTo' : 'showGoTo');
	});
	
	$d.bind('deck.init', function() {
		var opts = $[deck]('getOptions'),
		$datalist = $(opts.selectors.gotoDatalist),
		slideTest = $.map([
			opts.classes.before,
			opts.classes.previous,
			opts.classes.current,
			opts.classes.next,
			opts.classes.after
		], function(el, i) {
			return '.' + el;
		}).join(', '),
		rootCounter = 1;
		
		// Bind key events
		$d.unbind('keydown.deckgoto').bind('keydown.deckgoto', function(e) {
			var key = $[deck]('getOptions').keys.goto;
			
			if (e.which === key || $.inArray(e.which, key) > -1) {
				e.preventDefault();
				$[deck]('toggleGoTo');
			}
		});
		
		/* Populate datalist and work out countNested*/
		$.each($[deck]('getSlides'), function(i, $slide) {
			var id = $slide.attr('id'),
			$parentSlides = $slide.parentsUntil(opts.selectors.container, slideTest);
			
			if (id) {
				$datalist.append('<option value="' + id + '">');
			}
			
			if ($parentSlides.length) {
				$slide.removeData('rootIndex');
			}
			else if (!opts.countNested) {
				$slide.data('rootIndex', rootCounter);
				++rootCounter;
			}
		});
		
		// Process form submittal, go to the slide entered
		$(opts.selectors.gotoForm)
		.unbind('submit.deckgoto')
		.bind('submit.deckgoto', function(e) {
			var $field = $($[deck]('getOptions').selectors.gotoInput),
			ndx = parseInt($field.val(), 10);
			
			if (!$[deck]('getOptions').countNested) {
			  if (ndx >= rootCounter) return false;
				$.each($[deck]('getSlides'), function(i, $slide) {
					if ($slide.data('rootIndex') === ndx) {
						ndx = i + 1;
						return false;
					}
				});
			}
			
			$[deck]('go', isNaN(ndx) ? $field.val() : ndx - 1);
			$[deck]('hideGoTo');
			$field.val('');
			
			e.preventDefault();
		});
		
		// Dont let keys in the input trigger deck actions
		$(opts.selectors.gotoInput)
		.unbind('keydown.deckgoto')
		.bind('keydown.deckgoto', function(e) {
			e.stopPropagation();
		});
	});
})(jQuery, 'deck');

/*!
Deck JS - deck.hash
Copyright (c) 2011 Caleb Troughton
Dual licensed under the MIT license and GPL license.
https://github.com/imakewebthings/deck.js/blob/master/MIT-license.txt
https://github.com/imakewebthings/deck.js/blob/master/GPL-license.txt
*/

/*
This module adds deep linking to individual slides, enables internal links
to slides within decks, and updates the address bar with the hash as the user
moves through the deck. A permalink anchor is also updated. Standard themes
hide this link in browsers that support the History API, and show it for
those that do not. Slides that do not have an id are assigned one according to
the hashPrefix option. In addition to the on-slide container state class
kept by core, this module adds an on-slide state class that uses the id of each
slide.
*/
(function ($, deck, window, undefined) {
	var $d = $(document),
	$window = $(window),
	
	/* Collection of internal fragment links in the deck */
	$internals,
	
	/*
	Internal only function.  Given a string, extracts the id from the hash,
	matches it to the appropriate slide, and navigates there.
	*/
	goByHash = function(str) {
		var id = str.substr(str.indexOf("#") + 1),
		slides = $[deck]('getSlides');
		
		$.each(slides, function(i, $el) {
			if ($el.attr('id') === id) {
				$[deck]('go', i);
				return false;
			}
		});
		
		// If we don't set these to 0 the container scrolls due to hashchange
		$[deck]('getContainer').scrollLeft(0).scrollTop(0);
	};
	
	/*
	Extends defaults/options.
	
	options.selectors.hashLink
		The element matching this selector has its href attribute updated to
		the hash of the current slide as the user navigates through the deck.
		
	options.hashPrefix
		Every slide that does not have an id is assigned one at initialization.
		Assigned ids take the form of hashPrefix + slideIndex, e.g., slide-0,
		slide-12, etc.

	options.preventFragmentScroll
		When deep linking to a hash of a nested slide, this scrolls the deck
		container to the top, undoing the natural browser behavior of scrolling
		to the document fragment on load.
	*/
	$.extend(true, $[deck].defaults, {
		selectors: {
			hashLink: '.deck-permalink'
		},
		
		hashPrefix: 'slide-',
		preventFragmentScroll: true
	});
	
	
	$d.bind('deck.init', function() {
	   var opts = $[deck]('getOptions');
		$internals = $(),
		slides = $[deck]('getSlides');
		
		$.each(slides, function(i, $el) {
			var hash;
			
			/* Hand out ids to the unfortunate slides born without them */
			if (!$el.attr('id') || $el.data('deckAssignedId') === $el.attr('id')) {
				$el.attr('id', opts.hashPrefix + i);
				$el.data('deckAssignedId', opts.hashPrefix + i);
			}
			
			hash ='#' + $el.attr('id');
			
			/* Deep link to slides on init */
			if (hash === window.location.hash) {
				$[deck]('go', i);
			}
			
			/* Add internal links to this slide */
			$internals = $internals.add('a[href="' + hash + '"]');
		});
		
		if (!Modernizr.hashchange) {
			/* Set up internal links using click for the poor browsers
			without a hashchange event. */
			$internals.unbind('click.deckhash').bind('click.deckhash', function(e) {
				goByHash($(this).attr('href'));
			});
		}
		
		/* Set up first id container state class */
		if (slides.length) {
			$[deck]('getContainer').addClass(opts.classes.onPrefix + $[deck]('getSlide').attr('id'));
		};
	})
	/* Update permalink, address bar, and state class on a slide change */
	.bind('deck.change', function(e, from, to) {
		var hash = '#' + $[deck]('getSlide', to).attr('id'),
		hashPath = window.location.href.replace(/#.*/, '') + hash,
		opts = $[deck]('getOptions'),
		osp = opts.classes.onPrefix,
		$c = $[deck]('getContainer');
		
		$c.removeClass(osp + $[deck]('getSlide', from).attr('id'));
		$c.addClass(osp + $[deck]('getSlide', to).attr('id'));
		
		$(opts.selectors.hashLink).attr('href', hashPath);
		if (Modernizr.history) {
			window.history.replaceState({}, "", hashPath);
		}
	});
	
	/* Deals with internal links in modern browsers */
	$window.bind('hashchange.deckhash', function(e) {
		if (e.originalEvent && e.originalEvent.newURL) {
			goByHash(e.originalEvent.newURL);
		}
		else {
			goByHash(window.location.hash);
		}
	})
	/* Prevent scrolling on deep links */
	.bind('load', function() {
		if ($[deck]('getOptions').preventFragmentScroll) {
			$[deck]('getContainer').scrollLeft(0).scrollTop(0);
		}
	});
})(jQuery, 'deck', this);

/*!
Deck JS - deck.menu
Copyright (c) 2011 Caleb Troughton
Dual licensed under the MIT license and GPL license.
https://github.com/imakewebthings/deck.js/blob/master/MIT-license.txt
https://github.com/imakewebthings/deck.js/blob/master/GPL-license.txt
*/

/*
This module adds the methods and key binding to show and hide a menu of all
slides in the deck. The deck menu state is indicated by the presence of a class
on the deck container.
*/
(function($, deck, undefined) {
	var $d = $(document),
	rootSlides; // Array of top level slides
	
	/*
	Extends defaults/options.
	
	options.classes.menu
		This class is added to the deck container when showing the slide menu.
	
	options.keys.menu
		The numeric keycode used to toggle between showing and hiding the slide
		menu.
		
	options.touch.doubletapWindow
		Two consecutive touch events within this number of milliseconds will
		be considered a double tap, and will toggle the menu on touch devices.
	*/
	$.extend(true, $[deck].defaults, {
		classes: {
			menu: 'deck-menu'
		},
		
		keys: {
			menu: 77 // m
		},
		
		touch: {
			doubletapWindow: 400
		}
	});

	/*
	jQuery.deck('showMenu')
	
	Shows the slide menu by adding the class specified by the menu class option
	to the deck container.
	*/
	$[deck]('extend', 'showMenu', function() {
		var $c = $[deck]('getContainer'),
		opts = $[deck]('getOptions');
		
		if ($c.hasClass(opts.classes.menu)) return;
		
		// Hide through loading class to short-circuit transitions (perf)
		$c.addClass([opts.classes.loading, opts.classes.menu].join(' '));
		
		/* Forced to do this in JS until CSS learns second-grade math. Save old
		style value for restoration when menu is hidden. */
		if (Modernizr.csstransforms) {
			$.each(rootSlides, function(i, $slide) {
				$slide.data('oldStyle', $slide.attr('style'));
				$slide.css({
					'position': 'absolute',
					'left': ((i % 4) * 25) + '%',
					'top': (Math.floor(i / 4) * 25) + '%'
				});
			});
		}
		
		// Need to ensure the loading class renders first, then remove
		window.setTimeout(function() {
			$c.removeClass(opts.classes.loading)
				.scrollTop($[deck]('getSlide').offset().top);
		}, 0);
	});

	/*
	jQuery.deck('hideMenu')
	
	Hides the slide menu by removing the class specified by the menu class
	option from the deck container.
	*/
	$[deck]('extend', 'hideMenu', function() {
		var $c = $[deck]('getContainer'),
		opts = $[deck]('getOptions');
		
		if (!$c.hasClass(opts.classes.menu)) return;
		
		$c.removeClass(opts.classes.menu);
		$c.addClass(opts.classes.loading);
		
		/* Restore old style value */
		if (Modernizr.csstransforms) {
			$.each(rootSlides, function(i, $slide) {
				var oldStyle = $slide.data('oldStyle');

				$slide.attr('style', oldStyle ? oldStyle : '');
			});
		}
		
		window.setTimeout(function() {
			$c.removeClass(opts.classes.loading).scrollTop(0);
		}, 0);
	});

	/*
	jQuery.deck('toggleMenu')
	
	Toggles between showing and hiding the slide menu.
	*/
	$[deck]('extend', 'toggleMenu', function() {
		$[deck]('getContainer').hasClass($[deck]('getOptions').classes.menu) ?
		$[deck]('hideMenu') : $[deck]('showMenu');
	});

	$d.bind('deck.init', function() {
		var opts = $[deck]('getOptions'),
		touchEndTime = 0,
		currentSlide,
		slideTest = $.map([
			opts.classes.before,
			opts.classes.previous,
			opts.classes.current,
			opts.classes.next,
			opts.classes.after
		], function(el, i) {
			return '.' + el;
		}).join(', ');
		
		// Build top level slides array
		rootSlides = [];
		$.each($[deck]('getSlides'), function(i, $el) {
			if (!$el.parentsUntil(opts.selectors.container, slideTest).length) {
				rootSlides.push($el);
			}
		});
		
		// Bind key events
		$d.unbind('keydown.deckmenu').bind('keydown.deckmenu', function(e) {
			if (e.which === opts.keys.menu || $.inArray(e.which, opts.keys.menu) > -1) {
				$[deck]('toggleMenu');
				e.preventDefault();
			}
		});
		
		// Double tap to toggle slide menu for touch devices
		$[deck]('getContainer').unbind('touchstart.deckmenu').bind('touchstart.deckmenu', function(e) {
			currentSlide = $[deck]('getSlide');
		})
		.unbind('touchend.deckmenu').bind('touchend.deckmenu', function(e) {
			var now = Date.now();
			
			// Ignore this touch event if it caused a nav change (swipe)
			if (currentSlide !== $[deck]('getSlide')) return;
			
			if (now - touchEndTime < opts.touch.doubletapWindow) {
				$[deck]('toggleMenu');
				e.preventDefault();
			}
			touchEndTime = now;
		});
		
		// Selecting slides from the menu
		$.each($[deck]('getSlides'), function(i, $s) {
			$s.unbind('click.deckmenu').bind('click.deckmenu', function(e) {
				if (!$[deck]('getContainer').hasClass(opts.classes.menu)) return;

				$[deck]('go', i);
				$[deck]('hideMenu');
				e.stopPropagation();
				e.preventDefault();
			});
		});
	})
	.bind('deck.change', function(e, from, to) {
		var container = $[deck]('getContainer');
		
		if (container.hasClass($[deck]('getOptions').classes.menu)) {
			container.scrollTop($[deck]('getSlide', to).offset().top);
		}
	});
})(jQuery, 'deck');

/*!
Deck JS - deck.navigation
Copyright (c) 2011 Caleb Troughton
Dual licensed under the MIT license and GPL license.
https://github.com/imakewebthings/deck.js/blob/master/MIT-license.txt
https://github.com/imakewebthings/deck.js/blob/master/GPL-license.txt
*/

/*
This module adds clickable previous and next links to the deck.
*/
(function($, deck, undefined) {
	var $d = $(document),
	
	/* Updates link hrefs, and disabled states if last/first slide */
	updateButtons = function(e, from, to) {
		var opts = $[deck]('getOptions'),
		last = $[deck]('getSlides').length - 1,
		prevSlide = $[deck]('getSlide', to - 1),
		nextSlide = $[deck]('getSlide', to + 1),
		hrefBase = window.location.href.replace(/#.*/, ''),
		prevId = prevSlide ? prevSlide.attr('id') : undefined,
		nextId = nextSlide ? nextSlide.attr('id') : undefined;
		
		$(opts.selectors.previousLink)
			.toggleClass(opts.classes.navDisabled, !to)
			.attr('href', hrefBase + '#' + (prevId ? prevId : ''));
		$(opts.selectors.nextLink)
			.toggleClass(opts.classes.navDisabled, to === last)
			.attr('href', hrefBase + '#' + (nextId ? nextId : ''));
	};
	
	/*
	Extends defaults/options.
	
	options.classes.navDisabled
		This class is added to a navigation link when that action is disabled.
		It is added to the previous link when on the first slide, and to the
		next link when on the last slide.
		
	options.selectors.nextLink
		The elements that match this selector will move the deck to the next
		slide when clicked.
		
	options.selectors.previousLink
		The elements that match this selector will move to deck to the previous
		slide when clicked.
	*/
	$.extend(true, $[deck].defaults, {
		classes: {
			navDisabled: 'deck-nav-disabled'
		},
		
		selectors: {
			nextLink: '.deck-next-link',
			previousLink: '.deck-prev-link'
		}
	});

	$d.bind('deck.init', function() {
		var opts = $[deck]('getOptions'),
		slides = $[deck]('getSlides'),
		$current = $[deck]('getSlide'),
		ndx;
		
		// Setup prev/next link events
		$(opts.selectors.previousLink)
		.unbind('click.decknavigation')
		.bind('click.decknavigation', function(e) {
			$[deck]('prev');
			e.preventDefault();
		});
		
		$(opts.selectors.nextLink)
		.unbind('click.decknavigation')
		.bind('click.decknavigation', function(e) {
			$[deck]('next');
			e.preventDefault();
		});
		
		// Find where we started in the deck and set initial states
		$.each(slides, function(i, $slide) {
			if ($slide === $current) {
				ndx = i;
				return false;
			}
		});
		updateButtons(null, ndx, ndx);
	})
	.bind('deck.change', updateButtons);
})(jQuery, 'deck');

/*!
Deck JS - deck.status
Copyright (c) 2011 Caleb Troughton
Dual licensed under the MIT license and GPL license.
https://github.com/imakewebthings/deck.js/blob/master/MIT-license.txt
https://github.com/imakewebthings/deck.js/blob/master/GPL-license.txt
*/

/*
This module adds a (current)/(total) style status indicator to the deck.
*/
(function($, deck, undefined) {
	var $d = $(document),
	
	updateCurrent = function(e, from, to) {
		var opts = $[deck]('getOptions');
		
		$(opts.selectors.statusCurrent).text(opts.countNested ?
			to + 1 :
			$[deck]('getSlide', to).data('rootSlide')
		);
	};
	
	/*
	Extends defaults/options.
	
	options.selectors.statusCurrent
		The element matching this selector displays the current slide number.
		
	options.selectors.statusTotal
		The element matching this selector displays the total number of slides.
		
	options.countNested
		If false, only top level slides will be counted in the current and
		total numbers.
	*/
	$.extend(true, $[deck].defaults, {
		selectors: {
			statusCurrent: '.deck-status-current',
			statusTotal: '.deck-status-total'
		},
		
		countNested: true
	});
	
	$d.bind('deck.init', function() {
		var opts = $[deck]('getOptions'),
		slides = $[deck]('getSlides'),
		$current = $[deck]('getSlide'),
		ndx;
		
		// Set total slides once
		if (opts.countNested) {
			$(opts.selectors.statusTotal).text(slides.length);
		}
		else {
			/* Determine root slides by checking each slide's ancestor tree for
			any of the slide classes. */
			var rootIndex = 1,
			slideTest = $.map([
				opts.classes.before,
				opts.classes.previous,
				opts.classes.current,
				opts.classes.next,
				opts.classes.after
			], function(el, i) {
				return '.' + el;
			}).join(', ');
			
			/* Store the 'real' root slide number for use during slide changes. */
			$.each(slides, function(i, $el) {
				var $parentSlides = $el.parentsUntil(opts.selectors.container, slideTest);

				$el.data('rootSlide', $parentSlides.length ?
					$parentSlides.last().data('rootSlide') :
					rootIndex++
				);
			});
			
			$(opts.selectors.statusTotal).text(rootIndex - 1);
		}
		
		// Find where we started in the deck and set initial state
		$.each(slides, function(i, $el) {
			if ($el === $current) {
				ndx = i;
				return false;
			}
		});
		updateCurrent(null, ndx, ndx);
	})
	/* Update current slide number with each change event */
	.bind('deck.change', updateCurrent);
})(jQuery, 'deck');
`
