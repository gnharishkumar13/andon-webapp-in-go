package main

import (
	"github.com/user/andon-webapp-in-go/src/db"
	"github.com/user/andon-webapp-in-go/src/routes"
	"github.com/user/andon-webapp-in-go/src/view"
	"github.com/user/andon-webapp-in-go/src/wc"
	"net/http"

	"log"
)

func main() {
	database, err := db.GetDB()
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	wc.SetDB(database)

	view.RegisterStaticHandlers()
	routes.Register()

	log.Fatal(http.ListenAndServe(":3000", nil))
}
