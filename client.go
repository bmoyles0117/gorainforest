package rainforest

type Rainforest struct {
	ClientToken string
}

func NewRainforest(client_token string) *Rainforest {
	return &Rainforest{client_token}
}
