package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type TestCase struct {
	TestClient  SearchClient
	TestRequest SearchRequest
	Result      *SearchResponse
	StatusCode  int
	Error       error
}

func TestSearchingClients(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	cases := []TestCase{
		{
			TestClient: SearchClient{
				AccessToken: "Good",
				URL:         ts.URL,
			},
			TestRequest: SearchRequest{
				Limit:      -5,
				Offset:     0,
				Query:      "Boyd",
				OrderField: "Name",
				OrderBy:    0,
			},
			Result:     nil,
			StatusCode: 0,
			Error:      fmt.Errorf("limit must be > 0"),
		},
		{
			TestClient: SearchClient{
				AccessToken: "Good",
				URL:         ts.URL,
			},
			TestRequest: SearchRequest{
				Limit:      5,
				Offset:     -7,
				Query:      "Boyd",
				OrderField: "Name",
				OrderBy:    0,
			},
			Result:     nil,
			StatusCode: 0,
			Error:      fmt.Errorf("offset must be > 0"),
		},
		{
			TestClient: SearchClient{
				AccessToken: "Good",
				URL:         ts.URL,
			},
			TestRequest: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "Boyd",
				OrderField: "Name",
				OrderBy:    5,
			},
			Result:     nil,
			StatusCode: 400,
			Error:      fmt.Errorf("unknown bad request error: OrderBy should be 1, 0 or -1"),
		},
		{
			TestClient: SearchClient{
				AccessToken: "Good",
				URL:         ts.URL,
			},
			TestRequest: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "Boyd",
				OrderField: "Greek",
				OrderBy:    1,
			},
			Result:     nil,
			StatusCode: 400,
			Error:      fmt.Errorf("OrderField Greek invalid"),
		},
		{
			TestClient: SearchClient{
				AccessToken: "good",
				URL:         ts.URL,
			},
			TestRequest: SearchRequest{
				Limit:      1,
				Offset:     1,
				Query:      "a",
				OrderField: "Name",
				OrderBy:    1,
			},
			Result: &SearchResponse{
				Users: []User{
					{ID: 1, Name: "Hilda Mayer", Age: 21, About: "Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.", Gender: "female"},
				},
				NextPage: false,
			},
			StatusCode: 200,
			Error:      nil,
		},
		{
			TestClient: SearchClient{
				AccessToken: "Bad",
				URL:         ts.URL,
			},
			TestRequest: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "Boyd",
				OrderField: "Name",
				OrderBy:    1,
			},
			Result:     nil,
			StatusCode: 401,
			Error:      fmt.Errorf("bad AccessToken"),
		},
	}
	//fmt.Println(cases[4].Result)
	for caseNum, item := range cases {
		tClient := item.TestClient
		tReq := item.TestRequest
		res, err := tClient.FindUsers(tReq)
		if err == nil && item.Error != nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if err != nil && item.Error == nil {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if !reflect.DeepEqual(item.Result, res) {
			t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, item.Result, res)
		}
	}
	ts.Close()

}
