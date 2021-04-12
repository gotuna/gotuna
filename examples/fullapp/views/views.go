package views

import "embed"

// EmbededViews holds HTML templates.
//
//go:embed *
var EmbededViews embed.FS
