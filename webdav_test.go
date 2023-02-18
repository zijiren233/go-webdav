package gowebdav

import (
	"testing"
)

func TestWebdav(t *testing.T) {
	ser := NewWebdav()
	ser.DefaultClient("", ".")
	t.Log(ser.Run(":8080"))
}
