package gowebdav

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"
)

const (
	style     = `<style>table {border-collapse: separate;border-spacing: 1.5em 0.25em;}h1 {padding-left: 0.3em;}a {text-decoration: none;color: blue;}.left {text-align: left;}.mono {font-family: monospace;}.mw20 {min-width: 20em;}</style>`
	meta      = `<meta name="referrer" content="no-referrer" />`
	listIndex = `<tr><th class="left mw20">Name</th><th class="left">Last modified</th><th>Size</th></tr><tr><th colspan="3"><hr></th></tr>`
	homeDIr   = "<tr><td><a href=\"%s\">Home Dir</a></td><td>&nbsp;</td><td class=\"mono\" align=\"right\">[DIR]</td></tr>"
	perDir    = `<td><a href="..">Pre Dir</a></td><td>&nbsp;</td><td class="mono" align="right">[DIR]</td></tr>`
	fileuri   = "<tr><td><a href=\"%s\" >%s</a></td><td class=\"mono\">%s</td><td class=\"mono\" align=\"right\">%s</td></tr>"
)

func (client *Cli) handleDirList(fs webdav.FileSystem, ctx *gin.Context) bool {
	filePath := ctx.Params.ByName("webdav")
	f, err := fs.OpenFile(ctx, filePath, os.O_RDONLY, 0)
	if err != nil {
		return false
	}
	defer f.Close()
	if fi, err := f.Stat(); err != nil || fi == nil || !fi.IsDir() {
		return false
	}
	dirs, err := f.Readdir(-1)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return true
	}
	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	client.generateWeb(dirs, filePath, ctx.Writer)
	return true
}

func (client *Cli) generateWeb(dirs []fs.FileInfo, path string, writer io.Writer) {
	fmt.Fprintf(writer, "<html><head>%s<title>Index of %s</title>%s</head>", meta, path, style)
	if client.pathPrefix == "" {
		fmt.Fprintf(writer, "<body><h1>Index of /<a href=\"/\">Home</a>%s</h1><table>%s%s", path2index(path), listIndex, fmt.Sprintf(homeDIr, "/"))
	} else {
		fmt.Fprintf(writer, "<body><h1>Index of /<a href=\"%s\">Home</a>%s</h1><table>%s%s", client.pathPrefix, path2index(path), listIndex, fmt.Sprintf(homeDIr, client.pathPrefix))
	}
	if path != "/" {
		fmt.Fprint(writer, perDir)
	}
	for _, d := range dirs {
		name := d.Name()
		if d.IsDir() {
			name += "/"
			fmt.Fprintf(writer, fileuri, name, name, d.ModTime().Format("2006/1/2 15:04:05"), "[DIR]")
		} else {
			fmt.Fprintf(writer, fileuri, name, name, d.ModTime().Format("2006/1/2 15:04:05"), getsize(d.Size()))
		}
	}
	fmt.Fprint(writer, `</table></body></html>`)
}
