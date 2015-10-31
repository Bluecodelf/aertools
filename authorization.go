package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/Bluecodelf/epigo"
	"net/http"
	"strings"
	"time"
)

// Token size in bytes.
// It will be stored in hex in the database, thus needing a string with a
// length twice as big as this value.
const TokenSize = 32

type AuthInfo struct {
	id    int
	login string
	token string
	level int // Still unused. Will be used for permission management.
}

func GetAuthorizationLevel(request *http.Request) (AuthInfo, error) {
	auth := request.Header.Get("Authorization")
	if auth == "" {
		return AuthInfo{-1, "", "", -1}, errors.New("No authorization header")
	}

	if strings.HasPrefix(auth, "Basic") && len(auth) > 6 {
		auth = auth[6:]
		return GetAuthorizationLevelBasic(auth)
	} else if strings.HasPrefix(auth, "Token") && len(auth) > 6 {
		auth = auth[6:]
		return GetAuthorizationLevelToken(auth)
	}
	return AuthInfo{-1, "", "", -1}, errors.New("Invalid authorization method")
}

func GetAuthorizationLevelBasic(basic string) (auth AuthInfo, err error) {
	// Decode the RFC 1521 base64 credentials and split them to user/pass.
	data, _ := base64.StdEncoding.DecodeString(basic)
	basic = string(data)
	separatorIndex := strings.Index(basic, ":")
	username := basic[:separatorIndex]
	password := basic[separatorIndex+1:]
	auth.id = -1
	auth.level = -1

	// Check username and password on the Epitech intranet
	// REVIEW: We create a new Epitech Client at each connection. This is
	// fucking ugly. Need to add a CheckAuthentication function in epigo.
	epitech := epigo.Client{Host: "https://intra.epitech.eu/"}
	err = epitech.Authenticate(username, password)
	if err != nil {
		return
	}

	// Authorization is successful, we can get authlevel for user and
	// retrieve/create the authentication token.
	var rows *sql.Rows
	rows, err = db.Query("SELECT id, level FROM accounts WHERE login = ?",
		username)
	defer rows.Close()
	if err != nil || !rows.Next() {
		if err == nil {
			err = errors.New("User is not an authorized AER.")
		}
		return
	}
	rows.Scan(&auth.id, &auth.level)
	auth.token, _ = GetTokenForUser(auth.id)
	return
}

func GetAuthorizationLevelToken(token string) (auth AuthInfo, err error) {
	var rows *sql.Rows
	auth.id = -1
	auth.level = -1

	// Get UserID and expiration date from token
	rows, err = db.Query("SELECT user_id, expiration FROM auth_tokens "+
		"WHERE token = ?", token)
	if err != nil || !rows.Next() {
		// TODO: need to do some research on rows.Close here. Potential crash.
		if err == nil {
			err = errors.New("Invalid token")
		}
		return
	}
	var expiration time.Time
	rows.Scan(&auth.id, &expiration)
	rows.Close()

	if expiration.Before(time.Now()) {
		err = errors.New("Expired token")
		return
	}

	// Get authentication level from UserID
	rows, err = db.Query("SELECT level FROM accounts WHERE id = ?", auth.id)
	defer rows.Close()
	if err != nil || !rows.Next() {
		if err == nil {
			err = errors.New("The universe has collapsed") // wtf
		}
		return
	}
	rows.Scan(&auth.level)
	return
}

func GetTokenForUser(id int) (token string, err error) {
	var rows *sql.Rows
	rows, err = db.Query("SELECT token, expiration FROM auth_tokens WHERE "+
		"user_id = ?", id)
	defer rows.Close()

	// Check if token needs (re)generation
	var expiration time.Time
	generateToken := false
	if err != nil || !rows.Next() {
		generateToken = true
	}
	rows.Scan(&token, expiration)
	if expiration.Before(time.Now()) {
		generateToken = true
	}

	// Generate token, if needed
	if generateToken {
		tokenData := make([]byte, TokenSize, TokenSize)
		rand.Read(tokenData)
		token = hex.EncodeToString(tokenData)
		expiration = time.Now().AddDate(0, 1, 0)
		db.Exec("INSERT INTO auth_tokens(token, user_id, expiration)"+
			" VALUES(?, ?, ?) ON DUPLICATE KEY UPDATE token=?, expiration=?",
			token, id, expiration, token, expiration)
	}
	return
}
