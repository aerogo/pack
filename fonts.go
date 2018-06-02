package main

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/aerogo/http/client"
)

var fontsCSSChannel = make(chan string, 1)

func fontsInit() {
	// Cache folder
	os.Mkdir(path.Join(cacheFolder, "fonts"), 0777)

	go func() {
		if len(config.Fonts) > 0 {
			cached, err := ReadFile(path.Join(cacheFolder, "fonts", strings.Join(config.Fonts, "|")+".css"))

			if err == nil {
				fontsCSSChannel <- cached
			} else {
				fontsCSSChannel <- downloadFontsCSS(config.Fonts)
			}
		} else {
			fontsCSSChannel <- ""
		}
	}()
}

func downloadFontsCSS(fonts []string) string {
	const fontsUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.167 Safari/537.36"

	url := "https://fonts.googleapis.com/css?family=" + strings.Join(fonts, "|")
	response, err := client.Get(url).Header("User-Agent", fontsUserAgent).End()

	if err != nil {
		panic(err)
	}

	fontsCSS := response.String()
	fontsCSS = strings.Replace(fontsCSS, "\r", "", -1)
	fontsCSS = strings.Replace(fontsCSS, "\n", "", -1)
	fontsCSS = strings.Replace(fontsCSS, "  ", " ", -1)
	fontsCSS = strings.Replace(fontsCSS, "{ ", "{", -1)
	fontsCSS = strings.Replace(fontsCSS, " {", "{", -1)
	fontsCSS = strings.Replace(fontsCSS, ": ", ":", -1)
	fontsCSS = strings.Replace(fontsCSS, "; ", ";", -1)
	fontsCSS = strings.Replace(fontsCSS, ", ", ",", -1)

	// Remove CSS comments
	fontsCSS = cssCommentsRegex.ReplaceAllString(fontsCSS, "")

	// If font-display is not set, set it:
	if !strings.Contains(fontsCSS, "font-display") {
		fontsCSS = strings.Replace(fontsCSS, "@font-face{", "@font-face{font-display:swap;", -1)
	}

	// Save in cache
	ioutil.WriteFile(path.Join(cacheFolder, "fonts", strings.Join(config.Fonts, "|")+".css"), []byte(fontsCSS), 0777)

	return fontsCSS
}
