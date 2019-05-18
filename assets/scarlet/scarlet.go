package scarletpacker

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/aerogo/flow/jobqueue"
	"github.com/aerogo/pack"
	"github.com/aerogo/scarlet"
	"github.com/akyoto/color"
	"github.com/akyoto/stringutils/unsafe"
)

// ScarletPacker is a packer for scarlet files.
type ScarletPacker struct {
	// Root directory
	root string

	// Where embedded code will be stored
	outputDirectory string

	// The prefix used for terminal output on each file.
	prefix string

	// A list of styles that should be compiled first, the order matters.
	styles []string

	// A list of Google fonts whose definition we are going to download.
	fonts []string

	// A buffered channel that contains our download result.
	fontsChannel chan string
}

// New creates a new ScarletPacker.
func New(root string, styles []string, fonts []string) *ScarletPacker {
	outputDirectory := path.Join(root, "components", "css")
	err := os.MkdirAll(outputDirectory, os.ModePerm)

	if err != nil {
		panic(err)
	}

	packer := &ScarletPacker{
		root:            root,
		outputDirectory: outputDirectory,
		prefix:          color.YellowString(" ★ "),
		styles:          styles,
		fonts:           fonts,
		fontsChannel:    make(chan string, 1),
	}

	err = os.MkdirAll(path.Join(root, "components", "cache", "fonts"), os.ModePerm)

	if err != nil {
		panic(err)
	}

	go func() {
		defer close(packer.fontsChannel)

		if len(packer.fonts) == 0 {
			packer.fontsChannel <- ""
			return
		}

		cached, err := ioutil.ReadFile(path.Join(root, "components", "cache", "fonts", strings.Join(packer.fonts, "|")+".css"))

		if err == nil {
			packer.fontsChannel <- unsafe.BytesToString(cached)
			return
		}

		css, err := downloadFonts(packer.fonts)

		if err != nil {
			panic(err)
		}

		// Save in cache
		err = ioutil.WriteFile(path.Join(root, "components", "cache", "fonts", strings.Join(packer.fonts, "|")+".css"), unsafe.StringToBytes(css), 0777)

		if err != nil {
			panic(err)
		}

		// Send the compressed CSS for the fonts
		packer.fontsChannel <- css
	}()

	return packer
}

// Map maps each job to its processed output.
func (packer *ScarletPacker) Map(job interface{}) interface{} {
	contents, err := ioutil.ReadFile(job.(string))

	if err != nil {
		color.Red(err.Error())
		return nil
	}

	return unsafe.BytesToString(contents)
}

// Reduce combines all outputs.
func (packer *ScarletPacker) Reduce(results jobqueue.Results) {
	buffer := strings.Builder{}

	// Ordered styles
	for _, name := range packer.styles {
		name = "styles/" + name + ".scarlet"
		contents, exists := results[name]

		if !exists {
			color.Yellow("config.json references a file that doesn't exist: %s", name)
			continue
		}

		fmt.Println(packer.prefix, name)
		buffer.WriteString(contents.(string))
		buffer.WriteByte('\n')

		// Remove the referenced style from the unordered results
		// because we already wrote its contents into the buffer.
		delete(results, name)
	}

	// Create a slice which we will sort later
	unorderedStyles := make([]string, 0, len(results))

	// Unordered styles in styles directory
	for name, contents := range results {
		if !strings.HasPrefix(name.(string), "styles/") {
			continue
		}

		fmt.Println(packer.prefix, name)
		unorderedStyles = append(unorderedStyles, contents.(string))

		// Remove the referenced style from the unordered results
		delete(results, name)
	}

	// Unordered styles outside of styles directory
	for name, contents := range results {
		fmt.Println(packer.prefix, name)
		unorderedStyles = append(unorderedStyles, contents.(string))
	}

	// This makes the order deterministic.
	// Since the order is well defined and not random,
	// etag based caching will work.
	sort.Slice(unorderedStyles, func(i, j int) bool {
		a := unorderedStyles[i]
		b := unorderedStyles[j]

		if len(a) == len(b) {
			return pack.HashString(a) < pack.HashString(b)
		}

		return len(a) < len(b)
	})

	// Write all the remaining styles into the buffer
	for _, code := range unorderedStyles {
		buffer.WriteString(code)
		buffer.WriteByte('\n')
	}

	bundledScarlet := buffer.String()
	reader := strings.NewReader(bundledScarlet)
	bundledCSS, err := scarlet.Compile(reader, false)

	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

	// Prepend fonts
	bundledCSS = <-packer.fontsChannel + bundledCSS

	// Write JS bundle into components/css/css.go where it can be used as css.Bundle()
	embedFile := path.Join(packer.outputDirectory, "css.go")
	err = pack.EmbedData(embedFile, "css", "Bundle", bundledCSS)

	if err != nil {
		panic(err)
	}
}
