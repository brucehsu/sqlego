package sqlego

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

func Select(table string, columns []string) *Statement {
	node := &Statement{table: table, columns: columns}
	node.Type = STMT_SELECT
	return node
}

func Insert(table string, values map[string]string) *Statement {
	node := &Statement{table: table}
	for k, v := range values {
		node.columns = append(node.columns, k)
		node.Values = append(node.Values, v)
	}
	node.Type = STMT_INSERT
	return node
}

func Update(table string, values map[string]string) *Statement {
	node := &Statement{table: table}
	for k, v := range values {
		node.columns = append(node.columns, k)
		node.Values = append(node.Values, v)
	}
	node.Type = STMT_UPDATE
	return node
}

func Delete(table string) *Statement {
	node := &Statement{table: table}
	node.Type = STMT_DELETE
	return node
}

func newTernaryPredicate(operand string, operand_second string, operator uint) *Predicate {
	pred := &Predicate{first: operand, second: operand_second, operator: operator}
	pred.Type = NODE_PRED
	return pred
}

func ExplicitPredicates(preds ...*Predicate) *Predicate {
	dummy_pred := &Predicate{}
	dummy_pred.Type = NODE_PRED
	exp_pred := &Predicate{}
	exp_pred.Type = NODE_EXPPRED
	for _, pred := range preds {
		exp_pred.children = append(exp_pred.children, pred)
	}
	dummy_pred.children = append(dummy_pred.children, exp_pred)
	return dummy_pred
}

func Eq(operand string, operand_second string) *Predicate {
	return newTernaryPredicate(operand, operand_second, PRED_EQ)
}

func Neq(operand string, operand_second string) *Predicate {
	return newTernaryPredicate(operand, operand_second, PRED_NEQ)
}

func Gt(operand string, operand_second string) *Predicate {
	return newTernaryPredicate(operand, operand_second, PRED_GT)
}

func Gte(operand string, operand_second string) *Predicate {
	return newTernaryPredicate(operand, operand_second, PRED_GTE)
}

func Lt(operand string, operand_second string) *Predicate {
	return newTernaryPredicate(operand, operand_second, PRED_LT)
}

func Lte(operand string, operand_second string) *Predicate {
	return newTernaryPredicate(operand, operand_second, PRED_LTE)
}

func (pred *Predicate) And(preds ...*Predicate) *Predicate {
	and_pred := &Predicate{}
	and_pred.Type = NODE_AND
	for _, child := range preds {
		and_pred.children = append(and_pred.children, child)
	}
	pred.children = append(pred.children, and_pred)
	return pred
}

func (pred *Predicate) Or(preds ...*Predicate) *Predicate {
	or_pred := &Predicate{}
	or_pred.Type = NODE_OR
	for _, child := range preds {
		or_pred.children = append(or_pred.children, child)
	}
	pred.children = append(pred.children, or_pred)
	return pred
}

func (node *Statement) Where(preds ...*Predicate) *Statement {
	var where_clause *Clause
	for _, child := range node.children {
		if reflect.TypeOf(child) == reflect.TypeOf(&Clause{}) {
			if child.(*Clause).Type == CLAU_WHERE {
				where_clause = child.(*Clause)
				break
			}
		}
	}
	if where_clause == nil {
		where_clause = &Clause{}
		where_clause.Type = CLAU_WHERE
		node.children = append(node.children, where_clause)
	}
	for _, pred := range preds {
		where_clause.children = append(where_clause.children, pred)
	}
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
		buffer.WriteString(strings.Join(node.columns, ","))
		buffer.WriteString(") VALUES (")
		buffer.WriteString(strings.Join(node.Values, ","))
		buffer.WriteString(")")
	case STMT_UPDATE:
		buffer.WriteString("UPDATE ")
		buffer.WriteString(node.table)
		buffer.WriteString(" SET ")
		predicates := []string{}
		for idx, col := range node.columns {
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

func compileWhere(children []interface{}) string {
	var buffer bytes.Buffer
	for _, obj := range children {
		if reflect.TypeOf(obj) == reflect.TypeOf(&Clause{}) {
			clause := obj.(*Clause)
			if clause.Type == CLAU_WHERE {
				buffer.WriteString(" WHERE ")
				buffer.WriteString(compileImplicitPredicates(clause.children, NODE_AND))
			}
		}
	}
	return buffer.String()
}

func compileImplicitPredicates(children []interface{}, concat_mode uint) string {
	var buffer bytes.Buffer
	for idx, node := range children {
		pred := node.(*Predicate)

		if idx > 0 && pred.Type <= NODE_EXPPRED {
			switch concat_mode {
			case NODE_AND:
				buffer.WriteString(" AND ")
			case NODE_OR:
				buffer.WriteString(" OR ")
			}
		}

		switch pred.Type {
		case NODE_EXPPRED:
			buffer.WriteString(" ( ")
			buffer.WriteString(compileImplicitPredicates(pred.children, NODE_AND))
			buffer.WriteString(" ) ")
		case NODE_PRED:
			buffer.WriteString(compilePredicate(pred))
			if len(pred.children) > 0 {
				buffer.WriteString(compileImplicitPredicates(pred.children, NODE_AND))
			}
		case NODE_AND:
			buffer.WriteString(" AND ")
			buffer.WriteString(compileImplicitPredicates(pred.children, NODE_AND))
		case NODE_OR:
			buffer.WriteString(" OR ")
			buffer.WriteString(compileImplicitPredicates(pred.children, NODE_OR))
		}
	}
	return buffer.String()
}

func compilePredicate(pred *Predicate) string {
	var buffer bytes.Buffer
	if pred.first == "" {
		return ""
	}
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
	return buffer.String()
}
