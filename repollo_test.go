package repollo

import (
	"testing"
)

type User struct {
	Name string
	Age  int
}

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
