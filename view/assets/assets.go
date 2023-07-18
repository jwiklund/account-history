package assets

import "embed"

//go:embed edit.html
//go:embed head.html
//go:embed index.html
//go:embed nav.html
//go:embed scripts.html
var EmbedFs embed.FS
