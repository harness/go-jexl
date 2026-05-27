# go-jexl

A JEXL-compliant expression and script evaluator for Go. Compiles expressions to bytecode and runs them on a stack-based virtual machine.

```go
package main

import (
    "fmt"

    "github.com/harness/go-jexl/jexl"
)

type User struct {
    Name string
    Age  int
}

func main() {
    context := User{
        Name: "Alice",
        Age:   30,
    }

    script := `Name + " is " + Age`
    program, err := jexl.Compile(script,
        jexl.WithContext(context),
        jexl.WithMaxIterations(1000), // optional
        jexl.WithMaxMemory(1000),     // optional
        jexl.WithMaxNodes(100),       // optional
    )
    if err != nil {
        panic(err)
    }

    output, err := jexl.Run(program, env)
    if err != nil {
        panic(err)
    }

    fmt.Println(output) // Alice is 30
}
```

## Options

| Option                              | Description                                        |
|-------------------------------------|----------------------------------------------------|
| `WithContext(v)`                    | Register variables to the program                  |
| `WithClass(name, obj)`              | Register a class for `new`                         |
| `WithNamespace(name, obj)`          | Register a namespace object                        |
| `WithFunction(name, fn, types...)`  | Register custom functions                          |
| `WithMaxIterations(n)`              | Cap loop iterations (0 = unlimited)                |
| `WithMaxNodes(n)`                   | Cap ast node count at compile time (0 = unlimited) |
| `WithMaxMemory(n)`                  | Cap mv memory usage (0 = unlimited)                |
| `WithStrict()`                      | Error                                              |
| `WithSafe()`                        | Prevent expressions from mutating the context      |

## Errors

Compile errors include a source snippet with a caret pointing at the problem:

```
invalid operation: + (mismatched types string and int) (1:8)
 | name + age
 | .......^
```

## Fast Path

For simple proprty access we have a zero-allocation fast evaluator. Use it when you know the expression is a simple property chain.

```go
out, err := jexl.EvalPath("trigger.payload.crNumber", env)
```

The function returns sentinel error `ErrNotPropertyPath`if the expression is not a plain property path. Use this to fall through to `Eval` for mixed workloads:

```go
out, err := jexl.EvalPath(expr, env)
if errors.Is(err, jexl.ErrNotPropertyPath) {
    out, err = jexl.Eval(expr, env)
}
```

## Templates

The `template` package evaluates JEXL expressions embedded in a string, replacing each `<+expr>` with its result.

```go

import (
    "fmt"

    "github.com/harness/go-jexl/jexl"
    "github.com/harness/go-jexl/jexl/template"
)

type User struct {
    Name string
    Age  int
}

func main() {
    context := User{
        Name: "Alice",
        Age:   30,
    }

    tmpl := template.New(func(expr string) (any, error) {
        return jexl.Eval(expr, context)
    })

    output, err := tmpl.ExecString("hello <+name>")
    if err != nil {
        panic(err)
    }

    fmt.Println(output) // hello Alice
}
```

Note that the delimer can be customized:

```go
tmpl = tmpl.Delim("${{", "}}")
out, _ := tmpl.ExecString("hello ${{ name }}")
// out == "hello Alice"
```


## Syntax Reference

### Literals

| Form         | Example                            |
|--------------|------------------------------------|
| Integer      | `42`, `-7`                         |
| Float        | `3.14`, `-0.5`                     |
| Long         | `42l`, `42L`                       |
| Double       | `42.0d`, `42.0D`, `NaN`            |
| BigInteger   | `42h`, `42H`                       |
| BigDecimal   | `3.14b`, `3.14B`                   |
| Octal        | `010`                              |
| Hex          | `0x10`, `0X10`                     |
| Exponent     | `42.0E-1D`, `42.0E+3B`             |
| String       | `"hello"`, `'world'`               |
| Boolean      | `true`, `false`                    |
| Null         | `null`                             |
| Array        | `[1, 2, 3]`                        |
| List         | `[1, 2, 3,...]`                    |
| Map          | `{"key": "value"}`                 |
| Empty map    | `{:}`                              |
| Set          | `{1, 2, 3}`                        |
| Regex        | `~/pattern/`                       |

