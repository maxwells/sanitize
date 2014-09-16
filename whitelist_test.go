package sanitize

import (
	"testing"
	"strings"
)

func TestSanitizeRendersDoctypeCorrectly(t *testing.T) {
	htmlDoc := `<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">`
	expectedOutput := `<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">`
	config := `
	{
		"elements": {
		}
	}`

	whitelist, _ := NewWhitelist([]byte(config))
	output, _ := whitelist.SanitizeRemove(strings.NewReader(htmlDoc))

	if output != expectedOutput {
		t.Errorf("failed: %s != %s", output, expectedOutput)
	}
}

func TestSanitizeRemoveRemovesNonWhitelistedNodes(t *testing.T) {
	htmlDoc := `<!DOCTYPE html>
				<html>
					<head>
						<title>My Title</title>
					</head>
					<body>
						<div>
							<b>Bold</b>
							<i>Italic</i>
							<em>Emphatic</em>
						</div>
					</body>
				</html>`
	expectedOutput := `<!DOCTYPE html><html><head><title>My Title</title></head><body><div><i>Italic</i></div></body></html>`
	config := `{
		"stripWhitespace": true,
		"elements": {
			"html": [],
			"head": [],
			"title": [],
			"body": [],
			"div": [],
			"i": []
		}
	}`

	whitelist, _ := NewWhitelist([]byte(config))
	output, _ := whitelist.SanitizeRemove(strings.NewReader(htmlDoc))

	if output != expectedOutput {
		t.Errorf("failed: %s != %s", output, expectedOutput)
	}
}

func TestStripWhitespace(t *testing.T) {
	htmlDoc := `<!DOCTYPE html>
				<html>
					<head>
					</head>
					<body>
					</body>
				</html>`
	expectedOutput := `<!DOCTYPE html><html><head></head><body></body></html>`
	config := `{
		"stripWhitespace": true,
		"elements": {
			"html": [],
			"head": [],
			"body": []
		}
	}`

	whitelist, _ := NewWhitelist([]byte(config))
	output, _ := whitelist.SanitizeRemove(strings.NewReader(htmlDoc))

	if output != expectedOutput {
		t.Errorf("failed: %s != %s", output, expectedOutput)
	}
}

func TestStripComments(t *testing.T) {
	htmlDoc := `<!DOCTYPE html><!-- hello world --><html></html>`
	expectedOutput := `<!DOCTYPE html><html></html>`
	config := `{
		"stripComments": true,
		"elements": {
			"html": []
		}
	}`

	whitelist, _ := NewWhitelist([]byte(config))
	output, _ := whitelist.SanitizeRemove(strings.NewReader(htmlDoc))

	if output != expectedOutput {
		t.Errorf("failed: %s != %s", output, expectedOutput)
	}
}