package rainforest

// A TestBrowser simply represents an element of the parent's
// Browsers array to convey the browser and the current state
// of that browser's testing.
type TestBrowser struct {
	Name  string `json:"name"`
	State string `json:"state"`
}

// A Test is returned when running tests, that represents the
// current state of the test being run as well as a few other
// meta data related parameters.
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
