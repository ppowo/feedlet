package knowledge

// getJSPBits returns JSP (JavaServer Pages) knowledge tidbits
func getJSPBits() []KnowledgeBit {
	return []KnowledgeBit{
		// What is JSP?
		{
			Title:       "What is JSP?",
			Description: "JavaServer Pages. HTML template files (.jsp) that run Java code on the server to generate dynamic HTML. Server processes the .jsp file and sends plain HTML to browser.",
			Code:        "<!-- user.jsp -->\n<html>\n<body>\n    <%\n        // Java code runs on SERVER\n        String username = request.getParameter(\"name\");\n        int visits = getVisitCount();\n    %>\n    \n    <h1>Hello <%= username %></h1>\n    <p>You've visited <%= visits %> times</p>\n    \n    <% if (isAdmin(username)) { %>\n        <a href=\"admin.jsp\">Admin Panel</a>\n    <% } %>\n</body>\n</html>",
			Category:    "Basics",
		},
		{
			Title:       "Three JSP Syntaxes You'll See",
			Description: "Legacy code mixes 3 styles: Scriptlets (<% %>), Expressions (<%= %>), and JSTL tags (<c:...>). They all generate the same HTML, just different syntax.",
			Code:        "<!-- 1. SCRIPTLET <% %> - Java statements -->\n<%\n    String name = request.getParameter(\"name\");\n    int age = Integer.parseInt(request.getParameter(\"age\"));\n    if (age >= 18) {\n        // do something\n    }\n%>\n\n<!-- 2. EXPRESSION <%= %> - Print a value -->\n<p>Name: <%= name %></p>\n<p>Age: <%= age %></p>\n\n<!-- 3. JSTL TAGS - Cleaner syntax -->\n<c:if test=\"${age >= 18}\">\n    <p>Adult</p>\n</c:if>\n\n<!-- All three work, often mixed in same file -->\n<!-- JSTL is cleaner but needs <%@ taglib %> import -->",
			Category:    "Basics",
		},
		{
			Title:       "How JSP Works",
			Description: "Browser requests .jsp → Server translates to servlet (Java class) → Compiles it → Runs it → Sends HTML back. First request is slow (compile), rest are fast (already compiled).",
			Code:        "// What happens:\n// 1. Browser: GET /user.jsp?id=123\n// 2. Server: user.jsp → user_jsp.java (servlet)\n// 3. Server: Compile user_jsp.java → user_jsp.class\n// 4. Server: Run user_jsp.class (executes Java code)\n// 5. Server: Sends generated HTML to browser\n// 6. Browser: Receives plain HTML, displays it\n\n// Next time someone visits user.jsp:\n// - Skip steps 2-3 (already compiled)\n// - Just run the class and send HTML\n\n// That's why first load is slow, rest are fast\n// Server does the work, browser just displays HTML",
			Category:    "Basics",
		},

		// JSP Basics
		{
			Title:       "Scriptlet Basics: <% %> and <%= %>",
			Description: "<% %> executes Java code. <%= %> prints a value to HTML. Code inside <% %> runs on server for each page request. Local variables are request-scoped (not shared between users).",
			Code:        "<!-- <% %> SCRIPTLET - Execute Java code -->\n<%\n    // These are local variables - each user gets their own\n    String name = request.getParameter(\"name\");\n    int age = Integer.parseInt(request.getParameter(\"age\"));\n    \n    // Can do any Java here\n    if (name == null) {\n        name = \"Guest\";\n    }\n%>\n\n<!-- <%= %> EXPRESSION - Print a value -->\n<p>Hello <%= name %></p>\n<p>You are <%= age %> years old</p>\n<p>Next year: <%= age + 1 %></p>\n\n<!-- Expression is shorthand for: <% out.print(value); %> -->",
			Category:    "Basics",
		},
		{
			Title:       "<%@ %> Directives - Page Configuration",
			Description: "<%@ %> sets page-level configuration. Goes at the top of .jsp file. Most common: page (settings), include (files), taglib (JSTL).",
			Code:        "<!-- page directive - page settings -->\n<%@ page contentType=\"text/html;charset=UTF-8\" language=\"java\" %>\n<%@ page import=\"java.util.List, com.myapp.User\" %>\n<%@ page errorPage=\"error.jsp\" %>\n<%@ page session=\"true\" %> <!-- default, can use session -->\n\n<!-- include directive - paste another file's content here -->\n<%@ include file=\"header.jsp\" %>\n<%@ include file=\"/WEB-INF/includes/navigation.jsp\" %>\n\n<!-- taglib directive - import tag library (JSTL) -->\n<%@ taglib uri=\"http://java.sun.com/jsp/jstl/core\" prefix=\"c\" %>\n<%@ taglib uri=\"http://java.sun.com/jsp/jstl/fmt\" prefix=\"fmt\" %>",
			Category:    "Basics",
		},
		{
			Title:       "<%! %> Declaration - Shared State",
			Description: "<%! %> creates class-level variables. ALL users share the same variable! User A increments it, User B sees the new value. Dangerous for counters - multiple users at once = race condition.",
			Code:        "<%!\n    // This variable is SHARED by ALL users!\n    private int pageViews = 0;\n%>\n\n<%\n    // User A visits: pageViews becomes 1\n    // User B visits: pageViews becomes 2\n    // User A refreshes: sees 2 (not their own count!)\n    pageViews++;\n    out.println(\"Total visits: \" + pageViews);\n%>\n\n<!-- Problem: Two users click at exact same time -->\n<!-- User A reads: 5, increments to 6 -->\n<!-- User B reads: 5, increments to 6 -->\n<!-- Both write 6! Lost one count (should be 7) -->\n\n<!-- Safe alternative: use application scope -->\n<%\n    Integer count = (Integer) application.getAttribute(\"count\");\n    if (count == null) count = 0;\n    count++;\n    application.setAttribute(\"count\", count);\n    // Still shared, but explicit about it\n%>",
			Category:    "Basics",
		},
		{
			Title:       "JSP Comments",
			Description: "<%-- --%> for JSP comments (not sent to client). <!-- --> for HTML comments (sent to client). Use JSP comments for sensitive info.",
			Code:        "<%-- JSP comment - NOT in HTML source, not sent to browser --%>\n<%-- Password: secret123 --%>\n<%-- This code is commented out: <%= user.getPassword() %> --%>\n\n<!-- HTML comment - VISIBLE in page source, sent to browser -->\n<!-- TODO: fix this later -->\n<!-- User can see this in View Source -->",
			Category:    "Basics",
		},

		// Implicit Objects
		{
			Title:       "request Object",
			Description: "JSP gives you 'request' variable for free (like req in Express.js). Access URL params, form data, headers. Think of it like Express req or FastAPI request.",
			Code:        "<!-- URL: page.jsp?name=John&age=25 -->\n<%\n    // Get URL/form parameters (like req.query or req.body)\n    String name = request.getParameter(\"name\"); // \"John\"\n    String age = request.getParameter(\"age\");   // \"25\"\n    \n    // Like Express: app.get('/page', (req, res) => {\n    //   const name = req.query.name;\n    // })\n    \n    // Get headers (like req.headers)\n    String userAgent = request.getHeader(\"User-Agent\");\n    String ip = request.getRemoteAddr();\n%>",
			Category:    "Implicit Objects",
		},
		{
			Title:       "response Object",
			Description: "The 'response' variable (like res in Express). Redirect users, set cookies, send errors. In modern apps, you'd use res.redirect() or res.json(). JSP is similar.",
			Code:        "<%\n    // Redirect (like res.redirect() in Express)\n    response.sendRedirect(\"login.jsp\");\n    \n    // Set cookie (like res.cookie() in Express)\n    Cookie cookie = new Cookie(\"username\", \"john\");\n    cookie.setMaxAge(3600); // 1 hour\n    response.addCookie(cookie);\n    \n    // Send 404 (like res.status(404).send())\n    response.sendError(404, \"Page not found\");\n    \n    // Modern equivalent:\n    // res.status(404).json({error: \"Not found\"})\n%>",
			Category:    "Implicit Objects",
		},
		{
			Title:       "session Object",
			Description: "Server-side session storage (like req.session in Express). Store user data between requests (login state, cart). Each user gets their own session. Stored on server, not browser.",
			Code:        "<%\n    // Save to session (like req.session.username = 'john')\n    session.setAttribute(\"username\", \"john\");\n    session.setAttribute(\"cart\", cartObject);\n    \n    // Get from session (survives page reloads!)\n    String user = (String) session.getAttribute(\"username\");\n    if (user == null) {\n        // Not logged in, redirect\n        response.sendRedirect(\"login.jsp\");\n    }\n    \n    // Logout - destroy session\n    session.invalidate(); // like req.session.destroy()\n%>",
			Category:    "Implicit Objects",
		},
		{
			Title:       "application Object",
			Description: "Global app storage shared by ALL users (like global variable on server). Survives until server restart. Used for site-wide counters/config. Available in <% %> blocks.",
			Code:        "<%\n    // 'application' = global storage for entire app\n    // ALL users share this! (Not per-user like session)\n    \n    // Store total site visitors (shared by everyone)\n    application.setAttribute(\"totalVisits\", 1000);\n    \n    // Get it back\n    Integer visits = (Integer) application.getAttribute(\"totalVisits\");\n    out.println(\"Total site visits: \" + visits);\n    \n    // Like storing in a global variable on your server\n    // In Node: global.totalVisits = 1000;\n%>",
			Category:    "Implicit Objects",
		},
		{
			Title:       "out Object",
			Description: "'out' prints HTML from Java code (like res.write() in Express). Goes inside <% %> tags. OLD WAY - don't use this! Use EL ${} instead. Shown here so you recognize it in old code.",
			Code:        "<!-- OLD WAY (don't do this!) -->\n<%\n    // 'out' is available automatically in <% %> blocks\n    out.println(\"<h1>Hello</h1>\"); // prints HTML\n    out.print(\"Name: \" + name);     // concatenate and print\n%>\n\n<!-- MODERN WAY (do this instead!) -->\n<h1>Hello</h1>\n<p>Name: ${name}</p>\n\n<!-- out is like res.write() in Express:\n   res.write('<h1>Hello</h1>');\n   But JSP gives you 'out' for free inside <% %> */}",
			Category:    "Implicit Objects",
		},
		{
			Title:       "pageContext Object",
			Description: "Advanced object to access other implicit objects and control page flow. Rarely used - most code uses request/session directly. Only useful for framework/library code.",
			Code:        "<%\n    // Access other implicit objects (rarely needed)\n    HttpServletRequest req = (HttpServletRequest) pageContext.getRequest();\n    HttpSession sess = pageContext.getSession();\n    \n    // Forward to another page (like response.sendRedirect but different)\n    pageContext.forward(\"nextPage.jsp\");\n    \n    // Search for variable in all scopes\n    // Checks: page -> request -> session -> application\n    Object username = pageContext.findAttribute(\"username\");\n    \n    // You'll rarely use this - it's for advanced cases\n%>",
			Category:    "Implicit Objects",
		},
		{
			Title:       "config Object",
			Description: "ServletConfig instance. Access servlet init parameters from web.xml. Rarely used in JSP - more common in servlets.",
			Code:        "// Get init parameters from web.xml\nString dbUrl = config.getInitParameter(\"dbUrl\");\nString driver = config.getInitParameter(\"driver\");\n\n// Get servlet info\nString servletName = config.getServletName();\n\n// In web.xml:\n// <servlet>\n//   <servlet-name>MyJSP</servlet-name>\n//   <jsp-file>/page.jsp</jsp-file>\n//   <init-param>\n//     <param-name>dbUrl</param-name>\n//     <param-value>jdbc:mysql://localhost/mydb</param-value>\n//   </init-param>\n// </servlet>",
			Category:    "Implicit Objects",
		},
		{
			Title:       "page Object",
			Description: "Reference to current JSP page instance (this). Rarely used. Equivalent to 'this' in a servlet.",
			Code:        "// Rarely used - refers to current page instance\nObject currentPage = page;\n\n// More common: use pageContext instead\npageContext.setAttribute(\"data\", value); // page scope",
			Category:    "Implicit Objects",
		},
		{
			Title:       "exception Object",
			Description: "Only available in error pages (isErrorPage=\"true\"). Contains the exception that caused the error.",
			Code:        "<%@ page isErrorPage=\"true\" %>\n\n<h1>Error Occurred</h1>\n<p>Message: <%= exception.getMessage() %></p>\n<p>Type: <%= exception.getClass().getName() %></p>\n\n<%\n    // Stack trace\n    exception.printStackTrace();\n    \n    // Log it\n    application.log(\"Error: \" + exception.getMessage(), exception);\n%>",
			Category:    "Implicit Objects",
		},

		// Scopes
		{
			Title:       "The 4 Scopes Explained",
			Description: "JSP has 4 places to store data. Page = this file only. Request = this HTTP request. Session = this user. Application = entire app (all users). Each has different lifetime.",
			Code:        "<!-- 1. PAGE SCOPE - dies at end of THIS JSP file -->\n<%\n    pageContext.setAttribute(\"temp\", \"only here\");\n    // If you forward to another.jsp, it's gone\n%>\n\n<!-- 2. REQUEST SCOPE - lives during this HTTP request -->\n<%\n    request.setAttribute(\"data\", userData);\n    // Survives if you forward/include to another JSP\n    // Dies after response sent to browser\n%>\n\n<!-- 3. SESSION SCOPE - lives for this user's session -->\n<%\n    session.setAttribute(\"username\", \"john\");\n    session.setAttribute(\"cart\", shoppingCart);\n    // Survives across page refreshes, clicking links\n    // Dies after 30 min timeout or logout\n%>\n\n<!-- 4. APPLICATION SCOPE - shared by ALL users -->\n<%\n    application.setAttribute(\"siteConfig\", config);\n    // ALL users see this\n    // Lives until server restarts\n%>",
			Category:    "Scopes",
		},
		{
			Title:       "Request vs Session Scope (Most Important)",
			Description: "Request = data for ONE page load. Session = data that follows the user around. Most common confusion in JSP. Request dies after response, Session persists.",
			Code:        "<!-- REQUEST SCOPE - dies after this page loads -->\n<%\n    // Servlet puts user data in request\n    request.setAttribute(\"user\", userObject);\n    request.getRequestDispatcher(\"profile.jsp\").forward(request, response);\n%>\n\n<!-- profile.jsp -->\n<%\n    User user = (User) request.getAttribute(\"user\"); // works!\n%>\n<!-- User clicks link to settings.jsp -->\n<!-- settings.jsp -->\n<%\n    User user = (User) request.getAttribute(\"user\"); // NULL! Gone!\n%>\n\n<!-- SESSION SCOPE - survives across pages -->\n<%\n    // Login servlet puts username in session\n    session.setAttribute(\"username\", \"john\");\n%>\n<!-- User can navigate anywhere, username still there -->\n<!-- profile.jsp, settings.jsp, any page: -->\n<%\n    String name = (String) session.getAttribute(\"username\"); // works!\n%>",
			Category:    "Scopes",
		},
		{
			Title:       "When to Use Each Scope",
			Description: "Page = temp calculations. Request = passing data from servlet to JSP. Session = login info, shopping cart. Application = site-wide config, counters.",
			Code:        "<!-- PAGE - rarely used, just for temp variables -->\n<%\n    pageContext.setAttribute(\"tempCalc\", price * quantity);\n%>\n\n<!-- REQUEST - servlet → JSP data passing -->\n// UserServlet.java\nrequest.setAttribute(\"userList\", users);\nrequest.getRequestDispatcher(\"users.jsp\").forward(req, res);\n\n<!-- SESSION - user-specific data -->\n<%\n    // Login\n    session.setAttribute(\"userId\", user.getId());\n    session.setAttribute(\"username\", user.getName());\n    session.setAttribute(\"cart\", new ShoppingCart());\n%>\n\n<!-- APPLICATION - site-wide data -->\n<%\n    // Startup servlet sets config\n    application.setAttribute(\"dbConnections\", 100);\n    application.setAttribute(\"maxUploadSize\", 10 * 1024 * 1024);\n%>",
			Category:    "Scopes",
		},
		{
			Title:       "How EL Searches Scopes",
			Description: "When you use ${variable} in EL, it searches: page → request → session → application. First match wins. Be careful with name collisions!",
			Code:        "<%\n    request.setAttribute(\"name\", \"Request Name\");\n    session.setAttribute(\"name\", \"Session Name\");\n%>\n\n<!-- Which one prints? -->\n<p>${name}</p>  <!-- Prints \"Request Name\" (found first) -->\n\n<!-- To specify scope explicitly: -->\n<p>${requestScope.name}</p>   <!-- \"Request Name\" -->\n<p>${sessionScope.name}</p>    <!-- \"Session Name\" -->\n\n<%\n    // Only session has this variable\n    session.setAttribute(\"username\", \"john\");\n%>\n<p>${username}</p>  <!-- Works! Searches all scopes, finds in session -->\n\n<!-- Good practice: use explicit scope for clarity -->\n<p>${sessionScope.username}</p>",
			Category:    "Scopes",
		},

		// EL (Expression Language)
		{
			Title:       "Expression Language (EL) - ${ }",
			Description: "EL lets you access variables without Java code. Use ${ } instead of <%= %>. Cleaner and safer. Automatically escapes HTML to prevent XSS.",
			Code:        "<!-- Old way with scriptlets -->\n<p>Name: <%= request.getAttribute(\"name\") %></p>\n<p>Age: <%= ((User)request.getAttribute(\"user\")).getAge() %></p>\n\n<!-- New way with EL -->\n<p>Name: ${name}</p>\n<p>Age: ${user.age}</p>  <!-- Calls user.getAge() automatically -->\n\n<!-- EL can do operations -->\n<p>Total: ${price * quantity}</p>\n<p>Tax: ${price * 0.08}</p>\n<p>Is Adult: ${age >= 18}</p>\n<p>Status: ${active ? 'Active' : 'Inactive'}</p>",
			Category:    "Expression Language",
		},
		{
			Title:       "EL Implicit Objects",
			Description: "Access request params, headers, cookies, scopes. param for single value, paramValues for arrays. Different from JSP implicit objects!",
			Code:        "<!-- Access request parameters -->\n<p>Name: ${param.name}</p> <!-- single value -->\n<p>Hobbies: ${paramValues.hobby[0]}</p> <!-- array -->\n\n<!-- Access headers and cookies -->\n<p>User Agent: ${header['User-Agent']}</p>\n<p>Username Cookie: ${cookie.username.value}</p>\n\n<!-- Access specific scopes -->\n<p>Request: ${requestScope.data}</p>\n<p>Session: ${sessionScope.username}</p>\n<p>Application: ${applicationScope.config}</p>",
			Category:    "Expression Language",
		},
		{
			Title:       "EL Operators",
			Description: "Arithmetic (+, -, *, /, %), comparison (==, !=, <, >), logical (&&, ||, !), ternary (?:), empty check.",
			Code:        "<!-- Arithmetic -->\n${10 + 5}        <!-- 15 -->\n${price * 1.1}   <!-- 10% markup -->\n${total / count} <!-- average -->\n${count % 2}     <!-- even/odd check -->\n\n<!-- Comparison (also: eq, ne, lt, gt, le, ge) -->\n${age >= 18}     <!-- boolean -->\n${status == 'active' ? 'Yes' : 'No'} <!-- ternary -->\n\n<!-- Empty check -->\n${empty username} <!-- true if null or empty string -->\n${not empty cart} <!-- true if cart has items -->",
			Category:    "Expression Language",
		},
		{
			Title:       "EL Property Navigation",
			Description: "Dot notation for properties. ${user.name} calls user.getName(). ${map.key} gets map value. ${list[0]} accesses array/list.",
			Code:        "<!-- Bean properties (calls getters) -->\n${user.name}        <!-- calls user.getName() -->\n${user.address.city} <!-- calls user.getAddress().getCity() -->\n\n<!-- Map access -->\n${sessionScope.cart['item1']} <!-- map.get(\"item1\") -->\n${header['User-Agent']}       <!-- header with dash -->\n\n<!-- Array/List access -->\n${colors[0]}      <!-- first element -->\n${users[index].name} <!-- dynamic index -->",
			Category:    "Expression Language",
		},

		// JSTL (JSP Standard Tag Library)
		{
			Title:       "What is JSTL?",
			Description: "JSP Standard Tag Library. A set of tags that replace Java code in JSP. Instead of <% if %>, use <c:if>. Instead of <% for %>, use <c:forEach>. Makes JSP cleaner and easier to read.",
			Code:        "<!-- Import JSTL at top of .jsp file -->\n<%@ taglib uri=\"http://java.sun.com/jsp/jstl/core\" prefix=\"c\" %>\n\n<!-- WITHOUT JSTL (scriptlets) -->\n<%\n    if (user.getAge() >= 18) {\n%>\n    <p>Adult</p>\n<%\n    }\n%>\n\n<!-- WITH JSTL (cleaner) -->\n<c:if test=\"${user.age >= 18}\">\n    <p>Adult</p>\n</c:if>\n\n<!-- Both do the same thing!\n     JSTL just looks cleaner, less Java in HTML -->",
			Category:    "JSTL",
		},
		{
			Title:       "Why Use JSTL?",
			Description: "Separates presentation (HTML) from logic (Java). Designers can edit JSP without knowing Java. Code is cleaner, more readable. In your legacy project, you'll see both styles mixed.",
			Code:        "<!-- Problem with scriptlets: -->\n<%\n    List<User> users = (List<User>) request.getAttribute(\"users\");\n    for (int i = 0; i < users.size(); i++) {\n        User u = users.get(i);\n%>\n    <tr>\n        <td><%= u.getName() %></td>\n        <td><%= u.getEmail() %></td>\n    </tr>\n<%\n    }\n%>\n\n<!-- Same thing with JSTL: -->\n<c:forEach items=\"${users}\" var=\"u\">\n    <tr>\n        <td>${u.name}</td>\n        <td>${u.email}</td>\n    </tr>\n</c:forEach>\n\n<!-- Much cleaner! No Java syntax, just tags -->\n<!-- Both generate the exact same HTML -->",
			Category:    "JSTL",
		},
		{
			Title:       "c:if - Simple Conditionals",
			Description: "Show/hide content based on condition. IMPORTANT: c:if has NO else! For if-else, use c:choose. Test uses ${} expressions (EL syntax).",
			Code:        "<!-- Show content if condition is true -->\n<c:if test=\"${user.age >= 18}\">\n    <p>You can vote</p>\n</c:if>\n\n<!-- Common mistake: c:if has NO else clause! -->\n<!-- This won't work:\n<c:if test=\"${isAdmin}\">\n    Admin\n<c:else>  <!-- NO SUCH THING! -->\n    User\n</c:else>\n</c:if>\n-->\n\n<!-- For if-else, use c:choose (see next tidbit) -->",
			Category:    "JSTL",
		},
		{
			Title:       "c:choose - If/Else Logic",
			Description: "Like if-else-if or switch statement. <c:when> = if/else if, <c:otherwise> = else. Use this when you need if-else logic.",
			Code:        "<!-- If-else-if logic -->\n<c:choose>\n    <c:when test=\"${score >= 90}\">\n        <p>Grade: A</p>\n    </c:when>\n    <c:when test=\"${score >= 80}\">\n        <p>Grade: B</p>\n    </c:when>\n    <c:when test=\"${score >= 70}\">\n        <p>Grade: C</p>\n    </c:when>\n    <c:otherwise>\n        <p>Grade: F</p>\n    </c:otherwise>\n</c:choose>\n\n<!-- Simple if-else -->\n<c:choose>\n    <c:when test=\"${user.isAdmin}\">\n        <a href=\"admin.jsp\">Admin Panel</a>\n    </c:when>\n    <c:otherwise>\n        <p>Regular User</p>\n    </c:otherwise>\n</c:choose>",
			Category:    "JSTL",
		},
		{
			Title:       "c:forEach - Looping Over Collections",
			Description: "Loop through Lists, arrays, or any collection. The most common JSTL tag you'll see. Var is the loop variable name (like 'item' in for-each loops).",
			Code:        "<!-- Loop over a list of users -->\n<c:forEach items=\"${users}\" var=\"user\">\n    <tr>\n        <td>${user.name}</td>\n        <td>${user.email}</td>\n    </tr>\n</c:forEach>\n\n<!-- Loop with index/counter -->\n<c:forEach items=\"${products}\" var=\"p\" varStatus=\"status\">\n    <p>\n        ${status.index}: ${p.name}  <!-- index = 0,1,2... -->\n        (${status.count} of ${status.end + 1})  <!-- count = 1,2,3... -->\n    </p>\n    <c:if test=\"${status.first}\">First item!</c:if>\n    <c:if test=\"${status.last}\">Last item!</c:if>\n</c:forEach>\n\n<!-- Loop a specific number of times -->\n<c:forEach begin=\"1\" end=\"10\" var=\"i\">\n    <p>Row ${i}</p>\n</c:forEach>",
			Category:    "JSTL",
		},
		{
			Title:       "c:set and c:remove",
			Description: "Create/update variables in any scope. c:set creates, c:remove deletes.",
			Code:        "<!-- Set variable in page scope (default) -->\n<c:set var=\"name\" value=\"John\" />\n<p>${name}</p>\n\n<!-- Set in specific scope -->\n<c:set var=\"username\" value=\"admin\" scope=\"session\" />\n<c:set var=\"count\" value=\"${count + 1}\" scope=\"application\" />\n\n<!-- Set bean property -->\n<c:set target=\"${user}\" property=\"name\" value=\"Jane\" />\n\n<!-- Remove variable -->\n<c:remove var=\"name\" scope=\"session\" />",
			Category:    "JSTL",
		},
		{
			Title:       "c:url",
			Description: "Build URLs with automatic session ID encoding (if cookies disabled). Adds context path automatically.",
			Code:        "<!-- Simple URL -->\n<c:url value=\"/products.jsp\" />\n<!-- Outputs: /myapp/products.jsp (with context path) -->\n\n<!-- URL with parameters -->\n<c:url value=\"/search.jsp\" var=\"searchUrl\">\n    <c:param name=\"query\" value=\"${searchTerm}\" />\n    <c:param name=\"page\" value=\"1\" />\n</c:url>\n<a href=\"${searchUrl}\">Search</a>\n<!-- Outputs: /myapp/search.jsp?query=laptop&page=1 -->",
			Category:    "JSTL",
		},
		{
			Title:       "c:redirect and c:import",
			Description: "c:redirect sends HTTP redirect (like response.sendRedirect). c:import includes content from URL.",
			Code:        "<!-- Redirect to another page -->\n<c:if test=\"${empty username}\">\n    <c:redirect url=\"/login.jsp\" />\n</c:if>\n\n<!-- Redirect with parameters -->\n<c:redirect url=\"/error.jsp\">\n    <c:param name=\"code\" value=\"404\" />\n</c:redirect>\n\n<!-- Import/include content from URL -->\n<c:import url=\"/header.jsp\" />\n<c:import url=\"http://example.com/data.xml\" var=\"xmlData\" />",
			Category:    "JSTL",
		},
		{
			Title:       "c:out",
			Description: "Safely output text with HTML escaping. Prevents XSS attacks. Use instead of ${ } for user input.",
			Code:        "<!-- Unsafe - XSS vulnerable! -->\n<p>${userInput}</p> <!-- If userInput = \"<script>alert('XSS')</script>\" -->\n\n<!-- Safe - escapes HTML -->\n<c:out value=\"${userInput}\" />\n<!-- Outputs: &lt;script&gt;alert('XSS')&lt;/script&gt; -->\n\n<!-- Default value if null -->\n<c:out value=\"${username}\" default=\"Guest\" />\n\n<!-- Disable escaping (dangerous!) -->\n<c:out value=\"${htmlContent}\" escapeXml=\"false\" />",
			Category:    "JSTL",
		},
		{
			Title:       "fmt:formatDate and fmt:formatNumber",
			Description: "Format dates and numbers. Control patterns, locales, currencies.",
			Code:        "<!-- Format date -->\n<fmt:formatDate value=\"${order.date}\" pattern=\"yyyy-MM-dd\" />\n<fmt:formatDate value=\"${now}\" pattern=\"MM/dd/yyyy HH:mm:ss\" />\n<fmt:formatDate value=\"${date}\" type=\"date\" dateStyle=\"full\" />\n\n<!-- Format number -->\n<fmt:formatNumber value=\"${price}\" type=\"currency\" />\n<!-- $1,234.56 -->\n\n<fmt:formatNumber value=\"${percent}\" type=\"percent\" />\n<!-- 75% -->\n\n<fmt:formatNumber value=\"${pi}\" maxFractionDigits=\"2\" />\n<!-- 3.14 -->",
			Category:    "JSTL",
		},
		{
			Title:       "fn:length and fn:contains",
			Description: "JSTL functions for strings/collections. fn:length for size, fn:contains for substring check.",
			Code:        "<%@ taglib uri=\"http://java.sun.com/jsp/jstl/functions\" prefix=\"fn\" %>\n\n<!-- String length -->\n<c:if test=\"${fn:length(username) < 3}\">\n    Username too short\n</c:if>\n\n<!-- Collection size -->\n<p>Cart has ${fn:length(cart)} items</p>\n\n<!-- Substring check -->\n<c:if test=\"${fn:contains(email, '@')}\">\n    Valid email format\n</c:if>\n\n<!-- Case-insensitive contains -->\n<c:if test=\"${fn:containsIgnoreCase(message, 'error')}\">\n    Error detected\n</c:if>",
			Category:    "JSTL",
		},

		// Advanced
		{
			Title:       "Custom Tags",
			Description: "Create reusable JSP tags. Define in .tag file or TLD. Better than includes for reusable components with logic.",
			Code:        "<!-- In /WEB-INF/tags/hello.tag -->\n<%@ attribute name=\"name\" required=\"true\" %>\n<h1>Hello, ${name}!</h1>\n\n<!-- In JSP -->\n<%@ taglib prefix=\"my\" tagDir=\"/WEB-INF/tags\" %>\n<my:hello name=\"John\" />\n\n<!-- Tag with body -->\n<!-- In panel.tag -->\n<%@ attribute name=\"title\" required=\"true\" %>\n<div class=\"panel\">\n    <h2>${title}</h2>\n    <jsp:doBody/> <!-- insert body content here -->\n</div>\n\n<!-- Usage -->\n<my:panel title=\"User Info\">\n    <p>Name: ${user.name}</p>\n</my:panel>",
			Category:    "Advanced",
		},
		{
			Title:       "RequestDispatcher: forward vs include",
			Description: "forward = transfer control to another resource, stop current page. include = insert content, continue current page.",
			Code:        "<%\n// Forward - transfer control (URL stays same!)\nRequestDispatcher rd = request.getRequestDispatcher(\"/nextPage.jsp\");\nrd.forward(request, response);\n// Code after forward() DOESN'T run\n\n// Include - insert content and continue\nRequestDispatcher rd2 = request.getRequestDispatcher(\"/header.jsp\");\nrd2.include(request, response);\n// Code after include() DOES run\n%>\n\n<!-- JSTL alternatives -->\n<jsp:forward page=\"/nextPage.jsp\" />\n<jsp:include page=\"/header.jsp\" />",
			Category:    "Advanced",
		},
		{
			Title:       "Error Pages",
			Description: "Handle exceptions gracefully. Set errorPage in page directive or web.xml. Error page has access to exception object.",
			Code:        "<!-- In your JSP -->\n<%@ page errorPage=\"error.jsp\" %>\n\n<%\n    // If exception occurs, forwards to error.jsp\n    int x = 10 / 0; // ArithmeticException\n%>\n\n<!-- error.jsp -->\n<%@ page isErrorPage=\"true\" %>\n<h1>Oops! Something went wrong</h1>\n<p>Error: <%= exception.getMessage() %></p>\n<p>Type: <%= exception.getClass().getName() %></p>\n\n<!-- Or configure in web.xml -->\n<!-- <error-page>\n    <exception-type>java.lang.Exception</exception-type>\n    <location>/error.jsp</location>\n</error-page> -->",
			Category:    "Advanced",
		},
		{
			Title:       "JSP vs Servlet",
			Description: "JSP = HTML with Java (view). Servlet = Java with HTML (controller). JSP compiles to servlet. Use JSP for presentation, servlets for logic (MVC).",
			Code:        "// Servlet (Controller) - handles logic\npublic class UserServlet extends HttpServlet {\n    protected void doGet(HttpServletRequest request, \n                         HttpServletResponse response) {\n        // Get data\n        User user = userService.getUser(id);\n        \n        // Store in request\n        request.setAttribute(\"user\", user);\n        \n        // Forward to JSP\n        request.getRequestDispatcher(\"/user.jsp\")\n               .forward(request, response);\n    }\n}\n\n<!-- JSP (View) - displays data -->\n<h1>User: ${user.name}</h1>\n<p>Email: ${user.email}</p>",
			Category:    "Advanced",
		},
		{
			Title:       "Session Tracking Methods",
			Description: "4 ways to track users: Cookies (best), URL rewriting (if cookies disabled), Hidden fields, HttpSession API.",
			Code:        "// 1. Cookies (automatic with session)\nsession.setAttribute(\"user\", username);\n// Tomcat creates JSESSIONID cookie\n\n// 2. URL rewriting (if cookies disabled)\nString url = response.encodeURL(\"page.jsp\");\n// Adds ;jsessionid=ABC123 to URL\n\n// 3. Hidden form fields (manual tracking)\n<form>\n    <input type=\"hidden\" name=\"sessionId\" value=\"<%= session.getId() %>\" />\n</form>\n\n// 4. HttpSession API (what you normally use)\nsession.setAttribute(\"cart\", cart);\nCart cart = (Cart) session.getAttribute(\"cart\");",
			Category:    "Advanced",
		},
		{
			Title:       "Thread Safety in JSP",
			Description: "JSP is NOT thread-safe by default. Multiple users share ONE instance. Avoid instance variables (<%! %>). Use local variables in <% %>.",
			Code:        "<%!\n    // DANGEROUS - shared by ALL users!\n    private int counter = 0; // NOT thread-safe!\n    private String username; // All users share this!\n%>\n\n<%\n    // SAFE - each request gets its own local variable\n    int counter = 0; // OK - local to this request\n    String name = request.getParameter(\"name\"); // OK - local\n    \n    // SAFE - request/session scope\n    session.setAttribute(\"userCount\", count); // OK - per user\n    request.setAttribute(\"data\", value); // OK - per request\n%>\n\n<!-- To make JSP thread-safe for all: -->\n<%@ page isThreadSafe=\"false\" %>\n<!-- But this is slow - one request at a time! -->",
			Category:    "Advanced",
		},
		{
			Title:       "JSTL vs Scriptlets (Old vs Less Old)",
			Description: "Scriptlets = Java in HTML (VERY old, 1999). JSTL + EL = tags in HTML (less old, 2002). Both are legacy! Modern = React/Vue. But if stuck in JSP world, use JSTL.",
			Code:        "<!-- WORST - Scriptlets (1999 style) -->\n<%\n    List<User> users = (List<User>) request.getAttribute(\"users\");\n    for (User user : users) {\n%>\n    <p><%= user.getName() %></p>\n<%\n    }\n%>\n\n<!-- BETTER - JSTL (2002 style) -->\n<c:forEach items=\"${users}\" var=\"user\">\n    <p>${user.name}</p>\n</c:forEach>\n\n<!-- BEST - Modern framework (2024 style) -->\n<!-- {users.map(user => <p>{user.name}</p>)} -->",
			Category:    "Best Practices",
		},
		{
			Title:       "MVC in JSP World",
			Description: "Even old JSP tried to separate concerns! Servlet = Controller (handles request). JSP = View (displays HTML). JavaBean = Model (data). Like Express router + React component + data model.",
			Code:        "// CONTROLLER - UserServlet.java (like Express route)\npublic void doGet(HttpServletRequest req, HttpServletResponse res) {\n    // Like: app.get('/user/:id', (req, res) => {...})\n    User user = userDAO.getUser(id); // fetch data\n    req.setAttribute(\"user\", user);  // pass to view\n    req.getRequestDispatcher(\"user.jsp\").forward(req, res);\n    // Like: res.render('user', {user})\n}\n\n// VIEW - user.jsp (like React component)\n<h1>${user.name}</h1>  <!-- like {user.name} -->\n<p>Age: ${user.age}</p>",
			Category:    "Best Practices",
		},
		{
			Title:       "When You'll See JSP",
			Description: "Banks, insurance, government, old enterprise apps. If company has Java backend from 2000s-2010s, probably JSP. Spring Boot (modern Java) still CAN use JSP but most use React/Vue now.",
			Code:        "// You'll see JSP in:\n// - Legacy banking systems (security can't change easily)\n// - Government sites (slow to modernize)\n// - Old enterprise CRUD apps\n// - Maintaining 20-year-old codebases\n\n// Modern Java doesn't use JSP:\n// - Spring Boot + React (separate frontend/backend)\n// - REST APIs returning JSON (no HTML generation)\n// - Microservices (no server-side rendering)\n\n// If joining old company:\n// Day 1: \"We use JSP\" = you'll need this knowledge\n// Day 1: \"We use Spring Boot + React\" = you won't",
			Category:    "Introduction",
		},
		{
			Title:       "JSP Lifecycle",
			Description: "JSP translates to servlet, compiles once, handles requests. jspInit() once, _jspService() per request, jspDestroy() on shutdown.",
			Code:        "// 1. Translation: page.jsp -> page_jsp.java (servlet code)\n// 2. Compilation: page_jsp.java -> page_jsp.class\n// 3. Loading: Tomcat loads class\n// 4. Instantiation: Creates ONE instance (shared!)\n// 5. jspInit(): Called ONCE on first request\n// 6. _jspService(): Called EVERY request (your JSP code here)\n// 7. jspDestroy(): Called ONCE on shutdown\n\n<%!\n    public void jspInit() {\n        // Initialize resources (DB connection pool, etc)\n        System.out.println(\"JSP initialized\");\n    }\n    \n    public void jspDestroy() {\n        // Cleanup resources\n        System.out.println(\"JSP destroyed\");\n    }\n%>",
			Category:    "Advanced",
		},
	}
}
