# graceful
Package graceful is a Go 1.8+ package enabling graceful shutdown of http.Handler servers.

## Usage

Simply use `ListenAndServe` to create a http server:

```go
package main

import (
	"log"

	"github.com/ETZhangSX/graceful"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	if err := graceful.ListenAndServe(":8080", r.Handler()); err != nil {
		log.Fatal(err)
	}
}
```

The default timeout of shutdown is 5 seconds. You can configure it by using `WithShutdownTimeout`

```go
func main() {
	...
	...
	
	if err := graceful.ListenAndServe(
		":8080",
		handler,
		graceful.WithShutdownTimeout(10*time.Second),
		// add customer function before shutting down
		graceful.WithShutdownFunc(func() {
			log.Info("Http Server is shutting down...")
        }
	); err != nil {
		log.Fatal(err)
	}
}
```