Unicode escapes are supported in strings: `'a'` == `'a'`.

### Arithmetic

```js
1 + 2        // 3
10 - 4       // 6
3 * 4        // 12
7 / 2        // 3.5  (always float64)
7 div 2      // same as /
7 % 3        // 1
5 mod 2      // same as %
2 ** 10      // 1024.0
-a           // negate
+a           // positivize (promotes byte/short/char to int)
```

### Comparison

```js
a == b     // also: a eq b
a != b     // also: a ne b
a === b    // strict (no type coercion)
a !== b    // strict (no type coercion)
a < b      // also: a lt b
a <= b     // also: a le b
a > b      // also: a gt b
a >= b     // also: a ge b
```

### Logical

```js
a && b      // also: a and b
a || b      // also: a or b
!a          // also: not a
a ?? b      // b when a is null (does not coerce to boolean — false ?? "x" == false)
a ?: b      // b when a is false or null
```

### Bitwise

```js
a | b     // bitwise OR
a ^ b     // bitwise XOR
a & b     // bitwise AND
~a        // bitwise NOT
a << b    // left shift
a >> b    // signed right shift
a >>> b   // unsigned right shift
```

### Operator precedence (low → high)

| Precedence | Operators |
|------------|-----------|
| 10         | `\|\|` `or` |
| 15         | `&&` `and` |
| 20         | `==` `!=` `===` `!==` `<` `>` `<=` `>=` `in` `=~` `!~` `=^` `!^` `=$` `!$` |
| 22         | `\|` |
| 24         | `^` |
| 26         | `&` |
| 28         | `<<` `>>` `>>>` |
| 30         | `+` `-` |
| 60         | `*` `/` `%` |
| 100        | `**` (right-associative) |
| 500        | `??` |

### Ternary and conditional

Ternary conditions

```js
condition ? thenValue : elseValue
```

If / else condition blocks:

```js
if (condition) { value } else if (other) { value2 } else { fallback }
```

### Member access and safe navigation

```js
user.name             // field or map key access
user?.name            // optional chaining — null if user is null
arr[0]                // index access
arr.0                 // dot-numeric — same as arr[0]
arr?.[0]              // optional index
foo.'bar'             // quoted property — for names with spaces or reserved words
foo['new']            // bracket form for reserved word keys
foo.`${bar}`          // back-quoted interpolated property name
```

### Template Literals

```js
`Hello ${name}`           // string interpolation
`${a} + ${b} = ${a + b}`  // expressions inside ${}
```

### String Operators

```js
"hello" =^ "he"    // startsWith — true
"hello" =$ "lo"    // endsWith — true
"hello" =~ "hel.*" // matches — true
"hello" !^ "he"    // not startsWith
"hello" !$ "lo"    // not endsWith
"hello" !~ "xyz"   // not matches
"hello" =~ /^h/    // regex match
"hello" =~ "^h"    // string pattern match
```

### Membership

```js
"admin" in ["admin", "mod"]   // true (array membership)
"admin" in {"admin": true}    // true (map key membership)
1 =~ [1, 2, 3]                // true (alternate membership syntax)
```

### Range

```js
1..5    // [1, 2, 3, 4, 5]
```

### `instanceof`

```js
x instanceof Integer
x !instanceof String
```

### Builtins

```js
empty(x)    // true for null, "", 0, empty collection
empty x     // prefix form, same as empty(x)
size(x)     // length of string, array, slice, or map
size x      // prefix form, same as size(x)
```

### Variables

```js
var x = 10       // script-wide scope, can be redefined
let y = x + 5    // block-scoped, no redeclaration in same scope
const PI = 3.14  // block-scoped, prevents reassignment
```

### Assignment

```js
x = x * 2
x += 1
x -= 1
x *= 2
x /= 2
x %= 3
x &= y
x |= y
x ^= y
x <<= y
x >>= y
x >>>= y
x++
x--
++x
--x
```

### Loops

```js
// For each loop
for (var item : collection) { ... }
for (i : 1..10) { ... }

// For loop
for (var i = 0; i < n; i++) { ... }

// While loop
while (condition) { ... }

// Do while loop
do { ... } while (condition)
```

