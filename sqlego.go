package sqlego

import (
  "fmt"
  "bytes"
  "strings"
)

func Select(table string, columns []string) (*Statement) {
  node := &Statement{table: table, columns: columns}
  node.Type = STMT_SELECT
  return node
}

func Insert(table string, columns []string, values []string) (*Statement) {
  if len(columns) != len(values) {
    panic("Values cannot be mapped to columns")
  }
  node := &Statement{table: table, columns: columns}
  node.Type = STMT_INSERT
  node.Values = values
  return node
}

func Update(table string, columns []string, values []string) (*Statement) {
  if len(columns) != len(values) {
    panic("Values cannot be mapped to columns")
  }
  node := &Statement{table: table, columns: columns}
  node.Type = STMT_UPDATE
  node.Values = values
  return node
}

func Delete(table string) (*Statement) {
  node := &Statement{table: table}
  node.Type = STMT_DELETE
  return node
}

func (node *Statement) Compile() string {
  var buffer bytes.Buffer
  switch node.Type {
  case STMT_SELECT:
    buffer.WriteString("SELECT ")
    buffer.WriteString(strings.Join(node.columns, ","))
    buffer.WriteString(" FROM ")
    buffer.WriteString(node.table)
    buffer.WriteString(";")
  case STMT_INSERT:
    buffer.WriteString("INSERT INTO ")
    buffer.WriteString(node.table)
    buffer.WriteString(" (")
    buffer.WriteString(strings.Join(node.columns,","))
    buffer.WriteString(") VALUES (")
    buffer.WriteString(strings.Join(node.Values, ","))
    buffer.WriteString(");")
  case STMT_UPDATE:
    buffer.WriteString("UPDATE ")
    buffer.WriteString(node.table)
    buffer.WriteString(" SET ")
    predicates := []string{}
    for idx, col := range(node.columns) {
      predicates = append(predicates, fmt.Sprintf("%s=%s", col, node.Values[idx]))
    }
    buffer.WriteString(strings.Join(predicates, ","))
    buffer.WriteString(";")
  case STMT_DELETE:
    buffer.WriteString("DELETE FROM ")
    buffer.WriteString(node.table)
    buffer.WriteString(";")
  }

  return buffer.String()
}
