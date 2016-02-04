# goprox
A simple REST Proxy written in golang

#### Example

The simplest way to use GoProx is to route all requests to the REST endpoint.

```go
package main

import (
	"os"

	"github.com/damnpoet/goprox"
)

const (
	restURL = "http://sample-api.com/"
)

func main() {
	user, passwd := os.Getenv("OS_USERNAME"), os.Getenv("OS_PASSWORD")

	proxy := goprox.New(restURL, user, passwd)
	proxy.Listen(":8080")
}

```
