package main

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCase struct {
	Request *SearchRequest
}

var AccessToken string = fmt.Sprint(md5.Sum([]byte("AccessToken")))

func TestRequestErrorLimitAndOffset(t *testing.T) {
	searchClient := &SearchClient{}
	testCases := []TestCase{
		TestCase{
			Request: &SearchRequest{
				Limit: -1,
			},
		},
		TestCase{
			Request: &SearchRequest{
				Offset: -1,
			},
		},
	}

	testingServer := httptest.NewServer(http.HandlerFunc(SearchServer))
	for _, testCase := range testCases {
		_, err := searchClient.FindUsers(*testCase.Request)
		if err == nil {
			t.Errorf("[TestRequestErrorLimitAndOffset] expected error, got nil")
		}
	}
	testingServer.Close()
}

func TestTimeoutError(t *testing.T) {
	testCase := TestCase{
		Request: &SearchRequest{},
	}

	testingServer := httptest.NewServer(http.HandlerFunc(TimeoutErrorServer))
	searchClient := &SearchClient{
		URL: testingServer.URL,
	}

	_, err := searchClient.FindUsers(*testCase.Request)
	if err == nil {
		t.Errorf("[TestTimeoutError] expected error, got nil")
	}

	testingServer.Close()
}

func TestUnknownServerError(t *testing.T) {
	testCase := TestCase{
		Request: &SearchRequest{},
	}

	testingServer := httptest.NewServer(http.HandlerFunc(SearchServer))
	searchClient := &SearchClient{
		URL: "errorUrl",
	}

	_, err := searchClient.FindUsers(*testCase.Request)
	if err == nil {
		t.Errorf("[TestUnknownServerError] expected error, got nil")
	}

	testingServer.Close()
}

func TestUnauthorizeError(t *testing.T) {
	testCase := TestCase{
		Request: &SearchRequest{},
	}

	testingServer := httptest.NewServer(http.HandlerFunc(UnauthorizeErrorServer))
	searchClient := &SearchClient{
		AccessToken: "",
		URL:         testingServer.URL,
	}

	_, err := searchClient.FindUsers(*testCase.Request)
	if err == nil {
		t.Errorf("[TestUnauthorizeError] expected error, got nil")
	}

	testingServer.Close()
}

func TestInternalServerError(t *testing.T) {
	testCase := TestCase{
		Request: &SearchRequest{},
	}

	testingServer := httptest.NewServer(http.HandlerFunc(InternalErrorServer))

	searchClient := &SearchClient{
		AccessToken: AccessToken,
		URL:         testingServer.URL,
	}
	_, err := searchClient.FindUsers(*testCase.Request)
	if err == nil {
		t.Errorf("[TestInternalServerError] expected error, got nil")
	}

	testingServer.Close()
}

func TestBadRequestError(t *testing.T) {
	testCases := []TestCase{
		TestCase{
			Request: &SearchRequest{
				Query: "badReq_json",
			},
		},
		TestCase{
			Request: &SearchRequest{
				Query: "badReq_BadOrderField",
			},
		},
		TestCase{
			Request: &SearchRequest{
				Query: "badReq_unknown",
			},
		},
	}

	testingServer := httptest.NewServer(http.HandlerFunc(BadRequestErrorServer))

	searchClient := &SearchClient{
		AccessToken: AccessToken,
		URL:         testingServer.URL,
	}

	for _, testCase := range testCases {
		_, err := searchClient.FindUsers(*testCase.Request)
		if err == nil {
			t.Errorf("[TestInternalServerError] expected error, got nil")
		}
	}

	testingServer.Close()
}

func TestUnmarshJson(t *testing.T) {
	testCase := TestCase{
		Request: &SearchRequest{
			Query: "unmarshalError",
		},
	}

	testingServer := httptest.NewServer(http.HandlerFunc(SearchServer))
	searchClient := &SearchClient{
		AccessToken: AccessToken,
		URL:         testingServer.URL,
	}

	_, err := searchClient.FindUsers(*testCase.Request)
	if err == nil {
		t.Errorf("[TestUnmarshJson] expected error, got nil")
	}

	testingServer.Close()
}

func TestGetUsers(t *testing.T) {
	testCases := []TestCase{
		TestCase{
			Request: &SearchRequest{},
		},
		TestCase{
			Request: &SearchRequest{
				Limit: 26,
			},
		},
		TestCase{
			Request: &SearchRequest{
				Limit:  6,
				Offset: 0,
			},
		},
	}

	testingServer := httptest.NewServer(http.HandlerFunc(SearchServer))
	searchClient := &SearchClient{
		AccessToken: AccessToken,
		URL:         testingServer.URL,
	}

	for _, testCase := range testCases {
		_, err := searchClient.FindUsers(*testCase.Request)
		if err != nil {
			t.Errorf("[TestGetUsers] expected nil, got error: %s", err)
		}
	}

	testingServer.Close()
}
