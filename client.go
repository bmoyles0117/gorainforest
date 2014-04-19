package rainforest

import (
	"net/http"
)

type Rainforest struct {
	ClientToken string
	client      *http.Client
}

func NewRainforest(client_token string) *Rainforest {
	return &Rainforest{
		ClientToken: client_token,
		client:      http.DefaultClient,
	}
}
