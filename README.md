# result

[![Go Reference](https://pkg.go.dev/badge/github.com/wlrgo/result/v2.svg)](https://pkg.go.dev/github.com/wlrgo/result/v2)
[![CI](https://github.com/wlrgo/result/actions/workflows/ci.yml/badge.svg)](https://github.com/wlrgo/result/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/wlrgo/result)](https://github.com/wlrgo/result/blob/main/go.mod)
[![Release](https://img.shields.io/github/v/release/wlrgo/result)](https://github.com/wlrgo/result/releases)
[![License](https://img.shields.io/github/license/wlrgo/result)](LICENSE)

`Result[T, E]` for success or failure.

## Rust API in Go

This is a [wlrgo](https://github.com/wlrgo) package. wlrgo ports Rust
standard-library types to Go and **intentionally copies the Rust API** instead
of reshaping it into an idiomatic Go design. Later packages in the org follow
the same rule.

Names, combinators, and eager vs lazy evaluation follow
[`std::result::Result`](https://doc.rust-lang.org/std/result/enum.Result.html)
as closely as Go generics allow. A few Go-only helpers exist where Rust has no
equivalent: `From`, `Get`, `GetErr`, and `Unpack`.

## Install

Requires Go 1.27 or later.

```bash
go get github.com/wlrgo/result/v2
```

## Example

```go
package main

import (
	"fmt"
	"strconv"

	"github.com/wlrgo/result/v2"
)

func main() {
	res := result.From(strconv.Atoi("2"))
	if v, ok := res.Get(); ok {
		fmt.Println(v)
	}

	fmt.Println(res.Map(func(n int) int { return n * 10 }).UnwrapOr(-1))
}
```

The zero value is `Err`. `Ok` and `Err` construct a result; `From` converts
from `(T, error)`. Combinators such as `Map`, `And`, and `Or` are methods.
`Compare`, `Flatten`, and `Collect` remain package-level functions.

## API

| Group | Highlights |
| --- | --- |
| Query | `IsOk`, `IsErr`, `IsOkAnd`, `IsErrAnd` |
| Extract | `Unwrap`, `Expect`, `UnwrapOr`, `Get` |
| Combine | `And`, `AndThen`, `Or`, `OrElse` |
| Transform | `Map`, `MapErr`, `Flatten`, `Inspect` |
| Compare | `Compare`, `Equal` (`Ok` < `Err`) |
| Iterate | `Seq`, `Collect` |
| Convert | `From`, `Get`, `GetErr`, `Unpack` |

See [pkg.go.dev/github.com/wlrgo/result/v2](https://pkg.go.dev/github.com/wlrgo/result/v2)
for the full API and package contract.

## License

[MIT](LICENSE)
