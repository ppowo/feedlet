package knowledge

func getJSPBits() []KnowledgeBit {
	return []KnowledgeBit{
		// Basics
		{"What is JSP and how does it work?", "JSP tutorial beginners"},
		{"What are the different types of JSP syntax?", "JSP syntax tutorial"},
		{"Explain the JSP lifecycle", "JSP lifecycle tutorial"},
		{"What is the difference between <% %>, <%= %>, and <%! %>?", "JSP scriptlet vs expression tutorial"},
		{"What are JSP directives and what types are there?", "JSP directives tutorial"},
		{"What are JSP comments and how are they different from HTML comments?", "JSP comments tutorial"},

		// Implicit Objects
		{"What are JSP implicit objects?", "JSP implicit objects tutorial"},
		{"How does the request object work in JSP?", "JSP request object tutorial"},
		{"How does the response object work in JSP?", "JSP response object tutorial"},
		{"What is the session object and how is it used?", "JSP session object tutorial"},
		{"What is the application object in JSP?", "JSP application object tutorial"},
		{"What is the out object in JSP?", "JSP out object tutorial"},
		{"What is pageContext and when would you use it?", "JSP pageContext tutorial"},
		{"What is the exception object and when is it available?", "JSP exception handling tutorial"},

		// Scopes
		{"What are the four scopes in JSP?", "JSP scopes tutorial"},
		{"What is the difference between request and session scope?", "JSP request vs session scope tutorial"},
		{"When should you use each JSP scope?", "JSP scope best practices guide"},
		{"How does EL search through scopes?", "JSP EL scope tutorial"},

		// Expression Language
		{"What is Expression Language (EL) in JSP?", "JSP Expression Language tutorial"},
		{"How do you access request parameters with EL?", "JSP EL param tutorial"},
		{"What operators are available in EL?", "JSP EL operators tutorial"},
		{"How does property navigation work in EL?", "JSP EL property access tutorial"},

		// JSTL
		{"What is JSTL and why use it?", "JSTL tutorial beginners"},
		{"How does c:if work and why doesn't it have an else?", "JSTL c:if tutorial"},
		{"How does c:choose work for if-else logic?", "JSTL c:choose tutorial"},
		{"How does c:forEach work for iterating collections?", "JSTL c:forEach tutorial"},
		{"What do c:set and c:remove do?", "JSTL c:set tutorial"},
		{"How does c:out prevent XSS attacks?", "JSTL c:out tutorial"},
		{"How do you format dates and numbers with JSTL?", "JSTL formatting tutorial"},

		// Advanced
		{"What is the difference between forward and include?", "JSP forward vs include tutorial"},
		{"How do you handle errors in JSP?", "JSP error page tutorial"},
		{"What is the difference between JSP and Servlet?", "JSP vs Servlet tutorial"},
		{"What are the different session tracking methods?", "JSP session tracking tutorial"},
		{"Why is JSP not thread-safe by default?", "JSP thread safety tutorial"},
		{"How do you create custom tags in JSP?", "JSP custom tags tutorial"},

		// Best Practices
		{"Why is JSTL preferred over scriptlets?", "JSTL vs scriptlets tutorial"},
		{"How does MVC work in a JSP application?", "JSP MVC tutorial"},
		{"When would you encounter JSP in modern development?", "JSP modern development guide"},
	}
}
