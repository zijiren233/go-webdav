package gowebdav

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"
)

func (client *client) handleDirList(fs webdav.FileSystem, ctx *gin.Context) bool {
	filePath := ctx.Params.ByName("webdav")
	f, err := fs.OpenFile(ctx, filePath, os.O_RDONLY, 0)
	if err != nil {
		return false
	}
	defer f.Close()
	if fi, _ := f.Stat(); fi != nil && !fi.IsDir() {
		return false
	}
	dirs, err := f.Readdir(-1)
	if err != nil {
		log.Print(ctx, "Error reading directory", http.StatusInternalServerError)
		return false
	}
	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.Writer.Write(client.generateWeb(dirs, filePath).Bytes())
	return true
}

func (client *client) generateWeb(dirs []fs.FileInfo, path string) *bytes.Buffer {
	data := bytes.NewBuffer(nil)
	fmt.Fprintln(data, `<html><head>
<meta name=\"referrer\" content=\"no-referrer\" />`)
	fmt.Fprintf(data, "<title>Index of %s</title>\n", path)
	fmt.Fprintln(data, `<style>
table {
border-collapse: separate;
border-spacing: 1.5em 0.25em;
}
h1 {
padding-left: 0.3em;
}
a {
text-decoration: none;
color: blue;
}
.left {
text-align: left;
}
.mono {
font-family: monospace;
}
.mw20 {
min-width: 20em;
}
</style>
</head>
<body>`)
	fmt.Fprintf(data, "<h1>Index of %s</h1><table>\n", path)
	fmt.Fprintln(data, `<tr>
<th class="left mw20">Name</th>
<th class="left">Last modified</th>
<th>Size</th>
</tr>
<tr><th colspan="3"><hr></th></tr>
<tr>`)
	fmt.Fprintf(data, "<td><a href=\"%s\">Home Dir</a></td>\n", client.pathPrefix)
	fmt.Fprintln(data, `<td>&nbsp;</td>
<td class="mono" align="right">[DIR]    </td>
</tr>`)
	if path != "/" {
		fmt.Fprintln(data, `<td><a href="..">Pre Dir</a></td>
<td>&nbsp;</td>
<td class="mono" align="right">[DIR]    </td>
</tr>`)
	}
	for _, d := range dirs {
		name := d.Name()
		if d.IsDir() {
			name += "/"
			fmt.Fprintf(data, "<tr><td><a href=\"%s\" >%s</a></td><td class=\"mono\">%s</td><td class=\"mono\" align=\"right\">[DIR]    </td></tr>", name, name, d.ModTime().Format("2006/1/2 15:04:05"))
		} else {
			fmt.Fprintf(data, "<tr><td><a href=\"%s\" >%s</a></td><td class=\"mono\">%s</td><td class=\"mono\" align=\"right\">%s    </td></tr>", name, name, d.ModTime().Format("2006/1/2 15:04:05"), getsize(d.Size()))
		}
	}
	return data
}
