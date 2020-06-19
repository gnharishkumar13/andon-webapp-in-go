package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/user/andon-webapp-in-go/src/db"
	"github.com/user/andon-webapp-in-go/src/hash"
)

type auth struct {
	wrapped    http.Handler
	exceptions []string
	roles      []string
}

// NewAuth creates a middleware that ensures that the current user
// has the proper role to access the wrapped http.Handler.
//
// Authorization is verified against the list of roles that are provided.
// If nil is provided for the roles, then any logged in user is granted access.
//
// If a specific path should not be checked, then it should be added to the exceptions slice.
func NewAuth(roles, exceptions []string, wrapped http.Handler) http.Handler {
	if roles == nil {
		roles = make([]string, 0)
	}
	if exceptions == nil {
		exceptions = make([]string, 0)
	}
	if wrapped == nil {
		wrapped = http.DefaultServeMux
	}

	return &auth{
		wrapped:    wrapped,
		exceptions: exceptions,
		roles:      roles,
	}
}

func (a auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	//Add timeout and back to existing request context
	r = r.WithContext(ctx)
	for _, e := range a.exceptions {
		if strings.HasPrefix(r.URL.Path, e) {
			a.wrapped.ServeHTTP(w, r)
			return
		}
	}

	//no roles
	if len(a.roles) == 0 {
		a.wrapped.ServeHTTP(w, r)
		return
	}

	authCh := make(chan bool)
	errCh := make(chan error)

	go a.authorized(ctx, r, authCh, errCh)

	select {
	case authorized := <-authCh:
		if authorized {
			a.wrapped.ServeHTTP(w, r)
			return
		}
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	case <-errCh:
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	case <-ctx.Done():
		http.Error(w, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
		select {
		case <-authCh:
		case <-errCh:
		}

	}
}

func (a auth) authorized(ctx context.Context, r *http.Request, authCh chan bool, errCh chan error) {

	database, err := db.GetDB()
	if err != nil {
		log.Printf("unable to get instance of database %v", err)
		errCh <- err
		return
	}
	cookie, err := r.Cookie("token")
	if err != nil {
		log.Printf("unable to get valid token from request %v", err)
		errCh <- err
		return
	}

	token := cookie.Value
	result, err := database.QueryContext(ctx,
		`SELECT role 
		FROM roles
		JOIN users_roles 
			ON users_roles.role_id = roles.id
		JOIN logon_tokens 
			ON logon_tokens.user_id = users_roles.user_id
		WHERE logon_tokens.token = $1`, hash.Hash(token))
	if err != nil {
		log.Printf("failed to retrieve roles for provided token %q:%v", token, err)
		errCh <- err
		return
	}

	defer result.Close()
	for result.Next() {
		var role string
		err := result.Scan(&role)
		if err != nil {
			log.Printf("failed to retrieve role from database query %v", err)
			errCh <- err
			return
		}
		for _, authorizedRole := range a.roles {
			if authorizedRole == role {
				authCh <- true
				return
			}
		}
	}
	authCh <- false
}
