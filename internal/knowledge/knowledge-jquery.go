package knowledge

func getJQueryBits() []KnowledgeBit {
	return []KnowledgeBit{
		// Selectors
		{"What are the different types of jQuery selectors?", "jQuery selectors tutorial"},
		{"How do attribute selectors work in jQuery?", "jQuery attribute selectors example"},
		{"What are pseudo selectors in jQuery?", "jQuery pseudo selectors tutorial"},
		{"Explain hierarchy selectors in jQuery", "jQuery hierarchy selectors example"},
		{"How do form selectors work in jQuery?", "jQuery form selectors tutorial"},

		// DOM Manipulation
		{"What is the difference between .html(), .text(), and .val()?", "jQuery html text val methods tutorial"},
		{"Explain .append() vs .prepend() vs .after() vs .before()", "jQuery append prepend tutorial"},
		{"What is the difference between .remove(), .empty(), and .detach()?", "jQuery remove empty detach tutorial"},
		{"What is the difference between .attr() and .prop()?", "jQuery attr vs prop tutorial"},
		{"How do .addClass(), .removeClass(), and .toggleClass() work?", "jQuery addClass removeClass tutorial"},
		{"How does the .css() method work for getting and setting styles?", "jQuery css method tutorial"},
		{"How does .clone() work and what are its options?", "jQuery clone method example"},

		// Traversal
		{"What is the difference between .find(), .children(), and .parent()?", "jQuery traversal tutorial"},
		{"Explain .siblings(), .next(), and .prev() methods", "jQuery siblings next prev example"},
		{"What does .closest() do and how is it different from .parent()?", "jQuery closest vs parent tutorial"},
		{"How do .filter(), .not(), and .is() work?", "jQuery filter methods tutorial"},

		// Event Handling
		{"What is the difference between 'this' and '$(this)' in jQuery?", "jQuery this keyword tutorial"},
		{"Why don't arrow functions work correctly with jQuery's 'this'?", "jQuery arrow function this tutorial"},
		{"How does .each() work in jQuery?", "jQuery each method tutorial"},
		{"Explain .on() for event binding in jQuery", "jQuery on event binding tutorial"},
		{"How do you remove event handlers with .off()?", "jQuery off method tutorial"},
		{"What is event delegation in jQuery and why is it useful?", "jQuery event delegation tutorial"},
		{"What is the event object in jQuery and what properties does it have?", "jQuery event object tutorial"},
		{"How does .trigger() work for programmatic events?", "jQuery trigger method example"},

		// Effects & Animation
		{"What are the basic jQuery visibility methods?", "jQuery show hide toggle tutorial"},
		{"How do fade effects work in jQuery?", "jQuery fade effects tutorial"},
		{"Explain slide effects in jQuery", "jQuery slide effects tutorial"},
		{"How does the .animate() method work?", "jQuery animate tutorial"},
		{"What is animation queue buildup and how do you prevent it?", "jQuery animation queue tutorial"},
		{"How does .stop() work to control animations?", "jQuery stop animation example"},

		// AJAX
		{"How does $.ajax() work for HTTP requests?", "jQuery ajax tutorial"},
		{"What are the shorthand AJAX methods $.get() and $.post()?", "jQuery get post tutorial"},
		{"How do you fetch JSON data with jQuery?", "jQuery getJSON tutorial"},
		{"What does the .load() method do?", "jQuery load method tutorial"},
		{"How do you use promises with jQuery AJAX?", "jQuery ajax promise tutorial"},

		// Plugins
		{"How do you create a jQuery plugin?", "jQuery plugin development tutorial"},
		{"How do you handle plugin options with defaults?", "jQuery plugin options tutorial"},
		{"What is the chainability pattern in jQuery plugins?", "jQuery plugin chaining example"},

		// Deferred/Promises
		{"What is $.Deferred and how does it work?", "jQuery Deferred tutorial"},
		{"How do you wait for multiple AJAX requests with $.when?", "jQuery when multiple ajax example"},

		// Forms
		{"How do you serialize form data with jQuery?", "jQuery serialize form tutorial"},
		{"What is the difference between .serialize() and .serializeArray()?", "jQuery serialize vs serializeArray example"},

		// Dimensions & Position
		{"What is the difference between .offset() and .position()?", "jQuery offset vs position tutorial"},
		{"Explain .width() vs .innerWidth() vs .outerWidth()", "jQuery dimensions tutorial"},
		{"How do you work with scroll position in jQuery?", "jQuery scroll position tutorial"},

		// Best Practices
		{"Why should you cache jQuery selectors?", "jQuery performance tutorial"},
		{"What is method chaining in jQuery?", "jQuery method chaining tutorial"},
		{"What is $(document).ready() and when is it used?", "jQuery document ready tutorial"},
		{"What is the difference between $(document).ready() and window.onload?", "jQuery ready vs load explained"},
		{"How do you avoid conflicts with other libraries using $.noConflict()?", "jQuery noConflict tutorial"},
		{"How do you use namespaced events in jQuery?", "jQuery namespaced events example"},

		// Modern Alternatives
		{"How do you convert jQuery selectors to vanilla JavaScript?", "jQuery to vanilla JavaScript guide"},
		{"What is the vanilla JavaScript equivalent of $.ajax()?", "fetch API vs jQuery ajax tutorial"},
	}
}
