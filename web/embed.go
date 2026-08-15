// Package web embeds the static frontend files (index.html, style.css,
// app.js) into the Go binary via go:embed so the calculator can serve
// a web UI without external file dependencies.
package web

import "embed"

// FS is the embedded filesystem containing the web UI static assets.
//
//go:embed index.html style.css app.js
var FS embed.FS
