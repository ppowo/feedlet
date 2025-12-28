package knowledge

func getJSPBits() []KnowledgeBit {
	return []KnowledgeBit{
		// Basics
		{"What is JSP and how does it work?", "JSP JavaServer Pages"},
		{"What are the different types of JSP syntax?", "JSP scriptlet expression directive"},
		{"Explain the JSP lifecycle", "JSP lifecycle translation compilation"},
		{"What is the difference between <% %>, <%= %>, and <%! %>?", "JSP scriptlet expression declaration"},
		{"What are JSP directives and what types are there?", "JSP page directive taglib include"},
		{"What are JSP comments and how are they different from HTML comments?", "JSP comments security"},

		// Implicit Objects
		{"What are JSP implicit objects?", "JSP implicit objects"},
		{"How does the request object work in JSP?", "JSP request object getParameter"},
		{"How does the response object work in JSP?", "JSP response object redirect"},
		{"What is the session object and how is it used?", "JSP session object HttpSession"},
		{"What is the application object in JSP?", "JSP application object ServletContext"},
		{"What is the out object in JSP?", "JSP out object JspWriter"},
		{"What is pageContext and when would you use it?", "JSP pageContext object"},
		{"What is the exception object and when is it available?", "JSP exception object error page"},

		// Scopes
		{"What are the four scopes in JSP?", "JSP scopes page request session application"},
		{"What is the difference between request and session scope?", "JSP request vs session scope"},
		{"When should you use each JSP scope?", "JSP scope best practices"},
		{"How does EL search through scopes?", "JSP EL scope search order"},

		// Expression Language
		{"What is Expression Language (EL) in JSP?", "JSP Expression Language EL"},
		{"How do you access request parameters with EL?", "JSP EL param implicit object"},
		{"What operators are available in EL?", "JSP EL operators"},
		{"How does property navigation work in EL?", "JSP EL dot notation property access"},

		// JSTL
		{"What is JSTL and why use it?", "JSTL JSP Standard Tag Library"},
		{"How does c:if work and why doesn't it have an else?", "JSTL c:if conditional"},
		{"How does c:choose work for if-else logic?", "JSTL c:choose c:when c:otherwise"},
		{"How does c:forEach work for iterating collections?", "JSTL c:forEach loop varStatus"},
		{"What do c:set and c:remove do?", "JSTL c:set c:remove variable"},
		{"How does c:out prevent XSS attacks?", "JSTL c:out escapeXml XSS"},
		{"How do you format dates and numbers with JSTL?", "JSTL fmt:formatDate fmt:formatNumber"},

		// Advanced
		{"What is the difference between forward and include?", "JSP forward vs include RequestDispatcher"},
		{"How do you handle errors in JSP?", "JSP error page exception handling"},
		{"What is the difference between JSP and Servlet?", "JSP vs Servlet MVC"},
		{"What are the different session tracking methods?", "JSP session tracking cookies URL rewriting"},
		{"Why is JSP not thread-safe by default?", "JSP thread safety concurrency"},
		{"How do you create custom tags in JSP?", "JSP custom tags tag files"},

		// Best Practices
		{"Why is JSTL preferred over scriptlets?", "JSTL vs scriptlets best practices"},
		{"How does MVC work in a JSP application?", "JSP MVC pattern servlet controller"},
		{"When would you encounter JSP in modern development?", "JSP legacy enterprise applications"},
	}
}
