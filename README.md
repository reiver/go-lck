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
