package rainforest

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type Rainforest struct {
	ClientToken string
	client      *http.Client
}

func (r *Rainforest) RunTests(test_ids []int) (*Test, error) {
	var (
		data []byte
		err  error
		req  *http.Request
		res  *http.Response
	)

	if data, err = json.Marshal(map[string]interface{}{"tests": test_ids}); err != nil {
		return nil, err
	}

	if req, err = http.NewRequest("POST", "https://app.rainforestqa.com/api/1/runs", bytes.NewReader(data)); err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("CLIENT_TOKEN", r.ClientToken)

	if res, err = r.client.Do(req); err != nil {
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

func NewRainforest(client_token string) *Rainforest {
	return &Rainforest{
		ClientToken: client_token,
		client:      http.DefaultClient,
	}
}
