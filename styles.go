package main

import (
	"fmt"
	"os/exec"
	"strings"

	"regexp"

	"github.com/aerogo/aero"
	"github.com/fatih/color"
)

var css []string
var styleCount = 0
var cssChannel = make(chan *StylusCompileResult, 1024)
var cssCommentsRegex = regexp.MustCompile(`\/\*[^*]*\*+([^/*][^*]*\*+)*\/`)

func compileStyle(file string) {
	go func() {
		style, err := exec.Command("stylus", "-p", "--import", "styles/config.styl", "-c", file).Output()

		if err != nil {
			color.Red("Couldn't execute stylus.")
			color.Red(err.Error())
			cssChannel <- &StylusCompileResult{
				file: file,
				css:  "",
			}
			return
		}

		cssChannel <- &StylusCompileResult{
			file: file,
			css:  string(style),
		}
	}()

	styleCount++
}

func getBundledCSS() string {
	// Load config file
	app := aero.New()

	// Fonts
	fontsCSS := getFontsCSS()

	// Wait for stylus to finish compilation
	styles := make(map[string]string)

	for i := 0; i < styleCount; i++ {
		result := <-cssChannel
		styles[result.file] = result.css
	}

	// Ordered styles
	for _, styleName := range app.Config.Styles {
		styleName = "styles/" + styleName + ".styl"
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

	bundledCSS := fontsCSS + strings.Join(css, "")

	// Remove CSS comments
	bundledCSS = cssCommentsRegex.ReplaceAllString(bundledCSS, "")

	// Escape properly
	// bundledCSS = strings.Replace(bundledCSS, "\\", "\\\\", -1)
	// bundledCSS = strings.Replace(bundledCSS, "\"", "\\\"", -1)

	return bundledCSS
}
