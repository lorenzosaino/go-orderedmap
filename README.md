# Go ordered map

[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/lorenzosaino/go-orderedmap)
[![Build](https://github.com/lorenzosaino/go-sysctl/workflows/Build/badge.svg)](https://github.com/lorenzosaino/go-orderedmap/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/lorenzosaino/go-sysctl)](https://goreportcard.com/report/github.com/lorenzosaino/go-orderedmap)
[![License](https://img.shields.io/github/license/lorenzosaino/go-sysctl.svg)](https://github.com/lorenzosaino/go-orderedmap/blob/master/LICENSE)

Go implementation of an ordered map using generics.

## Implementation

An ordered map is a map that additionally maintains ordering among its entries.
This data structure can be used to solve a variety of problems: one very common use case is implementing LRU or LRU-like cache replacement policies.

This implementation supports O(1) lookup, removal, insertion to front/back, insertion before/after a specific key, move to front/back, move before/after a specific key.

Under the hood this is implemented as a combination of a map and doubly-linked list, whereby each value in the map is node of the list.
The list is implemented by forking the standard library [`container/list`](https://pkg.go.dev/container/list) package and adding support for generics.

This implementation is not safe for concurrent usage.

## Installation

You can get this library by invoking

    go get -u github.com/lorenzosaino/go-orderedmap

Since it requires generics support, you will need Go 1.18 or above.

## Documentation

See [Go doc](https://pkg.go.dev/github.com/lorenzosaino/go-orderedmap?tab=doc).

## License

[BSD 3-clause](LICENSE)
