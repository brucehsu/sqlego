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
  NODE_COND
  NODE_AND
  NODE_OR
)

type AST struct {
  Type uint
  Values []string
  children []*AST
}

type Statement struct {
  AST
  table string
  columns []string
}

type Clause struct {
  AST
}
