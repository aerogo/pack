package main

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/aerogo/flow/jobqueue"
	"github.com/aerogo/scarlet"
	"github.com/blitzprog/color"
)

var cssCommentsRegex = regexp.MustCompile(`\/\*[^*]*\*+([^/*][^*]*\*+)*\/`)
var styleAnnouncePrefix = " " + color.YellowString("★") + " "

func scarletWork(job interface{}) interface{} {
	file := job.(string)
	scarletCode, _ := ReadFile(file)
	return scarletCode
}

func scarletFinish(results jobqueue.Results) {
	// Convert to map[string]string
	styles := ToStringMap(results)

	// Bundled CSS
	bundledCSS := getBundledCSS(styles)

	// Write CSS bundle into $.css.go where it can be referenced as components.CSS
	EmbedData(path.Join(outputFolder, "css", "css.go"), "css", "Bundle", bundledCSS)
}

func getBundledCSS(styles map[string]string) string {
	scarletCodes := []string{}

	// Ordered styles
	for _, styleName := range config.Styles {
		styleName = "styles/" + styleName + ".scarlet"
		styleContent, exists := styles[styleName]

		if exists {
			announceStyle(styleName)
			scarletCodes = append(scarletCodes, styleContent)
			delete(styles, styleName)
		}
	}

	// Create a slice which we will sort later
	unorderedStyles := []string{}

	// Unordered styles in styles directory
	for styleName, styleContent := range styles {
		if strings.HasPrefix(styleName, "styles/") {
			announceStyle(styleName)
			unorderedStyles = append(unorderedStyles, styleContent)
			delete(styles, styleName)
		}
	}

	// Unordered styles
	for styleName, styleContent := range styles {
		announceStyle(styleName)
		unorderedStyles = append(unorderedStyles, styleContent)
	}

	// This doesn't really have any meaning besides making the order deterministic.
	// Since the order is well defined and not random, hash based caching will work.
	sort.Slice(unorderedStyles, func(i, j int) bool {
		a := unorderedStyles[i]
		b := unorderedStyles[j]

		if len(a) == len(b) {
			return HashString(a) < HashString(b)
		}

		return len(a) < len(b)
	})

	scarletCodes = append(scarletCodes, unorderedStyles...)

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
