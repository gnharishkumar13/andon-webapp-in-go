package admin

import (
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
		})
	})
}
