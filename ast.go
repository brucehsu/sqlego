package sqlego

const (
  STMT_INSERT uint = iota
  STMT_SELECT
  STMT_UPDATE
  STMT_DELETE
  CLAU_WHERE
  CLAU_FROM
  LIST_COLS
  LIST_CONDS
  NODE_PRED
  NODE_AND
  NODE_OR
  PRED_EQ
  PRED_NEQ
  PRED_GT
  PRED_GTE
  PRED_LT
  PRED_LTE
)

type AST struct {
  Type uint
  Values []string
  children []interface{}
}

type Statement struct {
  AST
  table string
  columns []string
}

type Clause struct {
  AST
}

type Predicate struct {
  AST
  operator uint
  first string
  second string
}
