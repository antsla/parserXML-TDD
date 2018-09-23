package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type FileParse struct {
	XMLName xml.Name    `xml:"root"`
	Users   []UserParse `xml:"row"`
}

type UserParse struct {
	Id     int    `xml:"id"`
	Name   string `xml:"first_name"`
	Age    int    `xml:"age"`
	About  string `xml:"about"`
	Gender string `xml:"gender"`
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")
	if query == "unmarshalError" {
		response, _ := json.Marshal(`{error_json}`)
		w.Write(response)
		return
	}

	db, err := os.Open("dataset.xml")
	defer db.Close()
	if err != nil {
		fmt.Println("open file error: ", err)
		return
	}

	file, err := ioutil.ReadAll(db)
	if err != nil {
		fmt.Println("cannot read xml: ", err)
		return
	}

	var parsed FileParse
	err = xml.Unmarshal(file, &parsed)
	if err != nil {
		fmt.Println("cannot parse xml structure: ", err)
		return
	}

	offset, errO := strconv.Atoi(r.FormValue("offset"))
	limit, errL := strconv.Atoi(r.FormValue("limit"))
	if errO != nil || errL != nil {
		fmt.Println("offset or limit parse error: ", errO, ", ", errL)
		return
	}
	if limit != 1 {
		response, _ := json.Marshal(parsed.Users[offset:limit])
		w.Write(response)
		return
	}

	response, err := json.Marshal(parsed.Users)
	w.Write(response)
}

func TimeoutErrorServer(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
}

func UnauthorizeErrorServer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
}

func InternalErrorServer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func BadRequestErrorServer(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")

	w.WriteHeader(http.StatusBadRequest)
	switch query {
	case "badReq_json":
		io.WriteString(w, `{"id": 10, "first_name": "John"`)
	case "badReq_BadOrderField":
		io.WriteString(w, `{"Error": "ErrorBadOrderField"}`)
	case "badReq_unknown":
		io.WriteString(w, `{"Error": "Unknown error"}`)
	}
}
