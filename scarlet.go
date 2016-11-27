package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aerogo/scarlet"
	"github.com/fatih/color"
)

var cssCommentsRegex = regexp.MustCompile(`\/\*[^*]*\*+([^/*][^*]*\*+)*\/`)

func scarletWork(job interface{}) interface{} {
	file := job.(string)
	// fmt.Println(" "+color.GreenString("☼"), file)

	scarletCode, _ := ReadFile(file)
	css, err := scarlet.Compile(scarletCode, false)

	if err != nil {
		color.Red("Scarlet error:")
		color.Red(err.Error())
	}

	return css
}

func getBundledCSS(styles map[string]string) string {
	css := []string{}

	// Ordered styles
	for _, styleName := range app.Config.Styles {
		styleName = "styles/" + styleName + scarletExtension
		styleContent := styles[styleName]

		if styleContent != "" {
			fmt.Println(" "+color.GreenString("☼"), styleName)
			css = append(css, styleContent)
			styles[styleName] = ""
		}
	}

	// Unordered styles in styles directory
	for styleName, styleContent := range styles {
		if strings.HasPrefix(styleName, "styles/") && styleContent != "" {
			fmt.Println(" "+color.GreenString("☼"), styleName)
			css = append(css, styleContent)
			styles[styleName] = ""
		}
	}

	// Unordered styles
	for styleName, styleContent := range styles {
		if styleContent != "" {
			fmt.Println(" "+color.GreenString("☼"), styleName)
			css = append(css, styleContent)
			styles[styleName] = ""
		}
	}

	// Fonts
	fontsCSS := <-fontsCSSChannel

	return fontsCSS + strings.Join(css, "")
}
