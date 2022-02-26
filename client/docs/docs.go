package docs

import "net/http"

// Docs represents the embedded doc file
//go:embed static
var Docs http.FileSystem
