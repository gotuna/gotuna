package static

import "embed"

//go:embed *
var EmbededStatic embed.FS
