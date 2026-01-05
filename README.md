# decodini

<div align="center">
    <img src="https://img.shields.io/badge/Written_In-Go-00acd7?style=for-the-badge&logo=go" alt="Go" />
</div>

<div align="center">
    <img width="300" src="/assets/decodini.png" alt="decodini" />
</div>

<br />

Decodini is a Go library for transmuting between different data structures. It uses an intermediate tree representation to convert between structs, maps, slices, and primitive types.

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

## Features

- **Direct Transmutation**: Convert between types without manual intermediate steps.
- **Recursive Transformation**: Full support for nested structures and collections.
- **Pointer Handling**: Automatic allocation and dereferencing of pointers.
- **Lazy Evaluation**: The intermediate `Tree` representation is evaluated only as needed.
- **Type Conversions**: Built-in logic for converting between strings and various slice types.

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

This project is licensed under the MIT License. See the `LICENSE` file for details.
