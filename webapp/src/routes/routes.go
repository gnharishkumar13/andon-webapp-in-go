package routes

import (
	"net/http"

	"github.com/user/andon-webapp-in-go/src/admin"
	"github.com/user/andon-webapp-in-go/src/wc"
)

//Register all handlers
func Register() {
	http.Handle("/", http.NotFoundHandler())
	http.Handle("/wc/", wc.NewViewHandler())
	http.Handle("/admin", admin.NewViewHandler())
	http.Handle("/admin/logon", admin.NewLogonViewHandler())
	http.Handle("/api/wc/", wc.NewAPIHandler())
}
