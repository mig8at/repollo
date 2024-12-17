# Repollo: Generic Data Collection Library for Go

Repollo is a simple and flexible Go library for managing collections of data. It provides methods to create, query, and manipulate collections using powerful filtering and transformation techniques. Designed for modularity and ease of use, it is ideal for projects requiring efficient and reusable data management.

---

## Features

- **Generic Collections**: Manage collections of any type using Go generics.
- **CRUD Operations**: Create, Read, Update, and Delete elements in collections.
- **Query Chaining**: Filter, sort, and limit data using a functional style.
- **Thread-Safe**: Built-in concurrency safety using mutexes.
- **Event Streams**: Subscribe to collection events using channels.

---

## Installation

To install the library, run:

```bash
go get github.com/mig8at/repollo
```

Import the library in your project:

```go
import "github.com/mig8at/repollo"
```

---

## Usage

### 1. Creating a Collection

Create a new collection to manage your data:

```go
type User struct {
	Name string
	Age  int
}

users := mycollection.NewCollection[User]()
users.Create("1", User{"Alice Smith", 25})
users.Create("2", User{"Bob Johnson", 30})
```

### 2. Querying Data

Use the `Where` method to filter data:

```go
filtered := users.Where(func(u User) bool {
    return u.Name == "Alice Smith"
}).Results()
fmt.Println(filtered) // User
```

### 3. Advanced Query Chaining

Perform more complex operations like limiting, offsetting, and sorting:

```go
sorted := users.Where(func(u User) bool {
    return u.Age > 20
}).Sort(func(a, b string) bool {
    return a < b
}).Limit(1).Results()
fmt.Println(sorted)
```

### 4. Event Subscription

Listen to collection events using the `Events` method:

```go
eventChannel := users.Events()
go func() {
    for event := range eventChannel {
        fmt.Printf("Event: %s, Key: %s, Value: %+v\n", event.Type, event.Key, event.Value)
    }
}()

users.Create("2", User{"Bob Johnson", 30})
```

### 5. Thread Safety

MyCollection ensures that operations are thread-safe, so you can use it in concurrent environments without additional locking mechanisms.

---

## API Reference

### Collection Methods

- `NewCollection[T any]() *Collection[T]`: Creates a new collection.
- `Create(key string, value T)`: Adds a new item to the collection.
- `Get(key string) (T, error)`: Retrieves an item by key.
- `Update(key string, value T) error`: Updates an existing item.
- `Delete(key string) error`: Deletes an item by key.
- `Where(predicate func(T) bool) *QueryResult[T]`: Filters the collection based on a predicate.
- `Events() <-chan Event[T]`: Returns a channel that emits collection events.

### QueryResult Methods

- `Results() []T`: Returns the final list of results.
- `Limit(n int) *QueryResult[T]`: Limits the number of results.
- `Offset(n int) *QueryResult[T]`: Skips the first `n` results.
- `Sort(less func(a, b T) bool) *QueryResult[T]`: Sorts results based on a comparison function.

---

## Running Tests

To test the library, use:

```bash
go test ./...
```

Example test case:

```go
func TestCollection(t *testing.T) {
	users := NewCollection[User]()
	users.Create("1", User{"Alice Smith", 25})
	users.Create("2", User{"Bob Johnson", 30})

	if len(users.data) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users.data))
	}

	filtered := users.Where(func(u User) bool {
		return u.Name == "Alice Smith"
	}).Results()

	if len(filtered) != 1 || filtered[0].Name != "Alice Smith" {
		t.Errorf("Expected 'Alice', got %+v", filtered)
	}
}

```

---

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests to improve the library.

### Steps to Contribute

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature-branch`).
3. Commit your changes (`git commit -m "Add feature"`).
4. Push to the branch (`git push origin feature-branch`).
5. Open a pull request.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

