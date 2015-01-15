package sqlego

import (
	s "strings"
	"testing"
)

func TestSelectStatement(t *testing.T) {
	node := Select("Users", []string{"id", "name", "email"})
	sql := node.Compile()
	if sql != "SELECT id,name,email FROM Users;" {
		t.Fatal(sql)
	}
}

func TestSelectStatementWithWhere(t *testing.T) {
	node := Select("Users", []string{"id", "name", "email"})
	node.Where(Gte("id", "10"))
	sql := node.Compile()
	if sql != "SELECT id,name,email FROM Users WHERE id>=10;" {
		t.Fatalf("Single condition:\n%s", sql)
	}

	// Test implicit concatenation
	node = Select("Users", []string{"id", "name", "email"})
	node.Where(Gte("id", "10"), Lt("id", "20"))
	sql = node.Compile()
	if sql != "SELECT id,name,email FROM Users WHERE id>=10 AND id<20;" {
		t.Fatalf("Implicit concatenation:\n%s", sql)
	}

	// Test explicit predicates
	node = Select("Users", []string{"id", "name", "email"})
	node.Where(ExplicitPredicates(Gte("id", "10"), Lt("id", "20")))
	sql = node.Compile()
	if sql != "SELECT id,name,email FROM Users WHERE  ( id>=10 AND id<20 ) ;" {
		t.Fatalf("Explicit predicates:\n%s", sql)
	}

	// Test explicit AND concatenation
	node = Select("Users", []string{"id", "name", "email"})
	node.Where(Gte("id", "10").And(Lt("id", "20")))
	sql = node.Compile()
	if sql != "SELECT id,name,email FROM Users WHERE id>=10 AND id<20;" {
		t.Fatalf("Explicit AND concatenation:\n%s", sql)
	}
}

func TestInsertStatement(t *testing.T) {
	node := Insert("Users", map[string]string{"name": "Bruce"})
	sql := node.Compile()
	if sql != "INSERT INTO Users (name) VALUES (Bruce);" {
		t.Fatal(sql)
	}
}

func TestUpdateStatement(t *testing.T) {
	node := Update("Users", map[string]string{"id": "2", "name": "Bruce", "email": "bruce@example.com"})
	sql := node.Compile()
	if !s.Contains(sql, "UPDATE Users SET") && !s.Contains(sql, "id=2") && !s.Contains(sql, "name=Bruce") && !s.Contains(sql, "email=bruce@example.com") {
		t.Fatal(sql)
	}
}

func TestDeleteStatement(t *testing.T) {
	node := Delete("Users")
	sql := node.Compile()
	if sql != "DELETE FROM Users;" {
		t.Fatal(sql)
	}
}
