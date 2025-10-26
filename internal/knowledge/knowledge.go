package knowledge

// KnowledgeBit represents a concept/tidbit to remember
type KnowledgeBit struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Code        string `json:"code,omitempty"`
	Category    string `json:"category"`
}

// GetKnowledgeBits returns all embedded knowledge tidbits
func GetKnowledgeBits() []KnowledgeBit {
	return []KnowledgeBit{
		// OOP Pillars
		{
			Title:       "Encapsulation",
			Description: "Hiding internal data behind methods. Private fields, public getters/setters.",
			Code:        "private double balance;\npublic double getBalance() { return balance; }",
			Category:    "OOP Pillars",
		},
		{
			Title:       "Inheritance",
			Description: "One class extends another. Child gets parent's methods/fields.",
			Code:        "public class Employee extends Person { }",
			Category:    "OOP Pillars",
		},
		{
			Title:       "Polymorphism",
			Description: "One type, many forms. Same variable type (Animal), different actual objects (Dog, Cat) = different behavior. The JVM picks which method to call at runtime based on actual object type.",
			Code:        "Animal a = new Dog();\na.makeSound(); // calls Dog's makeSound(), not Animal's\n\nAnimal b = new Cat();\nb.makeSound(); // calls Cat's makeSound()\n\n// Same variable type (Animal), different behavior at runtime\n// This is polymorphism!",
			Category:    "OOP Pillars",
		},
		{
			Title:       "Abstraction",
			Description: "Hiding complexity behind interfaces/abstract classes. Show only what matters.",
			Code:        "interface PaymentProcessor {\n    void process(double amount);\n}",
			Category:    "OOP Pillars",
		},

		// Polymorphism Deep Dive
		{
			Title:       "Overloading (Compile-time Polymorphism)",
			Description: "Same method name, different parameters. Compiler picks the method.",
			Code:        "void print(int x) { }\nvoid print(String x) { } // overloaded",
			Category:    "Polymorphism",
		},
		{
			Title:       "Overriding (Runtime Polymorphism)",
			Description: "Subclass replaces parent's method. JVM picks the method at runtime.",
			Code:        "@Override\nvoid greet() { System.out.println(\"Hey\"); }",
			Category:    "Polymorphism",
		},
		{
			Title:       "Method Overloading Rules",
			Description: "Must differ in number or type of parameters. Return type alone is not enough. Compile-time decision.",
			Code:        "void calc(int a) { }\nvoid calc(int a, int b) { } // OK\n// int calc(int x) { } // ERROR - return type not enough",
			Category:    "Polymorphism",
		},
		{
			Title:       "Method Overriding Rules",
			Description: "Must have same signature. Can't be more restrictive. Can throw narrower exceptions. Runtime decision.",
			Code:        "@Override\npublic void method() { } // was protected in parent",
			Category:    "Polymorphism",
		},
		{
			Title:       "Covariant Return Types",
			Description: "When overriding, you can return a subtype of the original return type. Java 5+ feature.",
			Code:        "class AnimalFactory {\n    Animal create() { return new Animal(); }\n}\n\nclass DogFactory extends AnimalFactory {\n    @Override\n    Dog create() { return new Dog(); } // OK! Dog is subtype of Animal\n}",
			Category:    "Polymorphism",
		},

		// Abstract Class vs Interface
		{
			Title:       "Abstract Class",
			Description: "Can have concrete methods, fields with state, single inheritance only. Use when classes share common behavior.",
			Code:        "abstract class Animal {\n    abstract void makeSound();\n    void breathe() { /* shared */ }\n}",
			Category:    "Abstraction",
		},
		{
			Title:       "Interface",
			Description: "Contract only (pure abstraction), no state, multiple inheritance allowed. Use for 'can do' relationships.",
			Code:        "interface Flyable {\n    void fly();\n}",
			Category:    "Abstraction",
		},
		{
			Title:       "Default Methods (Java 8+)",
			Description: "Interfaces can have method implementations now. Allows adding methods without breaking existing implementations.",
			Code:        "interface MyInterface {\n    default void log(String msg) {\n        System.out.println(msg);\n    }\n}",
			Category:    "Abstraction",
		},
		{
			Title:       "Functional Interface",
			Description: "Interface with exactly ONE abstract method. Can be used with lambdas. Annotate with @FunctionalInterface.",
			Code:        "@FunctionalInterface\ninterface Calculator {\n    int calculate(int a, int b);\n}",
			Category:    "Abstraction",
		},

		// Composition vs Inheritance
		{
			Title:       "Inheritance (is-a)",
			Description: "Dog IS-A Animal. Use for true subtype relationships.",
			Code:        "class Dog extends Animal { }",
			Category:    "Design Patterns",
		},
		{
			Title:       "Composition (has-a)",
			Description: "Car HAS-A Engine. Prefer composition over inheritance.",
			Code:        "class Car {\n    private Engine engine;\n}",
			Category:    "Design Patterns",
		},
		{
			Title:       "Liskov Substitution Principle",
			Description: "Objects of superclass should be replaceable with objects of subclass without breaking the application. Key to polymorphism.",
			Code:        "Animal a = new Dog(); // Dog can substitute Animal",
			Category:    "Design Patterns",
		},
		{
			Title:       "Dependency Injection",
			Description: "Pass dependencies to a class instead of creating them inside. Makes testing easier, reduces coupling.",
			Code:        "class UserService {\n    private Database db;\n    UserService(Database db) { this.db = db; } // injected\n}",
			Category:    "Design Patterns",
		},
		{
			Title:       "Singleton Pattern",
			Description: "Ensures only one instance of a class exists. Private constructor, static getInstance() method.",
			Code:        "class Database {\n    private static Database instance;\n    \n    private Database() {} // private constructor\n    \n    public static Database getInstance() {\n        if (instance == null) {\n            instance = new Database();\n        }\n        return instance;\n    }\n}",
			Category:    "Design Patterns",
		},
		{
			Title:       "Factory Pattern",
			Description: "Create objects without specifying exact class. Useful when creation logic is complex or needs to be centralized.",
			Code:        "interface Animal { void speak(); }\nclass Dog implements Animal { void speak() {...} }\nclass Cat implements Animal { void speak() {...} }\n\nclass AnimalFactory {\n    public static Animal create(String type) {\n        if (type.equals(\"dog\")) return new Dog();\n        if (type.equals(\"cat\")) return new Cat();\n        throw new IllegalArgumentException();\n    }\n}",
			Category:    "Design Patterns",
		},
		{
			Title:       "Builder Pattern",
			Description: "Build complex objects step by step. Common for objects with many optional parameters.",
			Code:        "class User {\n    String name;\n    int age;\n    String email;\n    \n    private User(Builder b) {\n        this.name = b.name;\n        this.age = b.age;\n        this.email = b.email;\n    }\n    \n    static class Builder {\n        String name, email;\n        int age;\n        Builder name(String n) { name = n; return this; }\n        Builder age(int a) { age = a; return this; }\n        User build() { return new User(this); }\n    }\n}\n\nUser user = new User.Builder().name(\"John\").age(30).build();",
			Category:    "Design Patterns",
		},

		// Modifiers
		{
			Title:       "Static",
			Description: "Belongs to class, not instance. Shared across all objects.",
			Code:        "static int count;\nstatic void doThing() { }",
			Category:    "Modifiers",
		},
		{
			Title:       "Final",
			Description: "final variable = can't reassign, final method = can't override, final class = can't inherit from.",
			Code:        "final int MAX = 100;\nfinal void method() { }\nfinal class Util { }",
			Category:    "Modifiers",
		},
		{
			Title:       "Access Modifiers",
			Description: "private (only this class), default (package), protected (package + subclasses), public (everyone).",
			Code:        "private int x;\nint y; // default/package\nprotected int z;\npublic int w;",
			Category:    "Modifiers",
		},
		{
			Title:       "Static vs Instance Initialization Blocks",
			Description: "static { } runs once when class loads. { } runs before each constructor call.",
			Code:        "class Foo {\n    static { System.out.println(\"Class loaded\"); }\n    { System.out.println(\"Instance created\"); }\n}",
			Category:    "Modifiers",
		},
		{
			Title:       "Transient Keyword",
			Description: "Marks fields to skip during serialization (saving object to file/network). Use for passwords, temp data, things you don't want saved. Only matters with Serializable.",
			Code:        "class User implements Serializable {\n    private String name;      // WILL be saved\n    private transient String password; // will NOT be saved\n}\n// When you save/send User, password stays null\n// Useful for sensitive data or calculated fields",
			Category:    "Modifiers",
		},
		{
			Title:       "Volatile Keyword",
			Description: "For multi-threaded code. Tells JVM: don't cache this variable, always read from main memory. Ensures all threads see latest value. Use for flags that multiple threads check.",
			Code:        "class Worker {\n    private volatile boolean running = true; // all threads see updates\n    \n    void stop() { running = false; } // thread 1\n    \n    void run() {\n        while (running) { work(); } // thread 2 sees change immediately\n    }\n}\n// Without volatile, thread 2 might cache 'running' and never see it change!",
			Category:    "Modifiers",
		},
		{
			Title:       "Synchronized Keyword",
			Description: "Lock a method/block so only ONE thread can run it at a time. Others wait. Prevents two threads from corrupting shared data (race condition). Makes code thread-safe but slower.",
			Code:        "class Counter {\n    private int count = 0;\n    \n    // Without synchronized - BROKEN\n    void increment() { count++; } // two threads = data corruption\n    \n    // With synchronized - SAFE\n    synchronized void increment() {\n        count++; // only one thread at a time\n    }\n}",
			Category:    "Modifiers",
		},

		// Exceptions
		{
			Title:       "Checked Exceptions",
			Description: "Compiler forces you to handle with try-catch or declare with throws. For recoverable errors. (IOException, SQLException)",
			Code:        "public void readFile() throws IOException {\n    FileReader fr = new FileReader(\"file.txt\");\n}",
			Category:    "Exceptions",
		},
		{
			Title:       "Unchecked Exceptions (Runtime Exceptions)",
			Description: "Compiler doesn't care. Subclass of RuntimeException. For programming bugs. (NullPointerException, IllegalArgumentException)",
			Code:        "String s = null;\ns.length(); // NPE at runtime",
			Category:    "Exceptions",
		},
		{
			Title:       "Try-with-resources",
			Description: "Automatically closes resources (files, connections). No need for finally block. Java 7+ feature.",
			Code:        "try (FileReader fr = new FileReader(\"file.txt\")) {\n    // use fr\n} // automatically closed",
			Category:    "Exceptions",
		},
		{
			Title:       "Multi-catch (Java 7+)",
			Description: "Catch multiple exception types in one catch block.",
			Code:        "try {\n    // code\n} catch (IOException | SQLException e) {\n    // handle both\n}",
			Category:    "Exceptions",
		},
		{
			Title:       "Exception Hierarchy",
			Description: "Throwable is the TOP parent of everything throwable. It has 2 children: Error (JVM crashes, don't catch) and Exception. Exception has 2 types: RuntimeException (unchecked) and everything else (checked).",
			Code:        "// Hierarchy:\n// Throwable (top parent)\n//   ├── Error (don't catch these - JVM problems)\n//   │     └── OutOfMemoryError, StackOverflowError\n//   └── Exception (catch these)\n//         ├── RuntimeException (unchecked)\n//         │     └── NullPointerException, IllegalArgumentException\n//         └── Others (checked)\n//               └── IOException, SQLException\n\n// You can catch/throw anything that extends Throwable\n// But you ONLY catch Exceptions, not Errors",
			Category:    "Exceptions",
		},

		// Collections
		{
			Title:       "ArrayList vs LinkedList",
			Description: "ArrayList: fast random access, slow insert/delete. LinkedList: slow access, fast insert/delete at ends.",
			Code:        "List<String> arr = new ArrayList<>(); // O(1) get\nList<String> link = new LinkedList<>(); // O(n) get",
			Category:    "Collections",
		},
		{
			Title:       "HashSet vs TreeSet vs LinkedHashSet",
			Description: "HashSet: no order, O(1) ops. TreeSet: sorted, O(log n) ops. LinkedHashSet: insertion order, O(1) ops.",
			Code:        "Set<String> hash = new HashSet<>(); // fastest\nSet<String> tree = new TreeSet<>(); // sorted\nSet<String> linked = new LinkedHashSet<>(); // ordered",
			Category:    "Collections",
		},
		{
			Title:       "HashMap vs TreeMap vs LinkedHashMap",
			Description: "HashMap: no order, O(1). TreeMap: sorted by keys, O(log n). LinkedHashMap: insertion order, O(1).",
			Code:        "Map<String, Integer> hash = new HashMap<>();\nMap<String, Integer> tree = new TreeMap<>();\nMap<String, Integer> linked = new LinkedHashMap<>();",
			Category:    "Collections",
		},
		{
			Title:       "Comparable vs Comparator",
			Description: "Comparable = ONE natural sort order built INTO the class (implement compareTo). Comparator = MANY custom sort orders OUTSIDE the class (pass to sort method). Use Comparable for default sort, Comparator for custom sorts.",
			Code:        "// Comparable - ONE natural order inside the class\nclass Person implements Comparable<Person> {\n    int age;\n    public int compareTo(Person p) { return this.age - p.age; }\n}\nCollections.sort(people); // uses compareTo\n\n// Comparator - MULTIPLE custom orders outside the class\nComparator<Person> byName = (p1, p2) -> p1.name.compareTo(p2.name);\nCollections.sort(people, byName);",
			Category:    "Collections",
		},
		{
			Title:       "equals() and hashCode() Contract",
			Description: "The rule: if equals() returns true, hashCode() MUST return same number for both objects. If you override equals(), you MUST also override hashCode(). If you don't, HashMap/HashSet will break.",
			Code:        "class Person {\n    String name;\n    int age;\n    \n    @Override\n    public boolean equals(Object o) {\n        Person p = (Person) o;\n        return name.equals(p.name) && age == p.age;\n    }\n    \n    @Override\n    public int hashCode() {\n        return Objects.hash(name, age); // MUST override both!\n    }\n}\n// If you only override equals(), HashMap won't find your objects!",
			Category:    "Collections",
		},
		{
			Title:       "Fail-Fast vs Fail-Safe Iterators",
			Description: "Fail-fast: throw ConcurrentModificationException if collection modified during iteration. Fail-safe: don't throw, work on copy.",
			Code:        "// Fail-fast (ArrayList, HashMap)\nList<String> list = new ArrayList<>();\nfor (String s : list) {\n    list.add(\"x\"); // ConcurrentModificationException!\n}\n\n// Fail-safe (CopyOnWriteArrayList)\nList<String> safe = new CopyOnWriteArrayList<>();\nfor (String s : safe) {\n    safe.add(\"x\"); // OK, but iterates over snapshot\n}",
			Category:    "Collections",
		},

		// Generics
		{
			Title:       "Generics Basics",
			Description: "Type parameters for compile-time type safety. Prevent ClassCastException.",
			Code:        "List<String> list = new ArrayList<>();\nlist.add(\"hello\");\n// list.add(123); // compile error",
			Category:    "Generics",
		},
		{
			Title:       "Bounded Type Parameters",
			Description: "Restrict what types a generic can accept. <T extends Number> means T must be Number or any subclass (Integer, Double, etc). Lets you call Number methods on T.",
			Code:        "class Calculator<T extends Number> {\n    double add(T a, T b) {\n        return a.doubleValue() + b.doubleValue(); // can call Number methods!\n    }\n}\nCalculator<Integer> calc = new Calculator<>(); // OK\n// Calculator<String> calc = new Calculator<>(); // ERROR - String not a Number",
			Category:    "Generics",
		},
		{
			Title:       "Wildcards: ? extends vs ? super",
			Description: "Confusing generics wildcards. ? extends T = READ only (producer). ? super T = WRITE only (consumer). Remember: PECS = Producer Extends, Consumer Super.",
			Code:        "// ? extends - PRODUCER (you read from it)\n// Accepts List<Integer>, List<Double>, etc\nvoid processNumbers(List<? extends Number> nums) {\n    Number n = nums.get(0); // OK - can read\n    // nums.add(5); // ERROR - can't add (what if it's List<Double>?)\n}\n\n// ? super - CONSUMER (you write to it)\n// Accepts List<Integer>, List<Number>, List<Object>, etc\nvoid addInts(List<? super Integer> list) {\n    list.add(5); // OK - can write Integer\n    // Integer x = list.get(0); // ERROR - returns Object\n}",
			Category:    "Generics",
		},
		{
			Title:       "Type Erasure",
			Description: "Generics only exist at COMPILE time. At RUNTIME, Java erases all generic type info. List<String> and List<Integer> become just 'List'. Can't check generic types at runtime.",
			Code:        "List<String> strings = new ArrayList<>();\nList<Integer> ints = new ArrayList<>();\n\n// At runtime, both are just ArrayList - type info erased!\nstrings.getClass() == ints.getClass() // true\n\n// Can't check generic type at runtime\nif (list instanceof List<String>) { } // COMPILE ERROR\nif (list instanceof List) { } // OK - can only check raw type",
			Category:    "Generics",
		},

		// Concurrency
		{
			Title:       "Thread vs Runnable",
			Description: "Both create parallel tasks. Thread = extend class (locks you in, can't extend anything else). Runnable = implement interface (better, can still extend other classes). Always prefer Runnable.",
			Code:        "// Bad - can't extend anything else now\nclass MyThread extends Thread {\n    public void run() { System.out.println(\"Running\"); }\n}\nnew MyThread().start();\n\n// Good - can still extend other classes\nclass MyTask implements Runnable {\n    public void run() { System.out.println(\"Running\"); }\n}\nnew Thread(new MyTask()).start();",
			Category:    "Concurrency",
		},
		{
			Title:       "ExecutorService",
			Description: "Manages a pool of worker threads that can run tasks. Instead of creating 100 threads (expensive), reuse 5 threads for 100 tasks. Used in real apps for background work.",
			Code:        "// Create pool of 5 reusable threads\nExecutorService executor = Executors.newFixedThreadPool(5);\n\n// Submit 100 tasks - they share the 5 threads\nfor (int i = 0; i < 100; i++) {\n    executor.submit(() -> processTask());\n}\n\nexecutor.shutdown(); // cleanup when done",
			Category:    "Concurrency",
		},
		{
			Title:       "Callable vs Runnable",
			Description: "Both run tasks in background. Runnable = fire and forget (void run()). Callable = get result back (T call()). Use Callable when you need the result of background work.",
			Code:        "// Runnable - no return value\nRunnable task1 = () -> { System.out.println(\"Done\"); };\n\n// Callable - returns value\nCallable<Integer> task2 = () -> {\n    // do expensive calculation\n    return 42;\n};\nFuture<Integer> result = executor.submit(task2);\nint value = result.get(); // blocks until ready",
			Category:    "Concurrency",
		},
		{
			Title:       "Future and CompletableFuture",
			Description: "Future = placeholder for result that's being computed in background. CompletableFuture = Future but you can chain operations (like Promise in JS).",
			Code:        "// Future - just wait for result\nFuture<String> future = executor.submit(() -> fetchData());\nString data = future.get(); // blocks until done\n\n// CompletableFuture - chain operations\nCompletableFuture.supplyAsync(() -> fetchUser())\n    .thenApply(user -> user.getName())\n    .thenAccept(name -> System.out.println(name));",
			Category:    "Concurrency",
		},
		{
			Title:       "CountDownLatch vs CyclicBarrier",
			Description: "CountDownLatch = wait for N things to complete (use once). CyclicBarrier = N threads wait for each other (reuse many times). Advanced concurrency tools for coordination.",
			Code:        "// CountDownLatch - wait for 3 workers to finish (one-time)\nCountDownLatch latch = new CountDownLatch(3);\nfor (int i = 0; i < 3; i++) {\n    executor.submit(() -> {\n        doWork();\n        latch.countDown(); // \"I'm done\"\n    });\n}\nlatch.await(); // wait for all 3\n\n// CyclicBarrier - 3 threads sync up (reusable)\nCyclicBarrier barrier = new CyclicBarrier(3);\nbarrier.await(); // all 3 wait here until everyone arrives",
			Category:    "Concurrency",
		},

		// Memory & JVM
		{
			Title:       "Stack vs Heap Memory",
			Description: "Stack = method calls and local primitives (fast, small, auto-cleaned). Heap = all objects created with 'new' (slower, large, garbage collected). Every 'new' goes on heap.",
			Code:        "void method() {\n    int x = 5; // x on stack (primitive)\n    String s = new String(\"hi\"); // s reference on stack\n                                  // \"hi\" object on heap\n} // x and s removed from stack, \"hi\" waits for GC",
			Category:    "Memory & JVM",
		},
		{
			Title:       "String Pool",
			Description: "JVM reuses string literals to save memory. \"hello\" is stored once, all references point to same object. new String() bypasses pool and creates duplicate.",
			Code:        "String s1 = \"hello\"; // goes in pool\nString s2 = \"hello\"; // reuses same object from pool\nSystem.out.println(s1 == s2); // true (same object!)\n\nString s3 = new String(\"hello\"); // creates NEW object\nSystem.out.println(s1 == s3); // false (different objects)",
			Category:    "Memory & JVM",
		},
		{
			Title:       "Garbage Collection Basics",
			Description: "JVM automatically deletes objects you're not using anymore. When no variables point to an object, it's 'garbage'. GC finds and deletes it. You can't force when, just suggest with System.gc().",
			Code:        "void method() {\n    String s = new String(\"hi\");\n    // s exists, object alive\n} // s goes out of scope, no references to \"hi\"\n  // object eligible for GC\n\nString s = new String(\"hi\");\ns = null; // also makes object eligible",
			Category:    "Memory & JVM",
		},
		{
			Title:       "Strong vs Weak vs Soft vs Phantom References",
			Description: "Control how GC treats objects. Strong (normal) = never GC'd while ref exists. Weak = GC can delete anytime. Soft = GC deletes only if memory is low. Advanced GC tuning concept.",
			Code:        "// Strong - normal reference, never GC'd\nString strong = new String(\"hi\");\n\n// Weak - can be GC'd even if ref exists (used in caches)\nWeakReference<String> weak = new WeakReference<>(strong);\nstrong = null;\n// GC can now collect the object\nString retrieved = weak.get(); // might be null now",
			Category:    "Memory & JVM",
		},

		// Streams & Lambdas (Java 8+)
		{
			Title:       "Lambda Expressions",
			Description: "Shorthand for anonymous functions. Instead of writing a whole anonymous class, just write the logic. Works only with interfaces that have ONE method (functional interfaces).",
			Code:        "// Old way - verbose\nlist.forEach(new Consumer<String>() {\n    public void accept(String s) {\n        System.out.println(s);\n    }\n});\n\n// Lambda - same thing, shorter\nlist.forEach(s -> System.out.println(s));\n// or even shorter\nlist.forEach(System.out::println);",
			Category:    "Streams & Lambdas",
		},
		{
			Title:       "Method References (::)",
			Description: "Even shorter than lambdas. When lambda ONLY calls one method, replace with ::. The :: means 'call this method on each element'. NOT the same as calling the method directly!",
			Code:        "// String::toUpperCase means: for each string, call its toUpperCase() method\n// It's shorthand for: s -> s.toUpperCase()\n\nlist.stream().map(s -> s.toUpperCase()); // lambda\nlist.stream().map(String::toUpperCase);   // method reference - SAME THING\n\n// System.out::println means: for each element, call System.out.println(element)\nlist.forEach(s -> System.out.println(s)); // lambda\nlist.forEach(System.out::println);        // method reference - SAME THING\n\n// Person::new means: for each element, call new Person(element)\nlist.stream().map(s -> new Person(s));  // lambda\nlist.stream().map(Person::new);          // constructor reference\n\n// NOT String.toUpperCase() - that doesn't exist! It's instance method.",
			Category:    "Streams & Lambdas",
		},
		{
			Title:       "Four Types of Method References",
			Description: "There are 4 ways to use :: shorthand. Type 3 is the confusing one (ClassName::instanceMethod) - it means call that method ON EACH object, not call it on the class.",
			Code:        "// 1. Static method - ClassName::staticMethod\nlist.stream().map(Integer::parseInt);  // same as: s -> Integer.parseInt(s)\n\n// 2. Instance method of specific object - object::method\nString prefix = \"Hello \";\nlist.stream().map(prefix::concat);  // same as: s -> prefix.concat(s)\n\n// 3. Instance method on arbitrary object - ClassName::instanceMethod\n// THE CONFUSING ONE! Calls method ON EACH element\nlist.stream().map(String::toUpperCase);  // same as: s -> s.toUpperCase()\nlist.stream().map(String::length);       // same as: s -> s.length()\n// NOT String.toUpperCase() - that doesn't exist!\n\n// 4. Constructor - ClassName::new\nlist.stream().map(Person::new);  // same as: s -> new Person(s)",
			Category:    "Streams & Lambdas",
		},
		{
			Title:       "Stream vs Collection",
			Description: "Collection = stores data (List, Set, Map). Stream = pipeline of operations on data (filter, map, etc). Stream doesn't store anything, just processes. Use streams to transform collections without loops.",
			Code:        "List<String> names = Arrays.asList(\"john\", \"jane\", \"bob\");\n\n// Old way - loop\nList<String> result = new ArrayList<>();\nfor (String n : names) {\n    if (n.startsWith(\"j\")) result.add(n.toUpperCase());\n}\n\n// Stream way - cleaner\nList<String> result = names.stream()\n    .filter(n -> n.startsWith(\"j\"))\n    .map(String::toUpperCase)\n    .collect(Collectors.toList());",
			Category:    "Streams & Lambdas",
		},
		{
			Title:       "Intermediate vs Terminal Operations",
			Description: "Intermediate ops are LAZY - they don't run until terminal op. Intermediate = return Stream (filter, map). Terminal = return result and END stream (collect, forEach). Nothing happens until terminal op!",
			Code:        "List<String> list = Arrays.asList(\"a\", \"bb\", \"ccc\");\n\n// This does NOTHING - no terminal operation\nlist.stream()\n    .filter(s -> s.length() > 1)\n    .map(String::toUpperCase); // nothing printed, nothing happens\n\n// This RUNS - has terminal operation\nlist.stream()\n    .filter(s -> s.length() > 1) // intermediate\n    .map(String::toUpperCase)    // intermediate\n    .collect(Collectors.toList()); // TERMINAL - now it runs!",
			Category:    "Streams & Lambdas",
		},
		{
			Title:       "Optional",
			Description: "Wrapper that says 'this might be null'. Instead of returning null and getting NPE, return Optional. Forces caller to handle missing value case. Prevents NullPointerException.",
			Code:        "// Old way - returns null, causes NPE\nString findUser(int id) {\n    return null; // oops\n}\nString name = findUser(1).toUpperCase(); // NPE!\n\n// Optional way - explicit about maybe-null\nOptional<String> findUser(int id) {\n    return Optional.empty();\n}\nString name = findUser(1).orElse(\"Unknown\"); // safe!\nfindUser(1).ifPresent(n -> System.out.println(n)); // only runs if present",
			Category:    "Streams & Lambdas",
		},

		// Miscellaneous
		{
			Title:       "Mutable vs Immutable",
			Description: "Mutable: can change after creation (ArrayList, StringBuilder). Immutable: can't change (String, Integer).",
			Code:        "String s = \"hello\"; // immutable\ns.toUpperCase(); // creates NEW string, s unchanged\ns = s.toUpperCase(); // must reassign",
			Category:    "Core Concepts",
		},
		{
			Title:       "== vs equals()",
			Description: "== compares references (same object?). equals() compares values (equivalent content?). Override equals() for custom comparison.",
			Code:        "String s1 = new String(\"hi\");\nString s2 = new String(\"hi\");\ns1 == s2 // false\ns1.equals(s2) // true",
			Category:    "Core Concepts",
		},
		{
			Title:       "Autoboxing and Unboxing",
			Description: "Autoboxing: automatic conversion from primitive to wrapper (int -> Integer). Unboxing: wrapper to primitive.",
			Code:        "Integer i = 5; // autoboxing\nint x = i; // unboxing\nList<Integer> list = new ArrayList<>(); // can't use List<int>",
			Category:    "Core Concepts",
		},
		{
			Title:       "Pass by Value",
			Description: "Java is ALWAYS pass-by-value. For objects, the VALUE of the reference is passed (not the object itself).",
			Code:        "void changePrimitive(int x) { x = 10; } // won't change original\n\nvoid modifyList(List<String> list) {\n    list.add(\"hi\"); // DOES change original (modifying object)\n}\n\nvoid reassignList(List<String> list) {\n    list = new ArrayList<>(); // DOESN'T change original (reassigning reference)\n}",
			Category:    "Core Concepts",
		},
		{
			Title:       "Diamond Problem (Multiple Inheritance)",
			Description: "Why Java doesn't allow multiple class inheritance. Interfaces solve this (Java 8+ default methods can still cause issues).",
			Code:        "class C extends A, B { } // NOT allowed\nclass C implements IA, IB { } // OK",
			Category:    "Core Concepts",
		},
		{
			Title:       "Serialization",
			Description: "Convert object to byte stream for storage/network. Class must implement Serializable. Use transient for fields to exclude.",
			Code:        "class User implements Serializable {\n    private static final long serialVersionUID = 1L;\n    private String name;\n}",
			Category:    "Core Concepts",
		},
		{
			Title:       "Reflection",
			Description: "Inspect and modify classes, methods, fields at runtime. Used by frameworks (Spring, Jackson). Slow, breaks encapsulation.",
			Code:        "// Get class info\nClass<?> clazz = person.getClass();\n\n// Call method by name\nMethod method = clazz.getMethod(\"setName\", String.class);\nmethod.invoke(person, \"John\");\n\n// Access private field\nField field = clazz.getDeclaredField(\"age\");\nfield.setAccessible(true); // bypass private\nfield.set(person, 30);",
			Category:    "Advanced",
		},
		{
			Title:       "Annotations",
			Description: "Metadata attached to code. @Override, @Deprecated, @SuppressWarnings are built-in. Can create custom annotations.",
			Code:        "@Override\npublic void method() { }\n@Deprecated(since=\"2.0\")\npublic void oldMethod() { }",
			Category:    "Advanced",
		},
		{
			Title:       "Varargs",
			Description: "Variable number of arguments. Must be last parameter. Internally an array.",
			Code:        "void print(String... args) {\n    for (String s : args) { ... }\n}\nprint(\"a\", \"b\", \"c\");",
			Category:    "Core Concepts",
		},
		{
			Title:       "String vs StringBuilder vs StringBuffer",
			Description: "String is IMMUTABLE (can't change). Every concat creates new object = SLOW. StringBuilder is MUTABLE (can change) = FAST. StringBuffer = StringBuilder but thread-safe (slower).",
			Code:        "// BAD - creates 1000 new String objects!\nString s = \"\";\nfor (int i = 0; i < 1000; i++) {\n    s += i; // creates new String each time\n}\n\n// GOOD - modifies same object\nStringBuilder sb = new StringBuilder();\nfor (int i = 0; i < 1000; i++) {\n    sb.append(i); // modifies existing\n}\nString result = sb.toString();",
			Category:    "Core Concepts",
		},
		{
			Title:       "Deep Copy vs Shallow Copy",
			Description: "Shallow copy = copy references (both point to same objects). Deep copy = copy actual objects (independent). clone() does shallow by default.",
			Code:        "class Person {\n    String name;\n    Address address; // object reference\n}\n\n// Shallow copy - both share same Address object\nPerson p2 = p1.clone(); // default clone()\np2.address.city = \"NYC\"; // CHANGES p1's address too!\n\n// Deep copy - independent Address objects\nPerson p2 = new Person();\np2.name = p1.name;\np2.address = new Address(p1.address); // new object\np2.address.city = \"NYC\"; // doesn't affect p1",
			Category:    "Core Concepts",
		},
		{
			Title:       "this vs super",
			Description: "this = refers to current object. super = refers to parent class. Used to disambiguate or call parent constructors/methods.",
			Code:        "class Parent {\n    int x = 10;\n    Parent() { System.out.println(\"Parent\"); }\n    void show() { System.out.println(\"Parent\"); }\n}\n\nclass Child extends Parent {\n    int x = 20;\n    Child() {\n        super(); // call parent constructor (must be first line)\n        this.x = 30; // this object's x\n    }\n    void show() {\n        System.out.println(this.x);  // 30 (child's x)\n        System.out.println(super.x); // 10 (parent's x)\n        super.show(); // call parent's show()\n    }\n}",
			Category:    "Core Concepts",
		},
		{
			Title:       "Static vs Dynamic Binding",
			Description: "Static binding = method call resolved at COMPILE time (private, static, final methods). Dynamic binding = resolved at RUNTIME (overridden methods). Dynamic = polymorphism.",
			Code:        "class Animal {\n    static void staticMethod() { System.out.println(\"Animal static\"); }\n    void instanceMethod() { System.out.println(\"Animal instance\"); }\n}\n\nclass Dog extends Animal {\n    static void staticMethod() { System.out.println(\"Dog static\"); }\n    @Override void instanceMethod() { System.out.println(\"Dog instance\"); }\n}\n\nAnimal a = new Dog();\na.staticMethod();   // \"Animal static\" - static binding (compile time)\na.instanceMethod(); // \"Dog instance\" - dynamic binding (runtime)",
			Category:    "Polymorphism",
		},
		{
			Title:       "Marker Interfaces",
			Description: "Empty interfaces with no methods. Just a 'marker' to tell JVM something about the class. Serializable, Cloneable, Remote are examples.",
			Code:        "// Marker interface - no methods!\ninterface Serializable { } \n\nclass User implements Serializable {\n    // Now JVM knows this class can be serialized\n}\n\n// Check at runtime\nif (obj instanceof Serializable) {\n    // can serialize this object\n}",
			Category:    "Core Concepts",
		},
		{
			Title:       "instanceof Operator",
			Description: "Checks if object is instance of a class/interface. Returns true if object is of that type or any subtype. Returns false if null.",
			Code:        "Object obj = \"Hello\";\n\nif (obj instanceof String) {\n    String s = (String) obj; // safe to cast\n    System.out.println(s.length());\n}\n\nAnimal a = new Dog();\na instanceof Dog    // true\na instanceof Animal // true\na instanceof Cat    // false\n\nObject nullObj = null;\nnullObj instanceof String // false (null is not instance of anything)",
			Category:    "Core Concepts",
		},
		{
			Title:       "Single Responsibility Principle (SOLID)",
			Description: "A class should have ONE reason to change. Don't mix multiple responsibilities in one class. Makes code easier to maintain and test.",
			Code:        "// BAD - UserManager does too much\nclass UserManager {\n    void saveUser() { /* DB logic */ }\n    void sendEmail() { /* email logic */ }\n    void generateReport() { /* report logic */ }\n}\n\n// GOOD - each class has one responsibility\nclass UserRepository { void save() { /* DB only */ } }\nclass EmailService { void send() { /* email only */ } }\nclass ReportGenerator { void generate() { /* report only */ } }",
			Category:    "Design Patterns",
		},
		{
			Title:       "Open/Closed Principle (SOLID)",
			Description: "Classes should be OPEN for extension but CLOSED for modification. Add new features by extending, not editing existing code. Use inheritance/interfaces.",
			Code:        "// BAD - have to modify Shape class for each new shape\nclass Shape {\n    String type;\n    double calculateArea() {\n        if (type.equals(\"circle\")) { /* circle logic */ }\n        else if (type.equals(\"square\")) { /* square logic */ }\n    }\n}\n\n// GOOD - extend, don't modify\nabstract class Shape { abstract double calculateArea(); }\nclass Circle extends Shape { double calculateArea() { /* circle */ } }\nclass Square extends Shape { double calculateArea() { /* square */ } }\n// Add new shape without modifying existing code!",
			Category:    "Design Patterns",
		},
		{
			Title:       "Primitives: 8 Types and Sizes",
			Description: "byte (8 bits), short (16), int (32), long (64), float (32), double (64), char (16), boolean (1 bit). Default numeric type is int. Default decimal is double.",
			Code:        "byte b = 127;    // -128 to 127\nshort s = 32000; // -32,768 to 32,767\nint i = 100;     // -2B to 2B (most common)\nlong l = 100L;   // -9 quintillion to 9 quintillion (need L suffix)\n\nfloat f = 3.14f;  // 32-bit decimal (need f suffix)\ndouble d = 3.14;  // 64-bit decimal (default for decimals)\n\nchar c = 'A';     // single character, 16-bit Unicode\nboolean flag = true; // true or false\n\nint x = 10;  // literal defaults to int\nlong y = 10; // OK, auto-widens to long",
			Category:    "Core Concepts",
		},
		{
			Title:       "Break and Continue with Labels",
			Description: "Normal break/continue only affects innermost loop. Use labels to break/continue outer loops. Rare but sometimes necessary.",
			Code:        "// Problem: how to break out of BOTH loops?\nouter: for (int i = 0; i < 10; i++) {\n    for (int j = 0; j < 10; j++) {\n        if (found) {\n            break outer; // breaks out of BOTH loops\n        }\n    }\n}\n\n// Continue with label\nouter: for (int i = 0; i < 10; i++) {\n    for (int j = 0; j < 10; j++) {\n        if (skip) {\n            continue outer; // skip to next iteration of outer loop\n        }\n    }\n}",
			Category:    "Core Concepts",
		},
	}
}
