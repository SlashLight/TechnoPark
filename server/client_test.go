package main

import (
	"net/http"
)

type TestCase struct {
	Limit      int
	Offset     int
	Query      string
	OrderField string
	OrderBy    int
	Result     *UserResult
	StatusCode string
}

type UserResult struct {
	Id        int
	Age       int
	FirstName string
	LastName  string
	Gender    string
	About     string
}

func ServerDummy(w http.ResponseWriter, r *http.Request) {

}
