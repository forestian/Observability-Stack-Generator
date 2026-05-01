package templates

import "embed"

// FS contains the generator templates used by internal/generator.
//
//go:embed *.tmpl
var FS embed.FS