### Control flow

```js
break
continue
return expr
```

### try / catch / finally / throw

```js
try {
    throw "oops"
} catch (let e) {
    e   // "oops"
} finally {
    // always runs
}
```

### switch

Arrow form (expression-style, returns a value):

```js
var r = switch (n) {
    case 1 -> "one"
    case 2, 3 -> "two or three"
    default -> "other"
}
```

Colon form (statement-style, uses `break`):

```js
switch (n) {
    case 1: r = "one"; break
    case 2: r = "two"; break
    default: r = "other"
}
```

Fall-through with colon form (empty case falls to next):

```js
switch (n) {
    case 1:
    case 2: r = "one or two"; break
    default: r = "other"
}
```

### Functions

```js
function add(x, y) { x + y }
add(3, 4)   // 7

// anonymous function expression
var add = function(x, y) { x + y }
```

## Lambdas

```js
(x) -> x * 2              // thin arrow
(x) => x * 2              // fat arrow
x -> x > 0                // single-arg shorthand
(x, y) -> x + y

(x) -> { var y = x * 2; y + 1 }   // block body

// parameter qualifiers
(let x, let y) -> x + y
(const x, const y) -> x + y
```

Lambdas are first-class values:

```js
var double = (x) -> x * 2
double(5)   // 10

var ops = {"double": (x) -> x * 2, "square": (x) -> x * x}
ops["double"](4)   // 8
```

Closures capture variable **values at definition time**:

```js
var t = 20
var s = function(x, y) { x + y + t }
t = 54
s(15, 7)   // 42 — t was captured as 20, not 54
```

## Builtins

Global builtin functions available to your expressions.

```js
abs(-5)              // 5
toUpperCase("hello") // "HELLO"
jsonMarshal({a: 1})  // '{"a":1}'
currentDate()        // "2024-01-15"
```

### Math

| Function           | Description                              |
|--------------------|------------------------------------------|
| `abs(v)`           | Absolute value                           |
| `ceil(v)`          | Smallest integer >= v                    |
| `floor(v)`         | Largest integer <= v                     |
| `round(v)`         | Nearest integer                          |
| `sqrt(v)`          | Square root                              |
| `pow(base, exp)`   | base ^ exp                               |
| `log(v)`           | Natural logarithm                        |
| `log2(v)`          | Binary logarithm                         |
| `log10(v)`         | Decimal logarithm                        |
| `min(a, b)`        | Smaller of a and b                       |
| `max(a, b)`        | Larger of a and b                        |
| `isNaN(v)`         | True if v is NaN                         |
| `isInfinite(v)`    | True if v is positive or negative infinity |

### String / Text

| Function                         | Description                                               |
|----------------------------------|-----------------------------------------------------------|
| `charAt(v, i)`                   | Unicode code point at index i                             |
| `codePointAt(v, i)`              | Code point value at index i                               |
| `compareTo(v, other)`            | Ordering of v vs other (-1, 0, 1)                         |
| `compareToIgnoreCase(v, other)`  | Case-insensitive ordering                                 |
| `concat(v, args...)`             | Concatenate v with each arg                               |
| `contains(v, sub)`               | True if v contains sub                                    |
| `endsWith(v, suffix)`            | True if v ends with suffix                                |
| `equals(v, other)`               | String equality                                           |
| `equalsIgnoreCase(v, other)`     | Case-insensitive equality                                 |
| `formatted(v, args...)`          | `fmt.Sprintf(v, args...)`                                 |
| `indexOf(v, sub)`                | First index of sub, or -1                                 |
| `indexOf(v, sub, from)`          | First index starting at from, or -1                       |
| `isBlank(v)`                     | True if v is empty or whitespace-only                     |
| `isEmpty(v)`                     | True if v is the empty string                             |
| `lastIndexOf(v, sub)`            | Last index of sub, or -1                                  |
| `length(v)`                      | Byte length                                               |
| `matches(v, pattern)`            | True if v fully matches regex pattern                     |
| `repeat(v, n)`                   | v repeated n times                                        |
| `replace(v, old, new)`           | All occurrences of old replaced by new                    |
| `replaceAll(v, pattern, repl)`   | All regex matches replaced by repl                        |
| `replaceFirst(v, pattern, repl)` | First regex match replaced by repl                        |
| `split(v, sep)`                  | Split v by sep                                            |
| `startsWith(v, prefix)`          | True if v starts with prefix                              |
| `strip(v)`                       | Leading and trailing whitespace removed                   |
| `stripLeading(v)`                | Leading whitespace removed                                |
| `stripTrailing(v)`               | Trailing whitespace removed                               |
| `substring(v, start)`            | v[start:]                                                 |
| `substring(v, start, end)`       | v[start:end]                                              |
| `substringBefore(v, delim)`      | Portion of v before the first occurrence of delim         |
| `toCharArray(v)`                 | Each character as a single-character string               |
| `toLowerCase(v)`                 | Lowercase                                                 |
| `toUpperCase(v)`                 | Uppercase — aliases: `toUpper`, `upper`                   |
| `trim(v)`                        | Leading and trailing whitespace removed                   |
| `trimLeft(v)`                    | Leading whitespace removed                                |
| `trimRight(v)`                   | Trailing whitespace removed                               |

