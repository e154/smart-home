package build

import "embed"

//go:embed public/*
//go:embed public/index.html
//go:embed public/css/*
//go:embed public/js/*
//go:embed public/fonts/*
//go:embed public/tinymce/*
//go:embed public/tinymce/langs/*
//go:embed public/tinymce/skins/*
//go:embed public/tinymce/skins/fonts/*
//go:embed public/img/*
//go:embed public/img/icons/*
//go:embed public/img/icons/*
var F embed.FS
