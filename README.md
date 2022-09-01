# go-webdav

```go
func main() {
	ser := gowebdav.NewWebdav()
	ser.NewClient("", ".")
	ser.Run(":8080")
}
```

```go
func main() {
	ser := gowebdav.NewWebdav()
	ser.NewClient("/prefix1", ".")
	c := ser.NewClient("/prefix2", "..")
	c.AddUser("admin", "admin", gowebdav.O_RDWR)
	ser.Run(":8080")
}
```