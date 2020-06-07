package wc

import (
	"github.com/user/andon-webapp-in-go/src/view"
	"log"
	"net/http"
	"regexp"
	"strconv"
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

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){

	matches := h.urlPattern.FindStringSubmatch(r.URL.Path)
	if len(matches) == 0 {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Print("failed to convert %q to integer: %v", matches[1], err )
		http.NotFound(w, r)
		return
	}


	t, err := view.Get("workcenter")
	if err != nil {
		log.Println(err)
		http.Error(w, "View template not found", http.StatusNotFound)
		return
	}
	wc, err := GetWorkCenter(id)
	if err != nil {
		log.Panicln(err)
		http.NotFound(w, r)
		return
	}

	w.Header().Add("Content-Type","text/html")
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