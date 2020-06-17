package admin

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/user/andon-webapp-in-go/src/hash"
)

type user struct {
	ID       int
	Username string
}

func verifyCredentials(ctx context.Context, username, password string) (bool, user, error) {
	u, err := findOne(ctx, username, hash.Hash(username, password))
	if err != nil {
		return false, user{}, err
	}
	return true, u, nil
}

func createLogonToken(ctx context.Context, userID int) (string, error) {
	ns := strconv.Itoa(time.Now().Nanosecond())
	token := base64.StdEncoding.EncodeToString([]byte(ns))
	tokenHash := hash.Hash(token)
	err := saveLogonToken(ctx, tokenHash, userID)
	if err != nil {
		msg := fmt.Sprintf("failed to save token: %v", err)
		log.Printf(msg)
		return "", fmt.Errorf(msg)
	}
	return token, nil
}

func validateToken(ctx context.Context, token string) bool {
	valid, err := isValidToken(ctx, hash.Hash(token))
	if err != nil {
		log.Printf("failed to validate logon token: %v", err)
		return false
	}
	return valid
}
