package main

import (
	"github.com/user/andon-webapp-in-go/src/view"
	"log"
	"net/http"
)

func main() {
	view.RegisterStaticHandlers()
	log.Fatal(http.ListenAndServe(":3000", nil))
}
