package knowledge

func getJavaBits() []KnowledgeBit {
	return []KnowledgeBit{
		// OOP Pillars
		{"What is encapsulation and why is it important in OOP?", "Java encapsulation tutorial example"},
		{"Explain inheritance in Java with an example", "Java inheritance tutorial example"},
		{"What is polymorphism? Explain with a real-world example", "Java polymorphism tutorial example"},
		{"What is abstraction and how does it differ from encapsulation?", "Java abstraction tutorial example"},

		// Polymorphism
		{"What is the difference between method overloading and overriding?", "Java overloading vs overriding example"},
		{"What is compile-time polymorphism vs runtime polymorphism?", "Java compile time vs runtime polymorphism example"},
		{"What are covariant return types in Java?", "Java covariant return types example"},
		{"What are the rules for method overriding in Java?", "Java method overriding rules example"},

		// Abstract Class vs Interface
		{"When would you use an abstract class vs an interface?", "Java abstract class vs interface tutorial"},
		{"What are default methods in interfaces? Why were they added?", "Java interface default methods tutorial"},
		{"What is a functional interface and how is it used with lambdas?", "Java functional interface lambda tutorial"},

		// Design Patterns
		{"Explain composition over inheritance with an example", "Java composition over inheritance example"},
		{"What is the Liskov Substitution Principle?", "Liskov Substitution Principle Java example"},
		{"Explain dependency injection and its benefits", "dependency injection Java tutorial"},
		{"How would you implement the Singleton pattern in Java?", "Java singleton pattern tutorial"},
		{"When would you use the Factory pattern?", "Java factory pattern tutorial"},
		{"Explain the Builder pattern with an example", "Java builder pattern tutorial"},

		// Modifiers
		{"What is the difference between static and instance variables?", "Java static vs instance variable example"},
		{"What does the final keyword mean for variables, methods, and classes?", "Java final keyword tutorial"},
		{"Explain Java access modifiers: public, private, protected, default", "Java access modifiers tutorial"},
		{"What is the transient keyword used for?", "Java transient keyword example"},
		{"When would you use the volatile keyword?", "Java volatile keyword tutorial"},
		{"What does synchronized do and when would you use it?", "Java synchronized keyword tutorial"},

		// Exceptions
		{"What is the difference between checked and unchecked exceptions?", "Java checked vs unchecked exceptions tutorial"},
		{"Explain try-with-resources and its benefits", "Java try with resources tutorial"},
		{"What is the exception hierarchy in Java?", "Java exception hierarchy explained"},

		// Collections
		{"When would you use ArrayList vs LinkedList?", "Java ArrayList vs LinkedList tutorial"},
		{"What is the difference between HashSet, TreeSet, and LinkedHashSet?", "Java Set comparison tutorial"},
		{"Explain the difference between HashMap, TreeMap, and LinkedHashMap", "Java Map comparison tutorial"},
		{"What is the difference between Comparable and Comparator?", "Java Comparable vs Comparator tutorial"},
		{"Why must you override hashCode when you override equals?", "Java equals hashCode contract explained"},
		{"What are fail-fast and fail-safe iterators?", "Java fail fast iterator tutorial"},

		// Generics
		{"What are generics and why are they useful?", "Java generics tutorial"},
		{"What are bounded type parameters in generics?", "Java bounded type parameters example"},
		{"Explain wildcards in generics: extends vs super", "Java generics wildcards PECS tutorial"},
		{"What is type erasure in Java generics?", "Java type erasure explained"},

		// Concurrency
		{"What is the difference between extending Thread and implementing Runnable?", "Java Thread vs Runnable tutorial"},
		{"What is an ExecutorService and why use thread pools?", "Java ExecutorService tutorial"},
		{"What is the difference between Callable and Runnable?", "Java Callable vs Runnable example"},
		{"Explain Future and CompletableFuture", "Java CompletableFuture tutorial"},
		{"What is the difference between CountDownLatch and CyclicBarrier?", "Java CountDownLatch CyclicBarrier tutorial"},

		// Memory & JVM
		{"Explain the difference between stack and heap memory", "Java stack vs heap memory explained"},
		{"What is the String pool in Java?", "Java String pool tutorial"},
		{"How does garbage collection work in Java?", "Java garbage collection tutorial"},
		{"What are weak references and when would you use them?", "Java weak reference tutorial"},

		// Streams & Lambdas
		{"How do lambda expressions work in Java?", "Java lambda expressions tutorial"},
		{"What are method references and when would you use them?", "Java method reference tutorial"},
		{"What is the difference between a Collection and a Stream?", "Java Stream vs Collection tutorial"},
		{"What are intermediate vs terminal operations in streams?", "Java stream operations tutorial"},
		{"What is Optional and how does it prevent NullPointerException?", "Java Optional tutorial"},

		// Core Concepts
		{"What is the difference between mutable and immutable objects?", "Java immutable objects tutorial"},
		{"Explain the difference between == and equals()", "Java equals vs == explained"},
		{"What is autoboxing and unboxing?", "Java autoboxing tutorial"},
		{"Is Java pass by value or pass by reference?", "Java pass by value explained"},
		{"Why doesn't Java support multiple inheritance with classes?", "Java diamond problem explained"},
		{"What is serialization in Java?", "Java serialization tutorial"},
		{"What is reflection and when would you use it?", "Java reflection tutorial"},
		{"What are annotations in Java?", "Java annotations tutorial"},
		{"When would you use StringBuilder vs String?", "Java StringBuilder vs String tutorial"},
		{"What is the difference between deep copy and shallow copy?", "Java deep copy vs shallow copy tutorial"},
		{"Explain this vs super keywords", "Java this vs super keyword tutorial"},
		{"What is static binding vs dynamic binding?", "Java static vs dynamic binding example"},
		{"What is a marker interface?", "Java marker interface explained"},

		// SOLID Principles
		{"Explain the Single Responsibility Principle", "Single Responsibility Principle tutorial"},
		{"Explain the Open/Closed Principle", "Open Closed Principle tutorial"},
	}
}
