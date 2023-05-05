![cover](assets/cover.png)

# go-herder
Tiny process management tool

## Quick start

```golang
package main

import (
	"github.com/kirsupi/go-herder"
	"log"
	"runtime"
)

func main() {
	h := herder.New(herder.Config{
		MaxWorkersCount:     runtime.NumCPU(),
		Logger:              log.Default(),
		DefaultMaxStdoutLen: 32768,
		DefaultMaxStderrLen: 32768,
	})
	h.Run()
}
```

## Methods
