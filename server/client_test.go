package main

import (
	"testing"
)

type TestCase struct {
	TestClient	*SearchClient
	TestRequest *SearchRequest
	Result 		*SearchResponse
	StatusCode 	int
	IsError bool
}

func TestSearchingClients(t *testing.T) {
	cases := []TestCase{
		{
			TestClient: &SearchClient {
				AccessToken: "Good",
				URL:         "",
			},
			TestRequest: &SearchRequest {
				Limit:      -5,
				Offset:     0,
				Query:      "Boyd",
				OrderField: "Name",
				OrderBy:    0,
			},
			Result: nil,
			StatusCode: 0,
			IsError: true,
,		},
		{
			TestClient: &SearchClient {
				AccessToken: "Good",
				URL:         "",
			},
			TestRequest: &SearchRequest {
				Limit:      5,
				Offset:     -7,
				Query:      "Boyd",
				OrderField: "Name",
				OrderBy:    0,
			},
			Result: nil,
			StatusCode: 0,
			IsError: true,

		},
		{
			TestClient: &SearchClient {
				AccessToken: "Good",
				URL:         "",
			},
			TestRequest: &SearchRequest {
				Limit:      1,
				Offset:     0,
				Query:      "Boyd",
				OrderField: "Name",
				OrderBy:    5,
			},
			Result: nil,
			StatusCode: 400,
			IsError: true,
		},
		{
			TestClient: &SearchClient {
				AccessToken: "Good",
				URL:         "",
			},
			TestRequest: &SearchRequest {
				Limit:      1,
				Offset:     0,
				Query:      "Boyd",
				OrderField: "Name",
				OrderBy:    5,
			},
			Result: nil,
			StatusCode: 0,
			IsError: true,
		},

	}
	{

	}
}
