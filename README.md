# go-lck

Package **lck** implements thread-safe locking-types, for the Go programming language.

## Documention

Online documentation, which includes examples, can be found at: http://godoc.org/github.com/reiver/go-lck

[![GoDoc](https://godoc.org/github.com/reiver/go-lck?status.svg)](https://godoc.org/github.com/reiver/go-lck)

## Examples

Here is an example **locking-type** that can hold a `bool`:

```go
import "github.com/reiver/go-lck"

//

var lockable lck.Locking[bool]

// ...

lockable.Set(true)

// ...

value := lockable.Get()
```

Here is another example  **locking-type**, where here it holds a `map`:

```go
import "github.com/reiver/go-lck"

//

var lockable lck.Locking[map[string]any]

// ...

lockable.Let(fn(m *map[string]any) {
	if nil == *m {
		*m = map[string]any{}
	}

	*m["something"] = 5
})

// ...

var value any

lockable.Let(fn(m *map[string]any) {
	if nil == *m {
		return
	}

	value = *m["something"]
})
```

## Import

To import package **lck** use `import` code like the follownig:
```
import "github.com/reiver/go-lck"
```

## Installation

To install package **lck** do the following:
```
GOPROXY=direct go get https://github.com/reiver/go-lck
```

## Author

Package **lck** was written by [Charles Iliya Krempeaux](http://reiver.link)
