package gowebdav

import (
	"fmt"
	"os"
)

func getsize(size int64) string {
	tmp := float64(size)
	if tmp < 1024 {
		return fmt.Sprintf("%d B", size)
	} else if tmp = tmp / 1024; tmp < 1024 {
		return fmt.Sprintf("%.2f KB", tmp)
	} else if tmp = tmp / 1024; tmp < 1024 {
		return fmt.Sprintf("%.2f MB", tmp)
	} else if tmp = tmp / 1024; tmp < 1024 {
		return fmt.Sprintf("%.2f GB", tmp)
	} else {
		return fmt.Sprintf("%.2f TB", tmp/1024)
	}
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			return ":" + port
		}
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}