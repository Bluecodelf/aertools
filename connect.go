package main

import (
	"encoding/json"
	"net/http"
)

type GETConnectPayload struct {
	ID      int    `json:"user_id"`
	Token   string `json:"token"`
	Level   int    `json:"auth_level"`
	Message string `json:"message"`
}

// Although it is a connection handler, this node uses the GET method.
// The user credentials may be sent in the "Authorization" header using the
// "Basic" authorization method, as defined in RFC 1945.
// It is also possible to connect using a user token using the custom "Token"
// authorization method. TODO: More documentation will be written on this later.
func HandlerGETConnect(writer http.ResponseWriter, request *http.Request) {
	var err error
	var payload GETConnectPayload
	var authInfo AuthInfo
	authInfo, err = GetAuthorizationLevel(request)
	payload.ID = authInfo.id
	payload.Token = authInfo.token
	payload.Level = authInfo.level
	if payload.Level == -1 {
		payload.Message = err.Error()
		writer.WriteHeader(401)
	} else {
		payload.Message = "Authorized"
	}
	data, _ := json.Marshal(&payload)
	writer.Write(data)
}
