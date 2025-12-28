package knowledge

func getJQueryBits() []KnowledgeBit {
	return []KnowledgeBit{
		// Selectors
		{"What are the different types of jQuery selectors?", "jQuery selectors"},
		{"How do attribute selectors work in jQuery?", "jQuery attribute selectors"},
		{"What are pseudo selectors in jQuery?", "jQuery pseudo selectors"},
		{"Explain hierarchy selectors in jQuery", "jQuery hierarchy selectors"},
		{"How do form selectors work in jQuery?", "jQuery form selectors"},

		// DOM Manipulation
		{"What is the difference between .html(), .text(), and .val()?", "jQuery html text val"},
		{"Explain .append() vs .prepend() vs .after() vs .before()", "jQuery append prepend after before"},
		{"What is the difference between .remove(), .empty(), and .detach()?", "jQuery remove empty detach"},
		{"What is the difference between .attr() and .prop()?", "jQuery attr vs prop"},
		{"How do .addClass(), .removeClass(), and .toggleClass() work?", "jQuery addClass removeClass toggleClass"},
		{"How does the .css() method work for getting and setting styles?", "jQuery css method"},
		{"How does .clone() work and what are its options?", "jQuery clone method"},

		// Traversal
		{"What is the difference between .find(), .children(), and .parent()?", "jQuery find children parent traversal"},
		{"Explain .siblings(), .next(), and .prev() methods", "jQuery siblings next prev"},
		{"What does .closest() do and how is it different from .parent()?", "jQuery closest vs parent"},
		{"How do .filter(), .not(), and .is() work?", "jQuery filter not is methods"},

		// Event Handling
		{"What is the difference between 'this' and '$(this)' in jQuery?", "jQuery this vs $(this)"},
		{"Why don't arrow functions work correctly with jQuery's 'this'?", "jQuery arrow function this"},
		{"How does .each() work in jQuery?", "jQuery each method"},
		{"Explain .on() for event binding in jQuery", "jQuery on method event binding"},
		{"How do you remove event handlers with .off()?", "jQuery off remove event handler"},
		{"What is event delegation in jQuery and why is it useful?", "jQuery event delegation"},
		{"What is the event object in jQuery and what properties does it have?", "jQuery event object properties"},
		{"How does .trigger() work for programmatic events?", "jQuery trigger method"},

		// Effects & Animation
		{"What are the basic jQuery visibility methods?", "jQuery show hide toggle"},
		{"How do fade effects work in jQuery?", "jQuery fadeIn fadeOut fadeTo"},
		{"Explain slide effects in jQuery", "jQuery slideDown slideUp slideToggle"},
		{"How does the .animate() method work?", "jQuery animate custom animation"},
		{"What is animation queue buildup and how do you prevent it?", "jQuery animation queue stop"},
		{"How does .stop() work to control animations?", "jQuery stop animation"},

		// AJAX
		{"How does $.ajax() work for HTTP requests?", "jQuery ajax method"},
		{"What are the shorthand AJAX methods $.get() and $.post()?", "jQuery get post ajax"},
		{"How do you fetch JSON data with jQuery?", "jQuery getJSON ajax JSON"},
		{"What does the .load() method do?", "jQuery load method HTML"},
		{"How do you use promises with jQuery AJAX?", "jQuery ajax promise done fail"},

		// Plugins
		{"How do you create a jQuery plugin?", "jQuery plugin development"},
		{"How do you handle plugin options with defaults?", "jQuery plugin options defaults"},
		{"What is the chainability pattern in jQuery plugins?", "jQuery plugin chaining"},

		// Deferred/Promises
		{"What is $.Deferred and how does it work?", "jQuery Deferred promise"},
		{"How do you wait for multiple AJAX requests with $.when?", "jQuery when multiple ajax"},

		// Forms
		{"How do you serialize form data with jQuery?", "jQuery serialize form data"},
		{"What is the difference between .serialize() and .serializeArray()?", "jQuery serialize vs serializeArray"},

		// Dimensions & Position
		{"What is the difference between .offset() and .position()?", "jQuery offset vs position"},
		{"Explain .width() vs .innerWidth() vs .outerWidth()", "jQuery width innerWidth outerWidth"},
		{"How do you work with scroll position in jQuery?", "jQuery scrollTop scrollLeft"},

		// Best Practices
		{"Why should you cache jQuery selectors?", "jQuery selector caching performance"},
		{"What is method chaining in jQuery?", "jQuery method chaining"},
		{"What is $(document).ready() and when is it used?", "jQuery document ready"},
		{"What is the difference between $(document).ready() and window.onload?", "jQuery ready vs load event"},
		{"How do you avoid conflicts with other libraries using $.noConflict()?", "jQuery noConflict"},
		{"How do you use namespaced events in jQuery?", "jQuery namespaced events"},

		// Modern Alternatives
		{"How do you convert jQuery selectors to vanilla JavaScript?", "jQuery to vanilla JavaScript querySelector"},
		{"What is the vanilla JavaScript equivalent of $.ajax()?", "fetch API vs jQuery ajax"},
	}
}
