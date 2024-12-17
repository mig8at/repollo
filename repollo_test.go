package repollo

import (
	"testing"
)

func TestCollection(t *testing.T) {
	users := NewCollection[string]()
	users.Create("1", "Alice")
	users.Create("2", "Bob")

	if len(users.data) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users.data))
	}

	filtered := users.Where(func(u string) bool {
		return u == "Alice"
	}).Results()

	if len(filtered) != 1 || filtered[0] != "Alice" {
		t.Errorf("Expected 'Alice', got %+v", filtered)
	}
}
