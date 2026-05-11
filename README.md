# GoAL — Go Algorithms Library

GoAL is a generic algorithms and data structures library for Go, developed at Hochschule München. It provides well-tested, interface-driven implementations suitable for algorithm education and applied use.

GoAL aims at broad coverage of classic algorithms and data structures, with a special emphasis on graph algorithms. These are built on a mixed graph model — a graph that may simultaneously contain directed arcs (with a source and a target) and undirected edges (connecting two vertices without a direction). This generalizes both directed and undirected graphs and allows algorithms to operate naturally on graphs that mix both connector types.

Refer to the [extensive documentation](https://hm.pages.gitlab.lrz.de/goal) for usage guides, API references, and code examples for all packages.

## Requirements

Go 1.21 or later.

## Installation

```
go get gitlab.lrz.de/hm/goal
```