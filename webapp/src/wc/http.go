package wc

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/user/andon-webapp-in-go/src/view"
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
	wc, err := GetWorkcenter(ctx, id)
	if err != nil {
		log.Panicln(err)
		http.NotFound(w, r)
		return
	}

	select {
	case <-ctx.Done():
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

type apiHandler struct {
	escalateRoutePattern *regexp.Regexp
}

// NewAPIHandler returns an http.Handler that is setup to respond to
// asynchronous calls that relate to workcenters.
func NewAPIHandler() http.Handler {
	return &apiHandler{
		escalateRoutePattern: regexp.MustCompile(`^\/api\/wc\/(\d+)/escalate$`),
	}
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	matches := h.escalateRoutePattern.FindStringSubmatch(r.URL.Path)
	if len(matches) == 0 {
		http.NotFound(w, r)
		return
	}
	wcID, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Printf("failed to convert workcenter ID %q to number: %v", matches[1], err)
		http.NotFound(w, r)
		return
	}
	h.escalate(ctx, w, r, wcID)
}

func (h apiHandler) escalate(ctx context.Context, w http.ResponseWriter,
	r *http.Request, id int) {

	doneCh := make(chan struct{})
	errCh := make(chan error)

	go func(doneCh chan<- struct{}, errChan chan<- error) {
		err := escalate(ctx, id)
		if err != nil {
			errCh <- err
			return
		}
		UpdateNow <- struct{}{}
		doneCh <- struct{}{}
	}(doneCh, errCh)

	select {
	case <-ctx.Done():
		http.Error(w, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
		select {
		case <-doneCh:
		case <-errCh:
		}
	case err := <-errCh:
		log.Printf("failed to escalate workcenter %q: %v", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	case <-doneCh:
		// function succeeded, nothing else to do!
	}

}

type wsHandler struct {
	upgrader    websocket.Upgrader
	idPattern   *regexp.Regexp
	connections map[int][]*websocket.Conn
	mConn       *sync.RWMutex
}

var once = sync.Once{}

// UpdateNow is a channel that is used to trigger an update message to be
// sent to all attached clients. This is normally used when the status of one
// or more workcenters has changed in a way that the client's should be aware of.
var UpdateNow = make(chan struct{})

// NewWebsocketHandler returns a handler that manages websockets that are communicating
// information about a single workcenter between a client and the server.
func NewWebsocketHandler() http.Handler {
	h := wsHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		idPattern:   regexp.MustCompile(`^\/ws\/wc\/(\d+)$`),
		connections: make(map[int][]*websocket.Conn),
		mConn:       &sync.RWMutex{},
	}

	once.Do(func() {
		go h.monitorWorkcenters()
	})
	return &h
}

func (h *wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	matches := h.idPattern.FindStringSubmatch(r.URL.Path)
	if len(matches) == 0 {
		http.NotFound(w, r)
		return
	}
	id, err := strconv.Atoi(matches[1])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, ok := h.connections[id]; !ok {
		h.connections[id] = []*websocket.Conn{conn}
	} else {
		h.connections[id] = append(h.connections[id], conn)
	}
	go h.listenToWebsocket(id, conn)
}

type statusChangePayload struct {
	WorkcenterID int `json:"workcenter"`
	Status       int `json:"status"`
}

type statusChangeMessage struct {
	Type    string              `json:"type"`
	Payload statusChangePayload `json:"payload"`
}

type currentStatusMessage struct {
	Status                     int    `json:"status"`
	StatusDescription          string `json:"statusDescription"`
	EscalationLevelDescription string `json:"escalationLevelDescription"`
	TimeAtStatus               string `json:"timeAtStatus"`
	TimeTillEscalation         string `json:"timeTillEscalation"`
}

func (h wsHandler) listenToWebsocket(id int, conn *websocket.Conn) {
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("unable to read message from client %v: %v", conn.RemoteAddr(), err)
			h.closeAndRemoveConnection(id, conn)
			break
		}
		switch {
		case bytes.Contains(data, []byte(`"type":"statusChange"`)):
			var msg statusChangeMessage
			err = json.Unmarshal(data, &msg)
			if err != nil {
				log.Printf("unable to parse message from client %v: %v", conn.RemoteAddr(), err)
				err = conn.WriteMessage(websocket.TextMessage, []byte("failed to read message"))
				if err != nil {
					log.Printf("unable to write message to client %v: %v", conn.RemoteAddr(), err)
				}
				continue
			}
			err = SetWorkcenterStatus(context.Background(),
				msg.Payload.WorkcenterID, msg.Payload.Status)
			if err != nil {
				log.Println(err)
			}
			UpdateNow <- struct{}{}
		default:
			log.Printf("Unknown message type received: %v", string(data))
		}
	}
}

func (h wsHandler) closeAndRemoveConnection(wcID int, conn *websocket.Conn) {
	h.mConn.Lock()
	for i, c := range h.connections[wcID] {
		if c == conn {
			h.connections[wcID] = append(h.connections[wcID][:i], h.connections[wcID][i+1:]...)
			break
		}
	}

	h.mConn.Unlock()
	conn.Close()
}

func (h wsHandler) monitorWorkcenters() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
		case <-UpdateNow:
		}
		h.mConn.RLock()
		for id, conns := range h.connections {
			// only check if something is listening
			if len(conns) == 0 {
				continue
			}
			wc, err := GetWorkcenter(context.Background(), id)
			if err != nil {
				log.Printf("unable to retrieve information for workcenter %q: %v", id, err)
				continue
			}
			msg := currentStatusMessage{
				Status:                     wc.Status,
				StatusDescription:          wc.StatusDescription(),
				EscalationLevelDescription: wc.EscalationLevelDescription(),
				TimeAtStatus:               view.DurationToHHMMSS(wc.TimeAtStatus()),
				TimeTillEscalation:         view.DurationToHHMMSS(wc.TimeTillEscalation()),
			}
			for _, c := range conns {
				err := c.WriteJSON(msg)
				if err != nil {
					log.Printf("unable to send workcenter data to websocket at address %v: %v", c.RemoteAddr(), err)
				}
			}
		}
		h.mConn.RUnlock()
	}
}
