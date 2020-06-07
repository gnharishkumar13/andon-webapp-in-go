package main

import (
	"log"
	"net/http"
	"time"

	"github.com/user/andon-webapp-in-go/src/view"
	"github.com/user/andon-webapp-in-go/src/wc"
)

func main() {
	view.RegisterStaticHandlers()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := view.Get("workcenter")
		if err != nil {
			log.Fatal(err)
		}
		workcenter := wc.Workcenter{
			ID:              1,
			Name:            "Assembly Line 1",
			CurrentProduct:  "Widgets",
			Status:          1,
			EscalationLevel: 2,
			StatusSetAt:     time.Now(),
		}
		err = t.Execute(w, struct {
			wc.Workcenter
			view.PipelineBase
		}{
			Workcenter:   workcenter,
			PipelineBase: view.PipelineBase{Title: workcenter.Name},
		})
		if err != nil {
			log.Fatal(err)
		}
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
