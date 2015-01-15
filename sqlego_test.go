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
