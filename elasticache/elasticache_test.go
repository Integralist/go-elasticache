package elasticache

import (
	"bytes"
	"os"
	"testing"
)

func TestElastiCacheEndpoint(t *testing.T) {
	expectation := "foo"
	os.Setenv("ELASTICACHE_ENDPOINT", expectation)
	response, _ := elasticache()

	if response != expectation {
		t.Errorf("The response '%s' didn't match the expectation '%s'", response, expectation)
	}
}

func TestParseNodes(t *testing.T) {
	expectation := "localhost|127.0.0.1|11211"
	cluster := `CONFIG cluster 0 25
1
localhost|127.0.0.1|11211

END`

	r := bytes.NewReader([]byte(cluster))

	response, _ := parseNodes(r)

	if response != expectation {
		t.Errorf("The response '%s' didn't match the expectation '%s'", response, expectation)
	}
}

func TestParseURLs(t *testing.T) {
	expectationLength := 3

	response, _ := parseURLs("host|foo|1 host|bar|2 host|baz|3")

	if len(response) != expectationLength {
		t.Errorf("The response length '%d' didn't match the expectation '%d'", len(response), expectationLength)
	}

	var suite = []struct {
		response    string
		expectation string
	}{
		{response[0], "foo:1"},
		{response[1], "bar:2"},
		{response[2], "baz:3"},
	}

	for _, v := range suite {
		if v.response != v.expectation {
			t.Errorf("The response '%s' didn't match the expectation '%s'", v.response, v.expectation)
		}
	}
}
