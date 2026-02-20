# decodini

<div align="center">
    <img src="https://img.shields.io/badge/Written_In-Go-00acd7?style=for-the-badge&logo=go" alt="Go" />
</div>

<div align="center">
    <img width="300" src="/assets/decodini.png" alt="decodini" />
</div>

<br />

Decodini is a small Go library for moving data between structs, maps, slices, and primitives.

It works in two steps internally: encode the source into a `Tree`, then decode that tree into the target type. Most users can just call `Transmute`.

## Installation

```bash
go get github.com/lukasl-dev/decodini
```

## Usage

The primary interface for the library is the `Transmute` function, which performs both encoding and decoding in a single step.

### Basic Example

```go
package main

import (
	"fmt"
	"github.com/lukasl-dev/decodini/pkg/decodini"
)

type UserSource struct {
	Username string
	Age      int
}

type UserTarget struct {
	Username string
	Age      int
}

func main() {
	src := UserSource{Username: "alice", Age: 30}
	
	// Transmute source into UserTarget
	dst, err := decodini.Transmute[UserTarget](nil, src)
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("%+v\n", dst)
}
```

### Struct Tags

Use the `decodini` struct tag to map fields that have different names in their respective structs. For a successful transmutation, the fields must resolve to the same name.

```go
type Source struct {
	InternalName string `decodini:"name"`
}

type Target struct {
	ExternalName string `decodini:"name"`
	IgnoredField string `decodini:"-"`
}
```

### Transmuting into Existing Values

Use `TransmuteInto` to populate an existing variable.

```go
var target UserTarget
err := decodini.TransmuteInto(nil, src, &target)
```

## Advanced Configuration

The `Transmutation` struct allows for customisation of the encoding and decoding behaviour.

```go
tm := &decodini.Transmutation{
	Decoding: &decodini.Decoding{
		StructTag: "custom_tag",
	},
}

dst, err := decodini.Transmute[UserTarget](tm, src)
```

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
