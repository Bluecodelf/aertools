package main

import (
	"database/sql"
	"encoding/json"
	"github.com/Bluecodelf/rets"
	"net/http"
	"strconv"
	"strings"
)

type POSTLockersPayload struct {
	Login     string  `json:"login"`
	Locker    string  `json:"locker_name"`
	Borrowing float64 `json:"borrowing_time"`
}

type PUTLockerPayload struct {
	Login     string  `json:"login"`
	Locker    string  `json:"locker_name"`
	Borrowing float64 `json:"borrowing_time"`
	Retrieval float64 `json:"retrieval_time"`
	State     string  `json:"state"`
}

type GETLockersPayload struct {
	ID        string  `json:"id"`
	Login     string  `json:"login"`
	Locker    string  `json:"locker_name"`
	Borrowing float64 `json:"borrowing_time"`
	Retrieval float64 `json:"retrieval_time"`
	State     string  `json:"state"`
}

func HandlerPOSTLockers(writer http.ResponseWriter, request *http.Request) {
	// Check if the user is authenticated as AER (level 1) or superior.
	auth, err := GetAuthorizationLevel(request)
	if err != nil || auth.id == -1 || auth.level < 1 {
		rets.HandlerHTTPUnauthorized(writer)
		return
	}

	// Retrieve the payload from the request, note that the POST payload can't
	// possess a retrieval or state field as those values are supposed to be
	// updated upon the retrieval of the locker.
	// TODO: Adding a RegEx-based check for login and locker would be nice.
	// Checking the validity of the borrowing time is essential too.
	payload := new(POSTLockersPayload)
	err = rets.UnmarshalHTTPBody(request, payload)
	if err != nil || payload.Login == "" || payload.Locker == "" {
		rets.HandlerHTTPBadRequest(writer)
		return
	}

	_, err = db.Query("INSERT INTO lockers(login, locker, borrowing) VALUES"+
		"(?, ?, ?)", payload.Login, payload.Locker, payload.Borrowing)
	if err != nil {
		rets.HandlerError(writer, err)
	} else {
		rets.HandlerHTTPOK(writer)
	}
	return
}

func HandlerPUTLocker(writer http.ResponseWriter, request *http.Request) {
	// Same permission check as POST /lockers
	auth, err := GetAuthorizationLevel(request)
	if err != nil || auth.id == -1 || auth.level < 1 {
		rets.HandlerHTTPUnauthorized(writer)
		return
	}

	// Retrieve the second node in the URL path, representing the locker ID
	var id float64
	id, err = strconv.ParseFloat(strings.Split(
		request.URL.Path[1:], "/")[1], 64)
	if err != nil {
		rets.HandlerError(writer, err)
		return
	}

	// TODO: some checks to do here as well
	payload := new(PUTLockerPayload)
	err = rets.UnmarshalHTTPBody(request, payload)
	if err != nil {
		rets.HandlerHTTPBadRequest(writer)
		return
	}

	_, err = db.Query("UPDATE lockers SET login=?, locker=?, borrowing=?,"+
		"retrieval=?, state=? WHERE id=?", payload.Login, payload.Locker,
		payload.Borrowing, payload.Retrieval, payload.State, id)
	if err != nil {
		rets.HandlerError(writer, err)
	} else {
		rets.HandlerHTTPOK(writer)
	}
	return
}

func HandlerGETLockers(writer http.ResponseWriter, request *http.Request) {
	// Same permissions...
	auth, err := GetAuthorizationLevel(request)
	if err != nil || auth.id == -1 || auth.level < 1 {
		rets.HandlerHTTPUnauthorized(writer)
		return
	}

	var rows *sql.Rows
	rows, err = db.Query("SELECT * FROM lockers")
	defer rows.Close()
	if err != nil {
		rets.HandlerError(writer, err)
	}

	// REVIEW: The append function actually allocates in power of two. It could
	// get messy really fast for semi-large requests. We may have to do
	// performance tests to know if implementing a home-made allocator is
	// better or not.
	var lockers []GETLockersPayload
	for rows.Next() {
		var locker GETLockersPayload
		rows.Scan(&locker.ID, &locker.Login, &locker.Locker,
			&locker.Borrowing, &locker.Retrieval, &locker.State)
		lockers = append(lockers, locker)
	}
	data, _ := json.Marshal(&lockers)
	writer.Write(data)
}
