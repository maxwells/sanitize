package sanitize

import (
	"testing"
	"strings"
	"fmt"
)

// presently failing
func TestSanitizeRemoveKeepsDoctype(t *testing.T) {
	htmlDoc := "<!DOCTYPE html><html><head></head><body><div></body></html>"
	config := `
	{
		"elements": {
			"html": [],
			"head": [],
			"body": []
		}
	}`

	whitelist, _ := NewWhitelist([]byte(config))
	output, _ := whitelist.SanitizeRemove(strings.NewReader(htmlDoc))

	if output != "<!DOCTYPE><html><head></head><body></body></html>" {
		t.Error("failed")
	}
}