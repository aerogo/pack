package main

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/aerogo/aero"
)

const fontsUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.71 Safari/537.36"

func downloadFontsCSS(fonts []string) string {
	fontsCSS, err := aero.Get("https://fonts.googleapis.com/css?family="+strings.Join(fonts, "|")).Header("User-Agent", fontsUserAgent).Send()

	if err != nil {
		return ""
	}

	fontsCSS = strings.Replace(fontsCSS, "\r", "", -1)
	fontsCSS = strings.Replace(fontsCSS, "\n", "", -1)
	fontsCSS = strings.Replace(fontsCSS, "  ", " ", -1)
	fontsCSS = strings.Replace(fontsCSS, "{ ", "{", -1)
	fontsCSS = strings.Replace(fontsCSS, ": ", ":", -1)
	fontsCSS = strings.Replace(fontsCSS, "; ", ";", -1)
	fontsCSS = strings.Replace(fontsCSS, ", ", ",", -1)

	// Remove CSS comments
	fontsCSS = cssCommentsRegex.ReplaceAllString(fontsCSS, "")

	// Save in cache
	ioutil.WriteFile(path.Join(cacheFolder, "fonts", strings.Join(app.Config.Fonts, "|")+".css"), []byte(fontsCSS), 0777)

	return fontsCSS
}
