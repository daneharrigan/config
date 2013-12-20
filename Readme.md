# config

A package in Go for parsing INI files

### Usage

Your ini file located in your current directory.

```ini
; ./.config
[section]
  key = value
```

Your Go file reading in the `.config` file.

```go
package main

import (
	"github.com/daneharrigan/config"
	"log"
)

func main() {
	c, err := config.New("./.config")
	if err != nil {
		log.Fatalf("fn=New error=%q", err)
	}

	v, err := c.Get("section", "key")
	if err != nil {
		log.Fatalf("fn=Get error=%q", err)
	}

	log.Printf("value=%q", v)
}
```
