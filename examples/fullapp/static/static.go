package static

import "embed"

// EmbededStatic holds static web server content.
//
//go:embed *
var EmbededStatic embed.FS
