package wc

import (
	"context"
	"github.com/user/andon-webapp-in-go/src/view"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type httpHandler struct {
	urlPattern *regexp.Regexp
}

//NewViewHandler returns a Handler that returns for HTML files related to WorkConters
func NewViewHandler() http.Handler {
	return &httpHandler{
		regexp.MustCompile(`^/wc/(\d+)$`),
	}
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//Add request context
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	matches := h.urlPattern.FindStringSubmatch(r.URL.Path)
	if len(matches) == 0 {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Print("failed to convert %q to integer: %v", matches[1], err)
		http.NotFound(w, r)
		return
	}

	//create the done channel

	done := make(chan struct{})

	go h.getView(ctx, id, r, w, done)

	select {
	case <-ctx.Done():
		http.Error(w, ctx.Err().Error(), http.StatusNotFound)
		<-done
	case <-done:
		cancel()
	}
}

func (h *httpHandler) getView(ctx context.Context,
	id int, r *http.Request,
	w http.ResponseWriter,
	done chan<- struct{}) {

	defer func() {
		done <- struct{}{}
	}()

	//Add this to test context deadline
	//time.Sleep(5*time.Second)

	t, err := view.Get("workcenter")
	if err != nil {
		log.Println(err)
		http.Error(w, "View template not found", http.StatusNotFound)
		return
	}
	wc, err := GetWorkCenter(ctx, id)
	if err != nil {
		log.Panicln(err)
		http.NotFound(w, r)
		return
	}

	select {
	case <- ctx.Done():
		return
	default:
	}
	w.Header().Add("Content-Type", "text/html")
	err = t.Execute(w, struct {
		Workcenter
		view.PipelineBase
	}{
		Workcenter:   wc,
		PipelineBase: view.PipelineBase{Title: wc.Name},
	})
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to generate view", http.StatusInternalServerError)
	}
}
