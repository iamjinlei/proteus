# Proteus

Proteus is a tool to generate static site from markdown.
The focus is on making note taking easy with some custom `tags` to support keywords, highlight, etc.

## Setup

No setup is required.
Just clone the repo.
```
https://github.com/iamjinlei/proteus.git
```

## Run

Local run:

```
go run ./cmd/ -s [path to markdown root]
```

This converts and serves HTML requests directly from markdown source on demand.
This is good for editing and testing.

Generation:

```
go run ./cmd/ -s [path to source markdown root] -d [path to destination root] -g
```

## Tags

TODO
