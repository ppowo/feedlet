package knowledge

func getJavaBits() []KnowledgeBit {
	return []KnowledgeBit{
		// OOP Pillars
		{"What is encapsulation and why is it important in OOP?", "Java encapsulation"},
		{"Explain inheritance in Java with an example", "Java inheritance"},
		{"What is polymorphism? Explain with a real-world example", "Java polymorphism"},
		{"What is abstraction and how does it differ from encapsulation?", "Java abstraction OOP"},

		// Polymorphism
		{"What is the difference between method overloading and overriding?", "Java overloading vs overriding"},
		{"What is compile-time polymorphism vs runtime polymorphism?", "Java compile time runtime polymorphism"},
		{"What are covariant return types in Java?", "Java covariant return types"},
		{"What are the rules for method overriding in Java?", "Java method overriding rules"},

		// Abstract Class vs Interface
		{"When would you use an abstract class vs an interface?", "Java abstract class vs interface"},
		{"What are default methods in interfaces? Why were they added?", "Java interface default methods"},
		{"What is a functional interface and how is it used with lambdas?", "Java functional interface lambda"},

		// Design Patterns
		{"Explain composition over inheritance with an example", "Java composition vs inheritance"},
		{"What is the Liskov Substitution Principle?", "Liskov Substitution Principle Java"},
		{"Explain dependency injection and its benefits", "dependency injection Java"},
		{"How would you implement the Singleton pattern in Java?", "Java singleton pattern"},
		{"When would you use the Factory pattern?", "Java factory pattern"},
		{"Explain the Builder pattern with an example", "Java builder pattern"},

		// Modifiers
		{"What is the difference between static and instance variables?", "Java static vs instance"},
		{"What does the final keyword mean for variables, methods, and classes?", "Java final keyword"},
		{"Explain Java access modifiers: public, private, protected, default", "Java access modifiers"},
		{"What is the transient keyword used for?", "Java transient keyword serialization"},
		{"When would you use the volatile keyword?", "Java volatile keyword threads"},
		{"What does synchronized do and when would you use it?", "Java synchronized keyword"},

		// Exceptions
		{"What is the difference between checked and unchecked exceptions?", "Java checked unchecked exceptions"},
		{"Explain try-with-resources and its benefits", "Java try with resources"},
		{"What is the exception hierarchy in Java?", "Java exception hierarchy Throwable"},

		// Collections
		{"When would you use ArrayList vs LinkedList?", "Java ArrayList vs LinkedList"},
		{"What is the difference between HashSet, TreeSet, and LinkedHashSet?", "Java HashSet TreeSet LinkedHashSet"},
		{"Explain the difference between HashMap, TreeMap, and LinkedHashMap", "Java HashMap TreeMap LinkedHashMap"},
		{"What is the difference between Comparable and Comparator?", "Java Comparable vs Comparator"},
		{"Why must you override hashCode when you override equals?", "Java equals hashCode contract"},
		{"What are fail-fast and fail-safe iterators?", "Java fail fast fail safe iterator"},

		// Generics
		{"What are generics and why are they useful?", "Java generics type safety"},
		{"What are bounded type parameters in generics?", "Java bounded type parameters"},
		{"Explain wildcards in generics: extends vs super", "Java generics wildcards PECS"},
		{"What is type erasure in Java generics?", "Java type erasure generics"},

		// Concurrency
		{"What is the difference between extending Thread and implementing Runnable?", "Java Thread vs Runnable"},
		{"What is an ExecutorService and why use thread pools?", "Java ExecutorService thread pool"},
		{"What is the difference between Callable and Runnable?", "Java Callable vs Runnable"},
		{"Explain Future and CompletableFuture", "Java Future CompletableFuture"},
		{"What is the difference between CountDownLatch and CyclicBarrier?", "Java CountDownLatch CyclicBarrier"},

		// Memory & JVM
		{"Explain the difference between stack and heap memory", "Java stack heap memory"},
		{"What is the String pool in Java?", "Java String pool intern"},
		{"How does garbage collection work in Java?", "Java garbage collection"},
		{"What are weak references and when would you use them?", "Java weak reference soft reference"},

		// Streams & Lambdas
		{"How do lambda expressions work in Java?", "Java lambda expressions"},
		{"What are method references and when would you use them?", "Java method reference"},
		{"What is the difference between a Collection and a Stream?", "Java Stream vs Collection"},
		{"What are intermediate vs terminal operations in streams?", "Java stream intermediate terminal operations"},
		{"What is Optional and how does it prevent NullPointerException?", "Java Optional null safety"},

		// Core Concepts
		{"What is the difference between mutable and immutable objects?", "Java mutable immutable String"},
		{"Explain the difference between == and equals()", "Java equals vs =="},
		{"What is autoboxing and unboxing?", "Java autoboxing unboxing"},
		{"Is Java pass by value or pass by reference?", "Java pass by value reference"},
		{"Why doesn't Java support multiple inheritance with classes?", "Java diamond problem multiple inheritance"},
		{"What is serialization in Java?", "Java serialization Serializable"},
		{"What is reflection and when would you use it?", "Java reflection"},
		{"What are annotations in Java?", "Java annotations"},
		{"When would you use StringBuilder vs String?", "Java StringBuilder vs String"},
		{"What is the difference between deep copy and shallow copy?", "Java deep copy shallow copy clone"},
		{"Explain this vs super keywords", "Java this super keyword"},
		{"What is static binding vs dynamic binding?", "Java static dynamic binding"},
		{"What is a marker interface?", "Java marker interface Serializable"},

		// SOLID Principles
		{"Explain the Single Responsibility Principle", "Single Responsibility Principle SOLID"},
		{"Explain the Open/Closed Principle", "Open Closed Principle SOLID"},
	}
}
