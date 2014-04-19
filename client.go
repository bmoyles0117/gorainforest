package rainforest

import (
	"bytes"
	"encoding/json"
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

	test := &Test{}

	if err = json.NewDecoder(res.Body).Decode(test); err != nil {
		return nil, err
	}

	return test, nil
}

func NewRainforest(client_token string) *Rainforest {
	return &Rainforest{
		ClientToken: client_token,
		client:      http.DefaultClient,
	}
}
