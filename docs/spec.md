# Overview

`texted` is a scriptable, headless text editor intended to be used for
making automated edits to files on disk.

## Using texted

A `texted` program is a series of instructions.
No means of abstraction are provided by `texted`,
these are to be implemented at a higher level, by the application using `texted`.

Much research has been done in this field already, so texted's editing language
is based on `emacs`.

### Example programs

Consider this file example.js:

```javascript
function doIt() {
  console.log("doing the thing")
}
```

We can turn it into this:

```javascript
function helloWorld() {
  console.log("hello, world")
}
```

Through the following `texted` script:

```sh
search-forward "doIt"
set-mark
search-forward "("
replace-region "helloWorld"
search-forward "doing the thing"
replace-match "hello, world"
```

This same script can also be encoded as an S-expression:

```lisp
(search-forward "doIt")
(set-mark)
(search-forward "(")
(replace-region "helloWorld")
(search-forward "doing the thing")
(replace-match "hello, world")
```

In fact, the regular texted script parser just uses a specialized reader:

1. Leading whitespace is stripped
2. If the next character is `(`, read a regular S-expression list,
   ignoring whitespace between elements.
3. Otherwise build a list reading invoking the reader repeatedly until
   a single `\n` is encountered.

#### JSON-encoding

The canonical format for encoding texted programs as JSON is this:

1. sexp-lists are JSON arrays,
2. sexp-numbers are JSON numbers,
3. sexp-strings are JSON strings,
4. sexp-symbols are JSON strings appearing as the first element of a list.
   Symbols are not allowed anywhere else.

According to these rules, the program above can be encoded as JSON like this:

```json
["search-forward", "doIt"]
["set-mark"]
["search-forward", "("]
["replace-region", "helloWorld"]
["search-forward", "doing the thing"]
["replace-match", "hello, world"]
```

## Architecture

The fundamental building block of texted is the `buffer`.

A buffer is a buffer of utf8-encoded text
unless otherwise specified by the buffer's `encoding`.

A buffer has a `point` (the current cursor position), starting at 1,
indicating the position *between* two characters.

The `mark` is "the other position" in a buffer,
which together with the `point` forms the `region`.

To obtain a buffer, call `texted.NewBuffer("initial contents")`

- `buf.String()` returns the contents of the buffer as a string
- `buf.Region()` returns the contents of the region
- `buf.PointMin()` returns the minimum value of the point (usually 1).
- `buf.PointMax()` returns the maximum value of the point.
- `buf.Mark()` returns the position of the mark.
- `buf.Encoding()` returns `"utf8"`
- `val, err := buf.Do(script)` executes script, returning the value of the last expression.

### Values

```go
package texted

type ValueKind interface {
  // a unique name identifying this kind
  KindName() string
}

type Value interface {
  Kind() ValueKind
}

// Example of a generic function operating on values and kinds
func IsA(value Value, kind Kind) bool {
  return value.Kind().KindName() == kind.KindName()
}

type SymbolKind struct {}

func (kind *SymbolKind) KindName() string { return "symbol" }

var TheSymbolKind = &SymbolKind{}

type Symbol struct {
  Name string 
}

func (sym *Symbol) Kind() ValueKind { 
  return TheSymbolKind
}

// same for numbers
// same for lists (backed by slices)
// same for strings
```
