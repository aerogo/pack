package scarletpacker

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/OneOfOne/xxhash"
	"github.com/aerogo/flow/jobqueue"
	"github.com/aerogo/scarlet"
	"github.com/akyoto/color"
	"github.com/akyoto/stringutils/unsafe"
)

// ScarletPacker is a packer for scarlet files.
type ScarletPacker struct {
	// Root directory
	root string

	// The prefix used for terminal output on each file.
	prefix string

	// A list of styles that should be compiled first, the order matters.
	styles []string
}

// New creates a new ScarletPacker.
func New(root string, styles []string) *ScarletPacker {
	return &ScarletPacker{
		root:   root,
		prefix: color.YellowString(" ★ "),
		styles: styles,
	}
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
	bundledCSS := strings.Builder{}

	// Ordered styles
	for _, name := range packer.styles {
		name = "styles/" + name + ".scarlet"
		contents, exists := results[name]

		if !exists {
			color.Yellow("config.json references a file that doesn't exist: %s", name)
			continue
		}

		fmt.Println(packer.prefix, name)
		bundledCSS.WriteString(contents.(string))

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

	// Unordered styles
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
			return hashString(a) < hashString(b)
		}

		return len(a) < len(b)
	})

	// Write all the remaining styles into the buffer
	for _, code := range unorderedStyles {
		bundledCSS.WriteString(code)
	}

	allCSSConcatenated := bundledCSS.String()
	reader := strings.NewReader(allCSSConcatenated)
	css, err := scarlet.Compile(reader, false)

	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

	fmt.Println(css)
}

// hashString hashes a long string.
func hashString(data string) uint64 {
	h := xxhash.NewS64(0)
	_, _ = h.WriteString(data)
	return h.Sum64()
}
