package public

import (
	"embed"
	"io/fs"
)

//go:embed all:public
var content embed.FS
var F, _ = fs.Sub(content, "public")
