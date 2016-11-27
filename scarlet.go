package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aerogo/scarlet"
	"github.com/fatih/color"
)

var cssCommentsRegex = regexp.MustCompile(`\/\*[^*]*\*+([^/*][^*]*\*+)*\/`)
var styleAnnouncePrefix = " " + color.YellowString("★") + " "

func scarletWork(job interface{}) interface{} {
	file := job.(string)
	scarletCode, _ := ReadFile(file)
	return scarletCode
	// css, err := scarlet.Compile(scarletCode, false)

	// if err != nil {
	// 	color.Red("Scarlet error:")
	// 	color.Red(err.Error())
	// }

	// return css
}

func announceStyle(name string) {
	fmt.Println(styleAnnouncePrefix, name)
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
		delete(styles, styleName)
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
