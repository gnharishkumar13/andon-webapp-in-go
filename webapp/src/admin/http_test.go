package admin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/andon-webapp-in-go/src/hash"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestLogonViewHandler(t *testing.T) {

	t.Run("logon", func(t *testing.T) {
		t.Run("redirects to /admin when valid credentials are provided", func(t *testing.T) {
			username, password := "user", "password"
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			database = db

			rows := sqlmock.NewRows([]string{"id", "username"})
			rows.AddRow(1, username)

			mock.ExpectQuery(`SELECT 
			id,
			username 
				FROM users
				WHERE username=\$1 AND password=\$2`).
				WithArgs(username, hash.Hash(username, password)).
				WillReturnRows(rows)

			//Expect a INSERT query - per the flow in saveLogonToken
			mock.ExpectExec(`^INSERT(.+)`).
				WillReturnResult(sqlmock.NewResult(1, 1))

			req := httptest.NewRequest(http.MethodPost, "/admin/logon",
				nil)

			req.Form = map[string][]string{}
			req.Form.Add("username", username)
			req.Form.Add("password", password)

			res := httptest.NewRecorder()

			lvh := logonViewHandler{}
			lvh.logon(context.Background(), res, req, make(chan struct{}, 1))

			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Database did not get expected queries %v", err)
			}

			if res.Code != http.StatusSeeOther {
				t.Errorf("Unexpected response code")
			}

			if len(res.Result().Cookies()) != 1 {
				t.Errorf("Expected Cookie not present")
			}

			if res.Header().Get("Location") != "/admin" {
				t.Errorf("Expected redirect to /admin did not happen")
			}
		})
	})
}
