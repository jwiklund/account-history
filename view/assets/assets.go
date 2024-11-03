package assets

import "embed"

//go:embed *.html *.ico
var EmbedFs embed.FS