### Type Conversion

| Function          | Description                                              |
|-------------------|----------------------------------------------------------|
| `booleanValue(v)` | Coerce to bool                                           |
| `byteValue(v)`    | Coerce to int8                                           |
| `shortValue(v)`   | Coerce to int16                                          |
| `intValue(v)`     | Coerce to int — alias: `toInteger`                       |
| `longValue(v)`    | Coerce to int64                                          |
| `floatValue(v)`   | Coerce to float32                                        |
| `doubleValue(v)`  | Coerce to float64                                        |
| `toString(v)`     | Coerce to string                                         |
| `default(v, fb)`  | Return v if non-nil and non-zero, otherwise fb           |

### Date / Time

| Function                  | Description                                              |
|---------------------------|----------------------------------------------------------|
| `currentDate()`           | Today's date as `YYYY-MM-DD`                             |
| `currentTime()`           | Current time as `HH:MM:SS`                               |
| `dateFormat(t, pattern)`  | Format time using a Java-style pattern (e.g. `yyyy-MM-dd HH:mm:ss`) |
| `plusMinutes(t, n)`       | Add n minutes to t                                       |
| `plusHours(t, n)`         | Add n hours to t                                         |
| `plusDays(t, n)`          | Add n days to t                                          |

### Encoding

| Function                    | Description                              |
|-----------------------------|------------------------------------------|
| `base64Encode(v)`           | Standard base64 encoding                 |
| `base64Decode(v)`           | Standard base64 decoding                 |
| `base64URLEncode(v)`        | URL-safe base64 encoding                 |
| `base64URLDecode(v)`        | URL-safe base64 decoding                 |
| `base64RawEncode(v)`        | Unpadded standard base64 encoding        |
| `base64RawDecode(v)`        | Unpadded standard base64 decoding        |
| `hexEncode(v)`              | Hexadecimal encoding                     |
| `hexDecode(v)`              | Hexadecimal decoding                     |

### Formatting

| Function                   | Description                              |
|----------------------------|------------------------------------------|
| `sprintf(format, args...)` | Formatted string (`fmt.Sprintf`)         |
| `sprint(args...)`          | Default string representation of args   |

### JSON

| Function                         | Description                                                          |
|----------------------------------|----------------------------------------------------------------------|
| `jsonMarshal(v)`                 | JSON encoding of v                                                   |
| `jsonUnmarshal(v)`               | Parse JSON string into a map                                         |
| `jsonMarshalIndent(v, pre, ind)` | Indented JSON encoding                                               |
| `jsonSelect(path, v)`            | Extract a value from a JSON string using a gjson/JSONPath expression |
| `jsonList(path, v)`              | Extract a list from a JSON string using a gjson/JSONPath expression  |

### Regex

| Function                   | Description                                              |
|----------------------------|----------------------------------------------------------|
| `regexExtract(pattern, v)` | First capture group from v matching pattern, or `""`     |

### XML

