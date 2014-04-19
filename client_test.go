package rainforest

import (
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

	return t.res, nil
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

func TestRainforestRunTests(t *testing.T) {
	rainforest := NewRainforest("ABC")

	tr := &dummyTransport{
		res: &http.Response{
			Status:     "201 Created",
			StatusCode: 201,
			Proto:      "HTTP/1.0",
			Body:       ioutil.NopCloser(strings.NewReader(`{"id":1,"object":"Run","created_at":"2014-04-19T06:06:47Z","environment_id":1770,"state_log":[],"state":"queued","result":"no_result","expected_wait_time":8100.0,"browsers":[{"name":"chrome","state":"disabled"},{"name":"firefox","state":"disabled"},{"name":"ie8","state":"disabled"},{"name":"ie9","state":"disabled"},{"name":"safari","state":"disabled"}],"requested_tests":[1,2,3]}`)),
		},
	}

	rainforest.client = &http.Client{Transport: tr}

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
