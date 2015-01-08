sqlego
======

SQL query building interface for Go.

# Usage
## Basic CRUD operations

```go
func Select(table string, columns []string) *Statement
```

```go
func Insert(table string, values map[string]string) *Statement
```

```go
func Update(table string, values map[string]string) *Statement
```

```go
func Delete(table string) *Statement
```

### Examples
- Select ``id, name, email`` from a table ``Users``
```go
sqlego.Select("Users", []string{"id", "name", "email"})
```
- Insert/Update a new record to table ``Users``
```go
sqlego.Insert("Users", map[string]string{"id": "2", "name": "brucehsu"}
sqlego.Update("Users", map[string]string{"id": "2", "name": "brucehsu"}
```
- Delete all records in table ``Users``
```go
sqlego.Delete("Users")
```

## Conditions
WHERE clause
```go
func (node *Statement) Where(preds ...*Predicate) *Statement
```

Adding condition predicates
```go
func Eq(operand string, operand_second string) *Predicate // =
```
```go
func Neq(operand string, operand_second string) *Predicate // <>
```
```go
func Gt(operand string, operand_second string) *Predicate // >
```
```go
func Gte(operand string, operand_second string) *Predicate // >=
```
```go
func Lt(operand string, operand_second string) *Predicate // <
```
```go
func Lte(operand string, operand_second string) *Predicate // <=
```

Concatenating conditions
```go
func (pred *Predicate) And(preds ...*Predicate) *Predicate
```
```go
func (pred *Predicate) Or(preds ...*Predicate) *Predicate
```

Grouping predicates explicitly
```go
func ExplicitPredicates(preds ...*Predicate) *Predicate
```

### Note

In functions that accept variadic arguments, given predicates would be concatenated with ``AND`` by default except for ``Or()``. 

### Examples
- Select the user record of ``Bruce`` with following columns ``id, name, email`` from a table ``Users``
```go
node := sqlego.Select("Users", []string{"id", "name", "email"})
node.Where(sqlego.Eq("name", "Bruce"))
node.Compile() // => SELECT id,name,email FROM Users WHERE name=Bruce;
```
- Select the user records whose ``id`` equals to ``1``  or range from ``5566`` to ``7788`` from a table ``Users``
```go
node := sqlego.Select("Users", []string{"id", "name", "email"})
node.Where(sqlego.Eq("id", "1").Or(sqlego.ExplicitPredicates(sqlego.Gte("id", "5566"), sqlego.Lte("id", "7788"))))
node.Compile() // => SELECT id,name,email FROM Users WHERE id=1 OR  ( id>=5566 AND id<=7788 ) ;
```
