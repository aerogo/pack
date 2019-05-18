package scarletpacker

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aerogo/http/client"
)

// This is the user agent we use to fetch the CSS from Google Fonts.
// Apparently you get different results if the user agent is empty,
// so we're faking a standard Chrome user agent.
const fontsUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.157 Safari/537.36"

// This regex will match all CSS multiline comments.
var cssCommentsRegex = regexp.MustCompile(`\/\*[^*]*\*+([^/*][^*]*\*+)*\/`)

// downloadFonts downloads the given Google fonts.
func downloadFonts(fonts []string) (string, error) {
	url := fmt.Sprintf("https://fonts.googleapis.com/css?family=%s&display=swap", strings.Join(fonts, "|"))
	response, err := client.Get(url).Header("User-Agent", fontsUserAgent).End()

	if err != nil {
		return "", err
	}

	if !response.Ok() {
		return "", fmt.Errorf("Fetching %s resulted in status %d", url, response.StatusCode())
	}

	css := response.String()
	css = strings.ReplaceAll(css, "\r", "")
	css = strings.ReplaceAll(css, "\n", "")
	css = strings.ReplaceAll(css, "  ", " ")
	css = strings.ReplaceAll(css, "{ ", "{")
	css = strings.ReplaceAll(css, " {", "{")
	css = strings.ReplaceAll(css, ": ", ":")
	css = strings.ReplaceAll(css, "; ", ";")
	css = strings.ReplaceAll(css, ", ", ",")

	// Remove CSS comments
	css = cssCommentsRegex.ReplaceAllString(css, "")

	return css, nil
}
