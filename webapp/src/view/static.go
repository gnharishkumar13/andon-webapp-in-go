package view

import (
	"github.com/user/andon-webapp-in-go/src/config"
	"net/http"
)

// RegisterStaticHandlers registers HTTP handlers that will serve static
// content such as CSS and JavaScript files
func RegisterStaticHandlers() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.Get().StaticRoot))))
}
