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

	cookie, err := r.Cookie("token")
	if err != nil {
		log.Printf("unable to retrieve token: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	valid := validateToken(ctx, cookie.Value)
	if !valid {
		log.Printf("invalid token received %q from address %q", cookie.Value, r.RemoteAddr)
		http.Redirect(w, r, "/admin/logon", http.StatusSeeOther)
		return
	}

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

type logonForm struct {
	Username string
}

//Logon page for admins
type logonViewHandler struct{}

func NewLogonViewHandler() http.Handler {
	return &logonViewHandler{}
}

//Initialize a handler, which would return the ServeHTTP
func (h logonViewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.get(w, r, logonForm{})
	case http.MethodPost:
		h.post(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h logonViewHandler) get(w http.ResponseWriter, r *http.Request, form logonForm) {
	t, err := view.Get("admin_logon")
	if err != nil {
		log.Println("unable to find view template for admin logon")
		http.Error(w, "view template not found", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "text/html")
	err = t.Execute(w, struct {
		view.PipelineBase
		logonForm
	}{
		PipelineBase: view.PipelineBase{Title: "Administrator Logon"},
		logonForm:    form,
	})
	if err != nil {
		log.Print(err)
		http.Error(w, "failed to generate view", http.StatusInternalServerError)
	}
}

func (h logonViewHandler) post(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	done := make(chan struct{})

	go h.logon(ctx, w, r, done)

	select {
	case <-ctx.Done():
		http.Error(w, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
		<-done
	case <-done:
		cancel()
	}
}

//receive channel chan<-struct{}
func (h logonViewHandler) logon(ctx context.Context, w http.ResponseWriter,
	r *http.Request, done chan<- struct{}) {
	defer func() {
		done <- struct{}{}
	}()
	username := r.FormValue("username")
	password := r.FormValue("password")
	verified, user, err := verifyCredentials(ctx, username, password)
	if err != nil {
		log.Printf("failed to logon user %q: %v", username, err)
		h.get(w, r, logonForm{Username: username})
		return
	}
	if !verified {
		h.get(w, r, logonForm{Username: username})
		return
	}
	logonToken, err := createLogonToken(ctx, user.ID)
	if err != nil {
		log.Printf("failed to generate logon token for user %q: %v", username, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		Name:    "token",
		Value:   logonToken,
		Expires: time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
