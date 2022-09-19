# Expression

This library provides an expression parser and evaluator. It has the following features:

* Literals
  * Numbers, i.e. `2`, `4.5`
  * Strings, i.e. `'hello'`
    * `\` and `'` must be escaped, i.e. `'a man called \'Dan\''`
    * Can otherwise contain any character
   * Booleans, i.e. `true`, `false`
* Indexing into an accompanying JSTN typed JSON object
  * Starts with `$`
  * Can use object key notation, i.e. `$.identifier`
    * Key must start with `[a-zA-Z_]` and otherwise only contain `[a-zA-Z0-9_]`
  * Can use index notation for object keys, i.e. `$['identifier']`
    * `\` and `'` must be escaped, i.e. `$['a man called \'Dan\'']`
    * Can otherwise contain any character
  * Can use index notation for array indeces, i.e. `$[2]`
  * Can be chained together, i.e. `$.indentifier[2]['id2']`
  * Can include a "value not found" operator `?`, i.e. `$.myBool ? true`
    * If this is not used and the value indicated by the path is not found in the object, the default value for the type will be used (`0` for number, `''` for string, `false` for bool, `null` for array)
* Unary operators
  * Logical not `!`
  * Number inverter `-`, i.e. `-42`
* Binary operators
  * Number addition `+`, i.e. `3 + 4`
  * Number subraction `-`, i.e. `5 - 2.2`
  * Number multiplication `*`, i.e. `4 * 4.2`
  * Number division `/`, i.e. `5 / 2`
  * Number comparisons `<`, `<=`, `>`, `>=`
  * Number/string/boolean equality `==`
  * Logical and `&&`, i.e. `(5 < 10) && (2 < 1)`
  * Logical or `||`
* Functions
  * Variadic argument number sum `sum`, i.e. `sum(1, 3, 5)`
  * Variadic argument number product `product`, i.e. `product(2, 3, 4)`
  * Variadic argument logical and `and`, i.e. `and($.thing1, $.thing2, $.thing3)`
  * Variadic argument logic or `or`
  * Array length `length`, i.e. `length($.myArray)`

## Language specification

```
expr := numExpr | boolExpr | strExpr | arrExpr

numExpr := NUMBER | PATH | numPathExpr | invExpr | addExpr | addFnExpr | subExpr | mulExpr | mulFnExpr | divExpr | lenExpr

boolExpr := BOOL | PATH | boolPathExpr | notExpr | cmpExpr | eqlExpr | andExpr | andFnExpr | orExpr | orFnExpr

strExpr := STRING | PATH | strPathExpr

arrExpr := PATH

numPathExpr := PATH IF_NOT_FOUND numParenExpr

invExpr := MINUS numParenExpr

addExpr := numParenExpr PLUS numParenExpr

addFnExpr := SUM LEFT_PAREN numExpr numExprList RIGHT_PAREN

subExpr := numParenExpr MINUS numParenExpr

mulExpr := numParenExpr TIMES numParenExpr

mulFnExpr := PRODUCT LEFT_PAREN numExpr numExprList RIGHT_PAREN

divExpr := numParenExpr DIVIDE numParenExpr

lenExpr := LENGTH LEFT_PAREN arrExpr RIGHT_PAREN

boolPathExpr := PATH IF_NOT_FOUND boolParenExpr

notExpr := NOT boolParenExpr

cmpExpr := numParenExpr LESS numParenExpr | numParenExpr LESS_EQ numParenExpr | numParenExpr MORE numParenExpr | numParenExpr MORE_EQ numParenExpr

eqlExpr := numParenExpr EQ numParenExpr | boolParenExpr EQ boolParenExpr | strParenExpr EQ strParenExpr

andExpr := boolParenExpr AND_OP boolParenExpr

andFnExpr := AND LEFT_PAREN boolExpr boolExprList RIGHT_PAREN

orExpr := boolParenExpr OR_OP boolParenExpr

orFnExpr := OR LEFT_PAREN boolExpr boolExprList RIGHT_PAREN

strPathExpr := PATH IF_NOT_FOUND strParenExpr

numParenExpr := NUMBER | PATH | LEFT_PAREN numPathExpr RIGHT_PAREN | LEFT_PAREN invExpr RIGHT_PAREN | LEFT_PAREN addExpr RIGHT_PAREN | addFnExpor | LEFT_PAREN subExpr RIGHT_PAREN | LEFT_PAREN mulExpr RIGHT_PAREN | mulFnExpr | LEFT_PAREN divExpr RIGHT_PAREN | lenExpr

numExprList := _ | COMMA numExpr numExprList

boolParenExpr := BOOL | PATH | LEFT_PAREN boolPathExpr RIGHT_PAREN | LEFT_PAREN notExpr RIGHT_PAREN | LEFT_PAREN cmpExpr RIGHT_PAREN | LEFT_PAREN eqlExpr RIGHT_PAREN | LEFT_PAREN andExpr RIGHT_PAREN | andFnExpr | LEFT_PAREN orExpr RIGHT_PAREN | orFnExpr

boolExprList := _ | COMMA boolExpr boolExprList

strParenExpr := STRING | PATH | LEFT_PAREN strPathExpr RIGHT_PAREN
```
### Tokens

```
PATH := must start with $
NUMBER := TODO
BOOL := true | false
STRING := '[any characters, with escaped ' and \]'
IF_NOT_FOUND := ?
MINUS := -
PLUS := +
SUM := sum
LEFT_PAREN := (
RIGHT_PAREN := )
TIMES := *
PRODUCT := product
DIVIDE := /
LENGTH := length
LESS := <
LESS_EQ := <=
MORE := >
MORE_EQ := >=
EQ := ==
AND_OP := &&
AND := and
OR_OP := ||
OR := or
COMMA := ,
```

## Potential Additions

Potential additions:

* Array literals, i.e. `['thing', 'other thing']`
* Array inclusion binary operator `in`, i.e. `$.questionType in ['radio', 'shortAnswer']`
* String concatenation function `concat`, i.e. `concat($.s1, $.s2)`
