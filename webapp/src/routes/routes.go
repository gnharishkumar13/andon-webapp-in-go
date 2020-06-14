package routes

import (
	"github.com/user/andon-webapp-in-go/src/admin"
	"github.com/user/andon-webapp-in-go/src/wc"
	"net/http"
)

//Register all handlers
func Register(){
	http.Handle("/", http.NotFoundHandler())
	http.Handle("/wc/", wc.NewViewHandler())
	http.Handle("/admin", admin.NewViewHandler())
	http.Handle("/api/wc/", wc.NewAPIHandler())
}
