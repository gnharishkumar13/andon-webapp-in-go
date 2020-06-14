package admin

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/user/andon-webapp-in-go/src/view"
	"github.com/user/andon-webapp-in-go/src/wc"
)

type viewHandler struct{}

// NewViewHandler returns a handler that manages requests to the admin page.
func NewViewHandler() http.Handler {
	return &viewHandler{}
}

func (h viewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.get(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h viewHandler) get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	done := make(chan struct{})
	go h.getView(ctx, w, r, done)

	select {
	case <-ctx.Done():
		http.Error(w, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
		<-done
	case <-done:
		cancel()
	}
}

func (h viewHandler) getView(ctx context.Context, w http.ResponseWriter,
	r *http.Request, done chan<- struct{}) {

	defer func() {
		done <- struct{}{}
	}()
	workcenters, err := wc.GetAllWorkcenters(ctx)
	if err != nil {
		log.Printf("unable to retrieve workcenters: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	t, err := view.Get("admin")
	if err != nil {
		log.Printf("unable to find view template for admin: %v", err)
		http.Error(w, "view template not found", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "text/html")
	err = t.Execute(w, struct {
		view.PipelineBase
		Workcenters []wc.Workcenter
	}{
		PipelineBase: view.PipelineBase{Title: "Administrator Dashboard"},
		Workcenters:  workcenters,
	})
	if err != nil {
		log.Print(err)
		http.Error(w, "failed to generate view", http.StatusInternalServerError)
	}
}
