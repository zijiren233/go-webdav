# go-webdav

```go
func main() {
	ser := gowebdav.NewSingleWebdav(".")
	ser.SetAuth("admin", "admin")
	ser.Run(":8080")
}
```

```go
func main() {
	ser := gowebdav.NewWebdav()
	client1 := ser.NewClient("/prefix", ".")
	client1.SetAuth("admin", "admin")
	ser.NewClient("/prefix2", "../")
	ser.Run()
}
```