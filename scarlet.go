package main

import (
	"fmt"
	"path"
	"regexp"
	"strings"

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

	// Write CSS bundle into $.css.go where it can be referenced as components.CSS
	EmbedData(path.Join(outputFolder, "$.css.go"), pixy.PackageName, "CSS", bundledCSS)
}

func getBundledCSS(styles map[string]string) string {
	scarletCodes := []string{}

	// Ordered styles
	for _, styleName := range app.Config.Styles {
		styleName = "styles/" + styleName + ".scarlet"
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