| Function                        | Description                                              |
|---------------------------------|----------------------------------------------------------|
| `xmlMarshal(v)`                 | XML encoding of v                                        |
| `xmlMarshalIndent(v, pre, ind)` | Indented XML encoding                                    |
| `xmlSelect(expr, doc)`          | Text content of all nodes matching the XPath expression  |


## Pragmas

Pragmas are directives at the top of an expression, parsed before compilation:

```sh
#pragma jexl.options 'strict'
#pragma jexl.import 'java.lang'
```

Supported pragma keys:

| Key | Values |
|-----|--------|
| `jexl.options` | `strict` / `+strict` (enable), `-strict` (disable) |
| `jexl.import` | `java.lang`, `java.math`, `java.util`, or a specific class |

## Namespaces

Namespaces let you call static methods and access constants using `Name:method(args)` syntax. Register any Go value as a namespace with `WithNamespace`:

```go
jexl.Compile(expr, jexl.WithNamespace("Math", javalang.NewMath()))
```

Namespaces are not added by default — every namespace must be explicitly registered.

## Java Classes

go-jexl emulates a subset of the Java standard library. These classes are not loaded automatically; use a pragma or register them with `WithClass`.

```sh
#pragma jexl.import 'java.lang'
#pragma jexl.import 'java.math'
#pragma jexl.import 'java.util'
```

Or register a specific class:

```go
jexl.Compile(expr,
    jexl.WithClass("StringBuilder", javalang.NewStringBuilder),
    jexl.WithClass("ArrayList", javautil.NewArrayList),
)
```

Constructable classes are instantiated with `new`. This works for emulated Java classes and any custom class registered via `WithClass`:

```js
var list = new('java.util.ArrayList')
var sb   = new('java.lang.StringBuilder', 'hello')
var pt   = new('Point', 3, 4)   // custom class registered with WithClass("Point", ...)
```

### java.lang.String

Static methods:

| Method                      | Description                 |
|-----------------------------|-----------------------------|
| `valueOf(v)`                | Coerce to string            |
| `format(fmt, args...)`      | `fmt.Sprintf(fmt, args...)` |

Instance methods:

| Method                              | Description                                      |
|-------------------------------------|--------------------------------------------------|
| `.charAt(i)`                        | Unicode code point at index i                    |
| `.codePointAt(i)`                   | Code point value at index i                      |
| `.compareTo(other)`                 | Ordering (-1, 0, 1)                              |
| `.compareToIgnoreCase(other)`       | Case-insensitive ordering                        |
| `.concat(args...)`                  | Concatenate                                      |
| `.contains(sub)`                    | True if sub is within the string                 |
| `.endsWith(suffix)`                 | True if ends with suffix                         |
| `.equals(other)`                    | String equality                                  |
| `.equalsIgnoreCase(other)`          | Case-insensitive equality                        |
| `.formatted(args...)`               | `fmt.Sprintf(v, args...)`                        |
| `.indexOf(sub)`                     | First index of sub, or -1                        |
| `.indexOf(sub, from)`               | First index starting at from, or -1              |
| `.isBlank()`                        | True if empty or whitespace-only                 |
| `.isEmpty()`                        | True if zero length                              |
| `.lastIndexOf(sub)`                 | Last index of sub, or -1                         |
| `.length()`                         | Byte length                                      |
| `.matches(pattern)`                 | True if fully matches regex                      |
| `.repeat(n)`                        | String repeated n times                          |
| `.replace(old, new)`                | All occurrences of old replaced by new           |
| `.replaceAll(pattern, repl)`        | All regex matches replaced                       |
| `.replaceFirst(pattern, repl)`      | First regex match replaced                       |
| `.split(sep)`                       | Split by sep                                     |
| `.startsWith(prefix)`               | True if starts with prefix                       |
| `.strip()`                          | Leading and trailing whitespace removed          |
| `.stripLeading()`                   | Leading whitespace removed                       |
| `.stripTrailing()`                  | Trailing whitespace removed                      |
| `.substring(from)`                  | v[from:]                                         |
| `.substring(from, to)`              | v[from:to]                                       |
| `.substringBefore(delim)`           | Portion before first occurrence of delim         |
| `.toCharArray()`                    | Each character as a single-char string           |
| `.toLowerCase()`                    | Lowercase                                        |
| `.toUpperCase()`                    | Uppercase                                        |
| `.trim()`                           | Leading and trailing whitespace removed          |

