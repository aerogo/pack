package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aerogo/aero"
	"github.com/aerogo/pixy"
	"github.com/fatih/color"
)

const (
	pixyExtension   = ".pixy"
	stylExtension   = ".styl"
	outputFolder    = "components"
	outputExtension = ".go"
)

// StylusCompileResult ...
type StylusCompileResult struct {
	file string
	css  string
}

func main() {
	// Load config file
	app := aero.New()

	pixy.PackageName = outputFolder

	var css []string

	styleCount := 0
	cssChannel := make(chan *StylusCompileResult, 1024)

	os.Mkdir(outputFolder, 0644)

	filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}

		switch filepath.Ext(path) {
		// Pixy
		case pixyExtension:
			fmt.Println(" "+color.GreenString("❀"), path)
			pixy.CompileFileAndSaveIn(path, outputFolder)

		// Stylus
		case stylExtension:
			go func() {
				style, err := exec.Command("stylus", "-p", "--import", "styles/config.styl", "-c", path).Output()

				if err != nil {
					color.Red("Couldn't execute stylus.")
					color.Red(err.Error())
					cssChannel <- &StylusCompileResult{
						file: path,
						css:  "",
					}
					return
				}

				cssChannel <- &StylusCompileResult{
					file: path,
					css:  string(style),
				}
			}()

			styleCount++
		}

		return nil
	})

	// Fonts
	fontsCSS := getFontsCSS()

	// CSS
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
	bundledCSS = strings.Replace(bundledCSS, "\\", "\\\\", -1)
	bundledCSS = strings.Replace(bundledCSS, "\"", "\\\"", -1)

	fmt.Println()
	fmt.Println("Done.")
}
