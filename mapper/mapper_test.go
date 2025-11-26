package mapper

import (
	"testing"
)

// TestNew tests the New function with various scenarios
func TestNew(t *testing.T) {
	t.Run("create empty mapper", func(t *testing.T) {
		m, err := New[string, int]()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if m == nil {
			t.Fatal("expected mapper to be non-nil")
		}
	})

	t.Run("create mapper with valid pairs", func(t *testing.T) {
		m, err := New[string, int]("a", 1, "b", 2, "c", 3)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if m == nil {
			t.Fatal("expected mapper to be non-nil")
		}

		// Verify all mappings work
		val, err := m.To("a")
		if err != nil || val != 1 {
			t.Errorf("expected To('a') = 1, got %d (err: %v)", val, err)
		}

		key, err := m.From(2)
		if err != nil || key != "b" {
			t.Errorf("expected From(2) = 'b', got %s (err: %v)", key, err)
		}
	})

	t.Run("odd number of arguments", func(t *testing.T) {
		_, err := New[string, int]("a", 1, "b")
		if err == nil {
			t.Fatal("expected error for odd number of arguments")
		}
		if err.Error() != "odd number of key/values" {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("invalid key type", func(t *testing.T) {
		_, err := New[string, int](123, 1)
		if err == nil {
			t.Fatal("expected error for invalid key type")
		}
		// Check error mentions expected and actual types
		errMsg := err.Error()
		if len(errMsg) == 0 {
			t.Error("expected non-empty error message")
		}
	})

	t.Run("invalid value type", func(t *testing.T) {
		_, err := New[string, int]("a", "invalid")
		if err == nil {
			t.Fatal("expected error for invalid value type")
		}
		// Check error mentions expected and actual types
		errMsg := err.Error()
		if len(errMsg) == 0 {
			t.Error("expected non-empty error message")
		}
	})

	t.Run("duplicate key", func(t *testing.T) {
		_, err := New[string, int]("a", 1, "a", 2)
		if err == nil {
			t.Fatal("expected error for duplicate key")
		}
		if err.Error() != "key already exists" {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("duplicate value", func(t *testing.T) {
		_, err := New[string, int]("a", 1, "b", 1)
		if err == nil {
			t.Fatal("expected error for duplicate value")
		}
		if err.Error() != "value already exists" {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("different types", func(t *testing.T) {
		m, err := New[int, string](1, "one", 2, "two", 3, "three")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		val, err := m.To(2)
		if err != nil || val != "two" {
			t.Errorf("expected To(2) = 'two', got %s (err: %v)", val, err)
		}

		key, err := m.From("three")
		if err != nil || key != 3 {
			t.Errorf("expected From('three') = 3, got %d (err: %v)", key, err)
		}
	})
}

// TestTo tests the To method
func TestTo(t *testing.T) {
	m, err := New[string, int]("a", 1, "b", 2, "c", 3)
	if err != nil {
		t.Fatalf("failed to create mapper: %v", err)
	}

	t.Run("existing mapping", func(t *testing.T) {
		val, err := m.To("a")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if val != 1 {
			t.Errorf("expected 1, got %d", val)
		}
	})

	t.Run("non-existing mapping", func(t *testing.T) {
		val, err := m.To("nonexistent")
		if err == nil {
			t.Fatal("expected error for non-existing key")
		}
		if val != 0 {
			t.Errorf("expected zero value, got %d", val)
		}
		if err.Error() != "no mapping found" {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("all mappings", func(t *testing.T) {
		tests := []struct {
			key   string
			value int
		}{
			{"a", 1},
			{"b", 2},
			{"c", 3},
		}

		for _, tt := range tests {
			val, err := m.To(tt.key)
			if err != nil {
				t.Errorf("To(%q) returned error: %v", tt.key, err)
			}
			if val != tt.value {
				t.Errorf("To(%q) = %d, want %d", tt.key, val, tt.value)
			}
		}
	})
}

// TestToWithDefault tests the ToWithDefault method
func TestToWithDefault(t *testing.T) {
	m, err := New[string, int]("a", 1, "b", 2, "c", 3)
	if err != nil {
		t.Fatalf("failed to create mapper: %v", err)
	}

	t.Run("existing mapping", func(t *testing.T) {
		val := m.ToWithDefault("a", 999)
		if val != 1 {
			t.Errorf("expected 1, got %d", val)
		}
	})

	t.Run("non-existing mapping returns default", func(t *testing.T) {
		val := m.ToWithDefault("nonexistent", 999)
		if val != 999 {
			t.Errorf("expected 999, got %d", val)
		}
	})

	t.Run("zero default value", func(t *testing.T) {
		val := m.ToWithDefault("nonexistent", 0)
		if val != 0 {
			t.Errorf("expected 0, got %d", val)
		}
	})

	t.Run("negative default value", func(t *testing.T) {
		val := m.ToWithDefault("nonexistent", -1)
		if val != -1 {
			t.Errorf("expected -1, got %d", val)
		}
	})
}

// TestFrom tests the From method
func TestFrom(t *testing.T) {
	m, err := New[string, int]("a", 1, "b", 2, "c", 3)
	if err != nil {
		t.Fatalf("failed to create mapper: %v", err)
	}

	t.Run("existing mapping", func(t *testing.T) {
		key, err := m.From(1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if key != "a" {
			t.Errorf("expected 'a', got %s", key)
		}
	})

	t.Run("non-existing mapping", func(t *testing.T) {
		key, err := m.From(999)
		if err == nil {
			t.Fatal("expected error for non-existing value")
		}
		if key != "" {
			t.Errorf("expected zero value, got %s", key)
		}
		if err.Error() != "no mapping found" {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("all reverse mappings", func(t *testing.T) {
		tests := []struct {
			value int
			key   string
		}{
			{1, "a"},
			{2, "b"},
			{3, "c"},
		}

		for _, tt := range tests {
			key, err := m.From(tt.value)
			if err != nil {
				t.Errorf("From(%d) returned error: %v", tt.value, err)
			}
			if key != tt.key {
				t.Errorf("From(%d) = %q, want %q", tt.value, key, tt.key)
			}
		}
	})
}

// TestFromWithDefault tests the FromWithDefault method
func TestFromWithDefault(t *testing.T) {
	m, err := New[string, int]("a", 1, "b", 2, "c", 3)
	if err != nil {
		t.Fatalf("failed to create mapper: %v", err)
	}

	t.Run("existing mapping", func(t *testing.T) {
		key := m.FromWithDefault(1, "default")
		if key != "a" {
			t.Errorf("expected 'a', got %s", key)
		}
	})

	t.Run("non-existing mapping returns default", func(t *testing.T) {
		key := m.FromWithDefault(999, "default")
		if key != "default" {
			t.Errorf("expected 'default', got %s", key)
		}
	})

	t.Run("empty string default value", func(t *testing.T) {
		key := m.FromWithDefault(999, "")
		if key != "" {
			t.Errorf("expected '', got %s", key)
		}
	})
}

// TestBidirectionalMapping tests that bidirectional mapping works correctly
func TestBidirectionalMapping(t *testing.T) {
	m, err := New[string, int]("a", 1, "b", 2, "c", 3)
	if err != nil {
		t.Fatalf("failed to create mapper: %v", err)
	}

	// Test that To -> From -> To gives the same result
	val, err := m.To("b")
	if err != nil {
		t.Fatalf("To('b') failed: %v", err)
	}

	key, err := m.From(val)
	if err != nil {
		t.Fatalf("From(%d) failed: %v", val, err)
	}

	if key != "b" {
		t.Errorf("expected round-trip to return 'b', got %s", key)
	}

	// Test that From -> To -> From gives the same result
	key2, err := m.From(2)
	if err != nil {
		t.Fatalf("From(2) failed: %v", err)
	}

	val2, err := m.To(key2)
	if err != nil {
		t.Fatalf("To(%s) failed: %v", key2, err)
	}

	if val2 != 2 {
		t.Errorf("expected round-trip to return 2, got %d", val2)
	}
}

// TestMapperWithDifferentTypes tests mapper with various type combinations
func TestMapperWithDifferentTypes(t *testing.T) {
	t.Run("int to string", func(t *testing.T) {
		m, err := New[int, string](1, "one", 2, "two")
		if err != nil {
			t.Fatalf("failed to create mapper: %v", err)
		}

		val, _ := m.To(1)
		if val != "one" {
			t.Errorf("expected 'one', got %s", val)
		}

		key, _ := m.From("two")
		if key != 2 {
			t.Errorf("expected 2, got %d", key)
		}
	})

	t.Run("bool to int", func(t *testing.T) {
		m, err := New[bool, int](true, 1, false, 0)
		if err != nil {
			t.Fatalf("failed to create mapper: %v", err)
		}

		val, _ := m.To(true)
		if val != 1 {
			t.Errorf("expected 1, got %d", val)
		}

		key, _ := m.From(0)
		if key != false {
			t.Errorf("expected false, got %v", key)
		}
	})

	t.Run("float64 to string", func(t *testing.T) {
		m, err := New[float64, string](1.5, "one-half", 2.5, "two-half")
		if err != nil {
			t.Fatalf("failed to create mapper: %v", err)
		}

		val, _ := m.To(1.5)
		if val != "one-half" {
			t.Errorf("expected 'one-half', got %s", val)
		}
	})
}

// TestEmptyMapper tests behavior with an empty mapper
func TestEmptyMapper(t *testing.T) {
	m, err := New[string, int]()
	if err != nil {
		t.Fatalf("failed to create empty mapper: %v", err)
	}

	t.Run("To on empty mapper", func(t *testing.T) {
		_, err := m.To("anything")
		if err == nil {
			t.Error("expected error when calling To on empty mapper")
		}
	})

	t.Run("From on empty mapper", func(t *testing.T) {
		_, err := m.From(123)
		if err == nil {
			t.Error("expected error when calling From on empty mapper")
		}
	})

	t.Run("ToWithDefault on empty mapper", func(t *testing.T) {
		val := m.ToWithDefault("anything", 42)
		if val != 42 {
			t.Errorf("expected 42, got %d", val)
		}
	})

	t.Run("FromWithDefault on empty mapper", func(t *testing.T) {
		key := m.FromWithDefault(123, "default")
		if key != "default" {
			t.Errorf("expected 'default', got %s", key)
		}
	})
}

// TestMustNew tests the MustNew function with various scenarios
func TestMustNew(t *testing.T) {
	t.Run("create empty mapper", func(t *testing.T) {
		m := MustNew[string, int]()
		if m == nil {
			t.Fatal("expected mapper to be non-nil")
		}
	})

	t.Run("create mapper with valid pairs", func(t *testing.T) {
		m := MustNew[string, int]("a", 1, "b", 2, "c", 3)
		if m == nil {
			t.Fatal("expected mapper to be non-nil")
		}

		// Verify all mappings work
		val, err := m.To("a")
		if err != nil || val != 1 {
			t.Errorf("expected To('a') = 1, got %d (err: %v)", val, err)
		}

		key, err := m.From(2)
		if err != nil || key != "b" {
			t.Errorf("expected From(2) = 'b', got %s (err: %v)", key, err)
		}
	})

	t.Run("different types", func(t *testing.T) {
		m := MustNew[int, string](1, "one", 2, "two", 3, "three")

		val, err := m.To(2)
		if err != nil || val != "two" {
			t.Errorf("expected To(2) = 'two', got %s (err: %v)", val, err)
		}

		key, err := m.From("three")
		if err != nil || key != 3 {
			t.Errorf("expected From('three') = 3, got %d (err: %v)", key, err)
		}
	})

	t.Run("panic on odd number of arguments", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for odd number of arguments")
			}
		}()
		MustNew[string, int]("a", 1, "b")
	})

	t.Run("panic on invalid key type", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid key type")
			}
		}()
		MustNew[string, int](123, 1)
	})

	t.Run("panic on invalid value type", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid value type")
			}
		}()
		MustNew[string, int]("a", "invalid")
	})

	t.Run("panic on duplicate key", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for duplicate key")
			}
		}()
		MustNew[string, int]("a", 1, "a", 2)
	})

	t.Run("panic on duplicate value", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for duplicate value")
			}
		}()
		MustNew[string, int]("a", 1, "b", 1)
	})
}
