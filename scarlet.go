package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"strings"

	"github.com/aerogo/aero"
	"github.com/aerogo/pixy"
	"github.com/aerogo/scarlet"
	"github.com/fatih/color"
)

var cssCommentsRegex = regexp.MustCompile(`\/\*[^*]*\*+([^/*][^*]*\*+)*\/`)
var styleAnnouncePrefix = " " + color.YellowString("★") + " "

func scarletWork(job interface{}) interface{} {
	file := job.(string)
	scarletCode, _ := ReadFile(file)
	return scarletCode
}

func scarletFinish(results WorkerPoolResults) {
	// Convert to map[string]string
	styles := ToStringMap(results)

	// Bundled CSS
	bundledCSS := getBundledCSS(styles)

	// Encode in Base64
	bundledCSS = base64.StdEncoding.EncodeToString(aero.StringToBytesUnsafe(bundledCSS))

	// Create Go code to load the embedded CSS
	cssCode := "package " + pixy.PackageName + "\n\nimport \"encoding/base64\"\n\n// CSS ...\nfunc CSS() string {\ncssEncoded := `\n" + bundledCSS + "\n`\ncssDecoded, _ := base64.StdEncoding.DecodeString(cssEncoded)\nreturn string(cssDecoded)\n}\n"

	// Write the loader to $.css.go
	ioutil.WriteFile(path.Join(outputFolder, "$.css.go"), aero.StringToBytesUnsafe(cssCode), 0644)
}

func getBundledCSS(styles map[string]string) string {
	scarletCodes := []string{}

	// Ordered styles
	for _, styleName := range app.Config.Styles {
		styleName = "styles/" + styleName + scarletExtension
		styleContent, exists := styles[styleName]

		if exists {
			announceStyle(styleName)
			scarletCodes = append(scarletCodes, styleContent)
			delete(styles, styleName)
		}
	}

	// Unordered styles in styles directory
	for styleName, styleContent := range styles {
		if strings.HasPrefix(styleName, "styles/") {
			announceStyle(styleName)
			scarletCodes = append(scarletCodes, styleContent)
			delete(styles, styleName)
		}
	}

	// Unordered styles
	for styleName, styleContent := range styles {
		announceStyle(styleName)
		scarletCodes = append(scarletCodes, styleContent)
	}

	allScarletCodes := strings.Join(scarletCodes, "\n")
	css, err := scarlet.Compile(allScarletCodes, false)

	if err != nil {
		color.Red(err.Error())
		return ""
	}

	// Fonts
	fontsCSS := <-fontsCSSChannel

	return fontsCSS + css
}

func announceStyle(name string) {
	fmt.Println(styleAnnouncePrefix, name)
}