### java.lang.StringBuilder / java.lang.StringBuffer

Constructable with `new`. Instance methods:

| Method                     | Description                        |
|----------------------------|------------------------------------|
| `.append(v)`               | Append value, returns `this`       |
| `.insert(i, v)`            | Insert at index, returns `this`    |
| `.delete(start, end)`      | Remove range, returns `this`       |
| `.deleteCharAt(i)`         | Remove one char, returns `this`    |
| `.replace(start, end, s)`  | Replace range, returns `this`      |
| `.reverse()`               | Reverse content, returns `this`    |
| `.charAt(i)`               | Code point at index                |
| `.indexOf(sub)`            | Index of substring, or -1          |
| `.length()`                | Code point count                   |
| `.substring(from, to)`     | Substring                          |
| `.toString()`              | Get string                         |

### java.math.BigDecimal

Arbitrary-precision decimal. Also available via the `3.14B` literal suffix.

Static methods:

| Method        | Description                     |
|---------------|---------------------------------|
| `new(v)`      | Construct from number or string |
| `valueOf(v)`  | Same as new                     |
| `ZERO`        | 0                               |
| `ONE`         | 1                               |
| `TEN`         | 10                              |

Instance methods:

| Method                   | Description                              |
|--------------------------|------------------------------------------|
| `.add(o)`                | Addition                                 |
| `.subtract(o)`           | Subtraction                              |
| `.multiply(o)`           | Multiplication                           |
| `.divide(o)`             | Division                                 |
| `.remainder(o)`          | Remainder                                |
| `.pow(exp)`              | Exponentiation                           |
| `.negate()`              | Negation                                 |
| `.abs()`                 | Absolute value                           |
| `.max(o)`                | Larger of this and o                     |
| `.min(o)`                | Smaller of this and o                    |
| `.compareTo(o)`          | -1, 0, or 1                              |
| `.equals(o)`             | Equality                                 |
| `.scale()`               | Scale                                    |
| `.setScale(n)`           | Set scale                                |
| `.stripTrailingZeros()`  | Remove trailing zeros                    |
| `.signum()`              | -1, 0, or 1                              |
| `.intValue()`            | Coerce to int                            |
| `.longValue()`           | Coerce to int64                          |
| `.doubleValue()`         | Coerce to float64                        |
| `.floatValue()`          | Coerce to float32                        |
| `.toString()`            | Decimal string                           |
| `.toPlainString()`       | Non-scientific decimal string            |

### java.math.BigInteger

Arbitrary-precision integer. Also available via the `42H` literal suffix.

Static methods:

| Method           | Description                         |
|------------------|-------------------------------------|
| `new(v)`         | Construct from number or string     |
| `new(s, base)`   | Construct from string in given base |
| `valueOf(v)`     | Construct from int64                |
| `ZERO`           | 0                                   |
| `ONE`            | 1                                   |
| `TWO`            | 2                                   |
| `TEN`            | 10                                  |

Instance methods:

| Method            | Description                |
|-------------------|----------------------------|
| `.add(o)`         | Addition                   |
| `.subtract(o)`    | Subtraction                |
| `.multiply(o)`    | Multiplication             |
| `.divide(o)`      | Division                   |
| `.remainder(o)`   | Remainder                  |
| `.mod(o)`         | Modulus (always positive)  |
| `.pow(exp)`       | Exponentiation             |
| `.negate()`       | Negation                   |
| `.abs()`          | Absolute value             |
| `.and(o)`         | Bitwise AND                |
| `.or(o)`          | Bitwise OR                 |
| `.xor(o)`         | Bitwise XOR                |
| `.not()`          | Bitwise NOT                |
| `.shiftLeft(n)`   | Left shift                 |
| `.shiftRight(n)`  | Right shift                |
| `.compareTo(o)`   | -1, 0, or 1                |
| `.equals(o)`      | Equality                   |
| `.max(o)`         | Larger of this and o       |
| `.min(o)`         | Smaller of this and o      |
| `.bitLength()`    | Number of bits             |
| `.signum()`       | -1, 0, or 1                |
| `.intValue()`     | Coerce to int              |
| `.longValue()`    | Coerce to int64            |
| `.doubleValue()`  | Coerce to float64          |
| `.toString()`     | Decimal string             |
| `.toString(base)` | String in given base       |

