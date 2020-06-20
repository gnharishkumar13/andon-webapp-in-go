package admin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/user/andon-webapp-in-go/src/hash"
)

func TestLogonViewHandler(t *testing.T) {
	t.Run("logon", func(t *testing.T) {
		t.Run("redirects to /admin when valid credentials are provided", func(t *testing.T) {
			ctx := context.Background()
			doneCh := make(chan struct{}, 1)
			username, password := "user", "password"
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock database instance: %v", err)
			}
			defer db.Close()
			database = db

			rows := sqlmock.NewRows([]string{"id", "username"}).
				AddRow(1, username)
			mock.ExpectQuery(`SELECT 
				id,
				username 
			FROM users
			WHERE username=\$1 AND password=\$2`).
				WithArgs(username, hash.Hash(username, password)).
				WillReturnRows(rows)

			mock.ExpectExec("^INSERT (.+)").
				WillReturnResult(sqlmock.NewResult(1, 1))

			req := httptest.NewRequest(http.MethodPost, "/admin/logon", nil)
			req.Form = map[string][]string{}
			req.Form.Add("username", username)
			req.Form.Add("password", password)

			res := httptest.NewRecorder()

			lvh := logonViewHandler{}
			lvh.logon(ctx, res, req, doneCh)

			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Database did not receive expected queries: %v\n", err)
			}
			if res.Code != http.StatusSeeOther {
				t.Error("Unexpected response code received\n")
			}
			if len(res.Result().Cookies()) != 1 {
				t.Error("Expected cookie not present in response\n")
			}
			if res.Header().Get("Location") != "/admin" {
				t.Error("Expected redirect to /admin did not occur\n")
			}
		})
	})
}
