package sqlego

import "testing"

func TestSelectStatement(t *testing.T) {
	node := Select("Users", []string{"id", "name", "email"})
	sql := node.Compile()
	if sql != "SELECT id,name,email FROM Users;" {
		t.Fatal(sql)
	}
}