### java.lang.Integer

| Method                | Description                |
|-----------------------|----------------------------|
| `parseInt(s)`         | Parse decimal string       |
| `parseInt(s, base)`   | Parse string in given base |
| `valueOf(v)`          | Construct from value       |
| `toString(v)`         | Decimal string             |
| `toString(v, base)`   | String in given base       |
| `compare(a, b)`       | -1, 0, or 1                |
| `max(a, b)`           | Larger value               |
| `min(a, b)`           | Smaller value              |
| `sum(a, b)`           | a + b                      |
| `toBinaryString(v)`   | Binary string              |
| `toHexString(v)`      | Hex string                 |
| `toOctalString(v)`    | Octal string               |
| `MAX_VALUE`           | 2147483647                 |
| `MIN_VALUE`           | -2147483648                |

### java.lang.Long

| Method              | Description                |
|---------------------|----------------------------|
| `parseLong(s)`      | Parse decimal string       |
| `parseLong(s, base)`| Parse string in given base |
| `valueOf(v)`        | Construct from value       |
| `toString(v)`       | Decimal string             |
| `toString(v, base)` | String in given base       |
| `compare(a, b)`     | -1, 0, or 1                |
| `max(a, b)`         | Larger value               |
| `min(a, b)`         | Smaller value              |
| `toBinaryString(v)` | Binary string              |
| `toHexString(v)`    | Hex string                 |
| `MAX_VALUE`         | 9223372036854775807        |
| `MIN_VALUE`         | -9223372036854775808       |

### java.lang.Double

| Method               | Description                 |
|----------------------|-----------------------------|
| `parseDouble(v)`     | Construct from value        |
| `valueOf(v)`         | Same as parseDouble         |
| `toString(v)`        | String representation       |
| `compare(a, b)`      | -1, 0, or 1                 |
| `max(a, b)`          | Larger value                |
| `min(a, b)`          | Smaller value               |
| `sum(a, b)`          | a + b                       |
| `isNaN(v)`           | True if NaN                 |
| `isInfinite(v)`      | True if infinite            |
| `isFinite(v)`        | True if finite              |
| `MAX_VALUE`          | math.MaxFloat64             |
| `MIN_VALUE`          | math.SmallestNonzeroFloat64 |
| `POSITIVE_INFINITY`  | +Inf                        |
| `NEGATIVE_INFINITY`  | -Inf                        |
| `NaN`                | NaN                         |

### java.lang.Float

| Method               | Description                  |
|----------------------|------------------------------|
| `parseFloat(v)`      | Construct from value         |
| `valueOf(v)`         | Same as parseFloat           |
| `toString(v)`        | String representation        |
| `compare(a, b)`      | -1, 0, or 1                  |
| `max(a, b)`          | Larger value                 |
| `min(a, b)`          | Smaller value                |
| `sum(a, b)`          | a + b                        |
| `isNaN(v)`           | True if NaN                  |
| `isInfinite(v)`      | True if infinite             |
| `isFinite(v)`        | True if finite               |
| `MAX_VALUE`          | math.MaxFloat32              |
| `MIN_VALUE`          | math.SmallestNonzeroFloat32  |
| `POSITIVE_INFINITY`  | +Inf                         |
| `NEGATIVE_INFINITY`  | -Inf                         |
| `NaN`                | NaN                          |

### java.lang.Boolean

| Method               | Description           |
|----------------------|-----------------------|
| `parseBoolean(v)`    | Construct from value  |
| `valueOf(v)`         | Same as parseBoolean  |
| `toString(v)`        | `"true"` or `"false"` |
| `compare(a, b)`      | -1, 0, or 1           |
| `logicalAnd(a, b)`   | a && b                |
| `logicalOr(a, b)`    | a \|\| b              |
| `logicalXor(a, b)`   | a != b                |
| `TRUE`               | true constant         |
| `FALSE`              | false constant        |

