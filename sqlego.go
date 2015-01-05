package sqlego

import (
  "fmt"
  "bytes"
  "strings"
  "reflect"
)

func Select(table string, columns []string) (*Statement) {
  node := &Statement{table: table, columns: columns}
  node.Type = STMT_SELECT
  return node
}

func Insert(table string, values map[string]string) (*Statement) {
  node := &Statement{table: table}
  for k, v := range values {
    node.columns = append(node.columns, k)
    node.Values = append(node.Values, v)
  }
  node.Type = STMT_INSERT
  return node
}

func Update(table string, values map[string]string) (*Statement) {
  node := &Statement{table: table}
  for k, v := range values {
    node.columns = append(node.columns, k)
    node.Values = append(node.Values, v)
  }
  node.Type = STMT_UPDATE
  return node
}

func Delete(table string) (*Statement) {
  node := &Statement{table: table}
  node.Type = STMT_DELETE
  return node
}

func newTernaryPredicate(operand string, operand_second string, operator uint) (*Predicate) {
  pred := &Predicate{first: operand, second: operand_second, operator: operator}
  pred.Type = NODE_PRED
  return pred
}

func Eq(operand string, operand_second string) (*Predicate) {
  return newTernaryPredicate(operand, operand_second, PRED_EQ)
}

func Neq(operand string, operand_second string) (*Predicate) {
  return newTernaryPredicate(operand, operand_second, PRED_NEQ)
}

func Gt(operand string, operand_second string) (*Predicate) {
  return newTernaryPredicate(operand, operand_second, PRED_GT)
}

func Gte(operand string, operand_second string) (*Predicate) {
  return newTernaryPredicate(operand, operand_second, PRED_GTE)
}

func Lt(operand string, operand_second string) (*Predicate) {
  return newTernaryPredicate(operand, operand_second, PRED_LT)
}

func Lte(operand string, operand_second string) (*Predicate) {
  return newTernaryPredicate(operand, operand_second, PRED_LTE)
}

func (node *Statement) Where(pred *Predicate) (*Statement) {
  where_clause := &Clause{}
  where_clause.Type = CLAU_WHERE
  node.children = append(node.children, where_clause)
  where_clause.children = append(where_clause.children, pred)
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
  case STMT_INSERT:
    buffer.WriteString("INSERT INTO ")
    buffer.WriteString(node.table)
    buffer.WriteString(" (")
    buffer.WriteString(strings.Join(node.columns,","))
    buffer.WriteString(") VALUES (")
    buffer.WriteString(strings.Join(node.Values, ","))
    buffer.WriteString(")")
  case STMT_UPDATE:
    buffer.WriteString("UPDATE ")
    buffer.WriteString(node.table)
    buffer.WriteString(" SET ")
    predicates := []string{}
    for idx, col := range(node.columns) {
      predicates = append(predicates, fmt.Sprintf("%s=%s", col, node.Values[idx]))
    }
    buffer.WriteString(strings.Join(predicates, ","))
  case STMT_DELETE:
    buffer.WriteString("DELETE FROM ")
    buffer.WriteString(node.table)
  }

  // Generate WHERE clause if exists
  buffer.WriteString(compileWhere(node.children))

  buffer.WriteString(";")
  return buffer.String()
}

func compileWhere(children []interface{}) (string) {
  var buffer bytes.Buffer
  for _, obj := range children {
    if reflect.TypeOf(obj) == reflect.TypeOf(&Clause{}) {
      clause := obj.(*Clause)
      if clause.Type == CLAU_WHERE {
        buffer.WriteString(" WHERE ")
        // We only have 1 predicate for now
        buffer.WriteString(compilePredicates(clause.children))
      }
    }
  }
  return buffer.String()
}

func compilePredicates(children []interface{}) (string) {
  // We assume given argument is a list of predicate nodes
  var buffer bytes.Buffer
  for _, node := range children {
    pred := node.(*Predicate)
    buffer.WriteString(pred.first)
    switch pred.operator {
    case PRED_EQ:
      buffer.WriteString("=")
    case PRED_NEQ:
      buffer.WriteString("<>")
    case PRED_GT:
      buffer.WriteString(">")
    case PRED_GTE:
      buffer.WriteString(">=")
    case PRED_LT:
      buffer.WriteString("<")
    case PRED_LTE:
      buffer.WriteString("<=")
    }
    buffer.WriteString(pred.second)
  }
  return buffer.String()
}
