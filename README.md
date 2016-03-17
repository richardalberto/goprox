# goprox [![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/damnpoet/goprox) [![Travis CI](https://travis-ci.org/damnpoet/goprox.svg?branch=master)](https://travis-ci.org/damnpoet/goprox) [![Coverage](http://gocover.io/_badge/github.com/damnpoet/goprox)](http://gocover.io/github.com/damnpoet/goprox)
A simple reverse proxy middleware written in golang

#### Example

The simplest way to use GoProx is to use it as the sole handler for all requests.

```go
package main

import (
	"os"
	"net/http"
	"log"

	"github.com/damnpoet/goprox"
)

func main() {
	mux := http.NewServeMux()

	proxy := goprox.New("http://foobar.com", goprox.Options{})
	handler := proxy.Handler(mux)

	log.Fatal(http.ListenAndServe(":8080", handler))
}

```
