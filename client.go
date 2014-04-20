package rainforest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const (
	ALL_TESTS = "all"

	RESULT_NONE   = "no_result"
	RESULT_PASSED = "passed"
	RESULT_FAILED = "failed"

	STATE_ENABLED     = "enabled"
	STATE_DISABLED    = "disabled"
	STATE_QUEUED      = "queued"
	STATE_VALIDATING  = "validating"
	STATE_IN_PROGRESS = "in_progress"
	STATE_PASSED      = "passed"
	STATE_FAILED      = "failed"
	STATE_COMPLETE    = "complete"
)

var InvalidTestIds = errors.New("Invalid test IDs passed, must be a string or array of ints")

// Rainforest stores the ClientToken and a few abstraction methods for
// implementing the API.
type Rainforest struct {
	ClientToken string
	client      *http.Client
}

func (r *Rainforest) doRequest(method, path string, body io.Reader) (*http.Response, error) {
	var (
		err error
		req *http.Request
	)

	if req, err = http.NewRequest(method, "https://app.rainforestqa.com/api/1"+path, body); err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("CLIENT_TOKEN", r.ClientToken)

	return r.client.Do(req)
}

// RunTests takes either a string, or an array of ints. The only string that
// should be passed for the time being is ALL_TESTS ("all"). When executed
// properly, a Test will be returned that represents the current state of the
// tests that were executed.
func (r *Rainforest) RunTests(test_filter interface{}) (*Test, error) {
	var (
		data []byte
		err  error
		res  *http.Response
	)

	if test_ids, ok := test_filter.([]int); ok {
		if data, err = json.Marshal(map[string]interface{}{"tests": test_ids}); err != nil {
			return nil, err
		}
	} else if test_criteria, ok := test_filter.(string); ok {
		if data, err = json.Marshal(map[string]interface{}{"tests": test_criteria}); err != nil {
			return nil, err
		}
	} else {
		return nil, InvalidTestIds
	}

	if res, err = r.doRequest("POST", "/runs", bytes.NewReader(data)); err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == 201 {
		test := &Test{}

		if err = json.NewDecoder(res.Body).Decode(test); err != nil {
			return nil, err
		}

		return test, nil
	} else {
		error_response := make(map[string]string)

		if err = json.NewDecoder(res.Body).Decode(&error_response); err != nil {
			return nil, err
		}

		return nil, errors.New(error_response["error"])
	}
}

// Generate a new Rainforest client with the client token specified, and
// associated to the default http client for executing requests.
func NewRainforest(client_token string) *Rainforest {
	return &Rainforest{
		ClientToken: client_token,
		client:      http.DefaultClient,
	}
}
