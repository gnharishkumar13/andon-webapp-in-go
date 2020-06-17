package admin

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/user/andon-webapp-in-go/src/db"
)

var database *sql.DB

// SetDB sets the database instance that will be used for all SQL commands.
// The package assumes that a properly configured sql.DB has been provided via this function
// before any of its other functions are executed.
func SetDB(db *sql.DB) {
	database = db
}

func findOne(ctx context.Context, username, pwdHash string) (user, error) {
	result := database.QueryRowContext(ctx,
		`SELECT 
			id,
			username 
		FROM users
		WHERE username=$1 AND password=$2`, username, pwdHash)
	var u user
	err := result.Scan(&u.ID, &u.Username)
	if err != nil {
		return user{}, fmt.Errorf("failed to retrieve user from database with username %q: %v", username, err)
	}
	return u, nil
}

func saveLogonToken(ctx context.Context, token string, userID int) error {
	expiration := db.ToTimestamp(time.Now().Add(24 * time.Hour))

	// TODO: create a job that will remove expired tokens
	_, err := database.ExecContext(ctx,
		`INSERT INTO logon_tokens (token, user_id, expiration)
		VALUES ($1, $2, $3)`, token, userID, expiration)
	if err != nil {
		msg := fmt.Sprintf("failed to save token: %v", err)
		log.Println(msg)
		return fmt.Errorf(msg)
	}
	return nil
}

func isValidToken(ctx context.Context, token string) (bool, error) {
	row := database.QueryRowContext(ctx,
		`SELECT count(token)
		FROM logon_tokens
		WHERE token=$1`, token)
	var count int
	err := row.Scan(&count)
	if err != nil {
		msg := fmt.Sprintf("failed to retrieve current tokens: %v", err)
		log.Println(msg)
		return false, err
	}
	return count > 0, nil
}
