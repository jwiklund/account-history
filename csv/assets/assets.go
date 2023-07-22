package assets

import "embed"

//go:embed *.yaml
var EmbedFS embed.FS
