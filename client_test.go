package rainforest

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type dummyTransport struct {
	req *http.Request
	res *http.Response
}

func (t *dummyTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	t.req = req

	if t.res != nil {
		return t.res, nil
	} else {
		return nil, errors.New("Dummy error")
	}
}

func TestNewRainforest(t *testing.T) {
	rainforest := NewRainforest("ABC")

	if rainforest.ClientToken != "ABC" {
		t.Errorf("Unexpected client token: %s", rainforest.ClientToken)
	}

	if rainforest.client == nil {
		t.Error("Rainforest client was not assigned an http client")
	}
}

func TestRainforestDoRequest(t *testing.T) {
	var (
		data       []byte
		err        error
		rainforest *Rainforest
	)

	rainforest = NewRainforest("ABC")

	tr := &dummyTransport{}

	rainforest.client = &http.Client{Transport: tr}

	// We don't care about the actual response here, just that the request was mutated
	rainforest.doRequest("POST", "/runs", strings.NewReader("hello"))

	if h := tr.req.Header.Get("Accept"); h != "application/json" {
		t.Errorf("Accept header was not set properly: %s", h)
	}

	if h := tr.req.Header.Get("Content-type"); h != "application/json" {
		t.Errorf("Content-type header was not set properly: %s", h)
	}

	if h := tr.req.Header.Get("CLIENT_TOKEN"); h != "ABC" {
		t.Errorf("CLIENT_TOKEN header was not set properly: %s", h)
	}

	if data, err = ioutil.ReadAll(tr.req.Body); err != nil {
		t.Errorf("Unexpected error reading the request's body: %s", err)
	}

	if string(data) != "hello" {
		t.Errorf("Unexpected value stored in the request's body: %s", string(data))
	}

	if tr.req.URL.RequestURI() != "/api/1/runs" {
		t.Errorf("Unexpected request URI: %s", tr.req.URL.RequestURI())
	}
}

func TestRainforestRunTests(t *testing.T) {
	var (
		err        error
		rainforest *Rainforest
	)

	rainforest = NewRainforest("ABC")

	tr := &dummyTransport{}

	rainforest.client = &http.Client{Transport: tr}

	tr.res = &http.Response{
		Status:     "403 Forbidden",
		StatusCode: 403,
		Proto:      "HTTP/1.0",
		Body:       ioutil.NopCloser(strings.NewReader(`{"error":"Invalid test ids"}`)),
	}

	if _, err = rainforest.RunTests([]int{1, 2, 3}); err == nil {
		t.Error("Expected to receive an error according to the mocked response")
	}

	if err.Error() != "Invalid test ids" {
		t.Errorf("Unexpected error: %s", err)
	}

	tr.res = &http.Response{
		Status:     "404 Not Found",
		StatusCode: 404,
		Proto:      "HTTP/1.0",
		Body:       ioutil.NopCloser(strings.NewReader(`{"error":"Account not found"}`)),
	}

	if _, err = rainforest.RunTests([]int{1, 2, 3}); err == nil {
		t.Error("Expected to receive an error according to the mocked response")
	}

	if err.Error() != "Account not found" {
		t.Errorf("Unexpected error: %s", err)
	}

	tr.res = &http.Response{
		Status:     "201 Created",
		StatusCode: 201,
		Proto:      "HTTP/1.0",
		Body:       ioutil.NopCloser(strings.NewReader(`{"id":1,"object":"Run","created_at":"2014-04-19T06:06:47Z","environment_id":1770,"state_log":[],"state":"queued","result":"no_result","expected_wait_time":8100.0,"browsers":[{"name":"chrome","state":"disabled"},{"name":"firefox","state":"disabled"},{"name":"ie8","state":"disabled"},{"name":"ie9","state":"disabled"},{"name":"safari","state":"disabled"}],"requested_tests":[1,2,3]}`)),
	}

	// Ensure that we can run all tests, and set the body correctly
	rainforest.RunTests(ALL_TESTS)

	if h := tr.req.Header.Get("Content-type"); h != "application/json" {
		t.Errorf("Content-type header was not set properly: %s", h)
	}

	if h := tr.req.Header.Get("Accept"); h != "application/json" {
		t.Errorf("Accept header was not set properly: %s", h)
	}

	if h := tr.req.Header.Get("CLIENT_TOKEN"); h != "ABC" {
		t.Errorf("CLIENT_TOKEN header was not set properly: %s", h)
	}

	if data, err := ioutil.ReadAll(tr.req.Body); err == nil {
		if string(data) != `{"tests":"all"}` {
			t.Errorf("Unexpected data stored in the request body: %s", string(data))
		}
	} else {
		t.Errorf("Unexpected error reading request body: %s", err)
	}

	// Run a selective few tests, and ensure the parameters are set properly
	tr.res = &http.Response{
		Status:     "201 Created",
		StatusCode: 201,
		Proto:      "HTTP/1.0",
		Body:       ioutil.NopCloser(strings.NewReader(`{"id":1,"object":"Run","created_at":"2014-04-19T06:06:47Z","environment_id":1770,"state_log":[],"state":"queued","result":"no_result","expected_wait_time":8100.0,"browsers":[{"name":"chrome","state":"disabled"},{"name":"firefox","state":"disabled"},{"name":"ie8","state":"disabled"},{"name":"ie9","state":"disabled"},{"name":"safari","state":"disabled"}],"requested_tests":[1,2,3]}`)),
	}

	test, _ := rainforest.RunTests([]int{1, 2, 3})

	if h := tr.req.Header.Get("Content-type"); h != "application/json" {
		t.Errorf("Content-type header was not set properly: %s", h)
	}

	if h := tr.req.Header.Get("Accept"); h != "application/json" {
		t.Errorf("Accept header was not set properly: %s", h)
	}

	if h := tr.req.Header.Get("CLIENT_TOKEN"); h != "ABC" {
		t.Errorf("CLIENT_TOKEN header was not set properly: %s", h)
	}

	if data, err := ioutil.ReadAll(tr.req.Body); err == nil {
		if string(data) != `{"tests":[1,2,3]}` {
			t.Errorf("Unexpected data stored in the request body: %s", string(data))
		}
	} else {
		t.Errorf("Unexpected error reading request body: %s", err)
	}

	if fmt.Sprintf("%d", test.RequestedTests) != fmt.Sprintf("%d", []int{1, 2, 3}) {
		t.Errorf("Unexpected test ids returned: %d", test.RequestedTests)
	}

	for i, browser_name := range []string{"chrome", "firefox", "ie8", "ie9", "safari"} {
		if test.Browsers[i].Name != browser_name {
			t.Errorf("Unexpected browser at sequence (%d): %s", i, test.Browsers[i].Name)
		}
	}
}
