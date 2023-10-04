package build

import "embed"

//go:embed public/*
//go:embed public/assets/*
//go:embed public/index.html
//go:embed public/favicon.ico
//go:embed public/logo.png
var F embed.FS