### java.lang.Character

| Method               | Description                   |
|----------------------|-------------------------------|
| `valueOf(v)`         | Construct from rune or string |
| `isDigit(v)`         | True if digit                 |
| `isLetter(v)`        | True if letter                |
| `isLetterOrDigit(v)` | True if letter or digit       |
| `isUpperCase(v)`     | True if uppercase             |
| `isLowerCase(v)`     | True if lowercase             |
| `isWhitespace(v)`    | True if whitespace            |
| `toUpperCase(v)`     | Uppercase version             |
| `toLowerCase(v)`     | Lowercase version             |
| `toString(v)`        | Single-character string       |
| `compare(a, b)`      | -1, 0, or 1                   |
| `MAX_VALUE`          | Unicode max rune (0x10FFFF)   |
| `MIN_VALUE`          | Null character (0)            |

### java.util.ArrayList

Ordered list, constructable with `new`. Array and list literals `[1, 2, 3]` also support these methods.

| Method              | Description                              |
|---------------------|------------------------------------------|
| `.size()`           | Length                                   |
| `.isEmpty()`        | True if empty                            |
| `.contains(x)`      | True if x is in the list                 |
| `.get(i)`           | Item at index                            |
| `.set(i, v)`        | Replace item, returns old value          |
| `.indexOf(x)`       | Index of x, or -1                        |
| `.add(x)`           | Append                                   |
| `.add(i, x)`        | Insert at index                          |
| `.remove(i)`        | Remove at index, returns old value       |
| `.clear()`          | Remove all                               |
| `.addAll(other)`    | Append all from other list               |
| `.subList(from, to)`| Slice [from, to)                         |
| `.toArray()`        | Copy as slice                            |
| `.addFirst(x)`      | Insert at beginning                      |
| `.addLast(x)`       | Append at end                            |
| `.getFirst()`       | First item                               |
| `.getLast()`        | Last item                                |
| `.peek()`           | First item or nil                        |
| `.poll()`           | Remove and return first item             |
| `.toString()`       | `"[a, b, c]"`                            |

### java.util.HashMap

Key/value map, constructable with `new`. Map literals `{"key": "value"}` also support these methods.

| Method                    | Description                              |
|---------------------------|------------------------------------------|
| `.size()`                 | Number of entries                        |
| `.isEmpty()`              | True if empty                            |
| `.containsKey(k)`         | True if key exists                       |
| `.containsValue(v)`       | True if value exists                     |
| `.get(k)`                 | Value or nil                             |
| `.put(k, v)`              | Set key, returns old value               |
| `.remove(k)`              | Remove key, returns old value            |
| `.clear()`                | Remove all entries                       |
| `.keySet()`               | Array of keys                            |
| `.values()`               | Array of values                          |
| `.entrySet()`             | Array of `{"key":k,"value":v}` maps      |
| `.getOrDefault(k, def)`   | Value or default                         |
| `.putIfAbsent(k, v)`      | Set if absent, returns current value     |
| `.putAll(other)`          | Copy all from other map                  |
| `.toString()`             | `"{k=v, ...}"`                           |

### java.util.Stack

LIFO stack, constructable with `new`.

| Method        | Description                                  |
|---------------|----------------------------------------------|
| `.push(x)`    | Push and return x                            |
| `.pop()`      | Remove and return top (error if empty)       |
| `.peek()`     | Top without removing (error if empty)        |
| `.empty()`    | True if empty                                |
| `.size()`     | Number of items                              |
| `.search(x)`  | 1-based position from top, or -1            |

### java.util.HashSet

Unique-element set, constructable with `new`.

| Method          | Description                          |
|-----------------|--------------------------------------|
| `.add(x)`       | Insert, returns true if new          |
| `.remove(x)`    | Delete, returns true if was present  |
| `.contains(x)`  | True if member                       |
| `.size()`       | Number of elements                   |
| `.isEmpty()`    | True if empty                        |
| `.clear()`      | Remove all                           |
| `.toArray()`    | Elements in insertion order          |
| `.toString()`   | `"[a, b, c]"`                        |
