package rainforest

type TestBrowser struct {
	Name  string `json:"name"`
	State string `json:"state"`
}

type Test struct {
	Id               int            `json:"id"`
	Object           string         `json:"object"`
	CreatedAt        string         `json:"created_at"`
	EnvironmentId    int            `json:"environment_id"`
	State            string         `json:"state"`
	Result           string         `json:"result"`
	ExpectedWaitTime float32        `json:"expected_wait_time"`
	Browsers         []*TestBrowser `json:"browsers"`
	RequestedTests   []int          `json:"requested_tests"`
}
