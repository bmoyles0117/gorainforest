package rainforest

import (
	"testing"
)

func TestNewRainforest(t *testing.T) {
	rainforest := NewRainforest("ABC")

	if rainforest.ClientToken != "ABC" {
		t.Errorf("Unexpected client token: %s", rainforest.ClientToken)
	}

	if rainforest.client == nil {
		t.Error("Rainforest client was not assigned an http client")
	}
}
