package main

import (
	"log"
	"net/http"

	"github.com/user/andon-webapp-in-go/src/middleware"

	"github.com/user/andon-webapp-in-go/src/admin"
	"github.com/user/andon-webapp-in-go/src/db"
	"github.com/user/andon-webapp-in-go/src/routes"
	"github.com/user/andon-webapp-in-go/src/view"
	"github.com/user/andon-webapp-in-go/src/wc"
)

func main() {
	database, err := db.GetDB()
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	wc.SetDB(database)
	admin.SetDB(database)

	view.RegisterStaticHandlers()
	routes.Register()

	gzipMiddleware := middleware.NewGzip([]string{"/ws"}, nil)
	log.Fatal(http.ListenAndServe(":3000", gzipMiddleware))

}
