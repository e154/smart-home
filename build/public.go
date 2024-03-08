package public

import (
	"embed"
	"io/fs"
)

//go:embed public/*
//go:embed public/assets/*
//go:embed public/index.html
//go:embed public/favicon.ico
var content embed.FS
var F, _ = fs.Sub(content, "public")
