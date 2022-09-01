# go-webdav

```go
package main

import (
	gowebdav "github.com/zijiren233/go-webdav"
)

func main() {
	ser := gowebdav.NewWebdav()
	ser.NewClient("", ".")
	ser.Run(":8080")
}
```

<br>

```go
package main

import (
	gowebdav "github.com/zijiren233/go-webdav"
)

func main() {
	ser := gowebdav.NewWebdav()
	ser.NewClient("/prefix1", ".")
	c := ser.NewClient("/prefix2", "../")
	c.AddUser("admin", "admin", gowebdav.O_RDWR)
	ser.Run(":8080")
}
```