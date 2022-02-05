# Go ordered map

[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/lorenzosaino/go-orderedmap)
[![Build](https://github.com/lorenzosaino/go-orderedmap/workflows/Build/badge.svg)](https://github.com/lorenzosaino/go-orderedmap/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/lorenzosaino/go-orderedmap)](https://goreportcard.com/report/github.com/lorenzosaino/go-orderedmap)
[![License](https://img.shields.io/github/license/lorenzosaino/go-orderedmap.svg)](https://github.com/lorenzosaino/go-orderedmap/blob/master/LICENSE)

Go implementation of an ordered map using generics.

## Implementation

An ordered map is a map that additionally maintains ordering among its entries.
This data structure can be used to solve a variety of problems: one very common use case is implementing LRU or LRU-like cache replacement policies.

This implementation supports O(1) lookup, removal, insertion to front/back, insertion before/after a specific key, move to front/back, move before/after a specific key.

Under the hood this is implemented as a combination of a map and doubly-linked list, whereby each value in the map is node of the list.
The list is implemented by forking the standard library [`container/list`](https://pkg.go.dev/container/list) package and adding support for generics.

This implementation is not safe for concurrent usage.

## Installation

You can get this library by invoking from your terminal

```
go get -u github.com/lorenzosaino/go-orderedmap
```

and then importing it in your code with

```go
import orderedmap "github.com/lorenzosaino/go-orderedmap"
```

Since it requires generics support, you will need Go 1.18 or above.

## Documentation

See [Go doc](https://pkg.go.dev/github.com/lorenzosaino/go-orderedmap?tab=doc).

## Development

You can invoke `make help` to see all make targets provided.

```bash
$ make help
all                            Run all checks and tests
mod-upgrade                    Upgrade all vendored dependencies
mod-update                     Ensure all used dependencies are tracked in go.{mod|sum} and vendored
fmt-check                      Validate that all source files pass "go fmt"
lint                           Run go lint
vet                            Run go vet
staticcheck                    Run staticcheck
test                           Run all tests
container-shell                Open a shell on a Docker container
container-%                    Run any target of this Makefile in a Docker container
help                           Print help
```

## License

[BSD 3-clause](LICENSE)
