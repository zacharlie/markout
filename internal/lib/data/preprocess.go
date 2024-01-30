//go:build ignore

package main

import (
	"bytes"
	"embed"
	"fmt"
	"os"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
)

var CssFileNames = []string{
	"pandoc",
	"retro",
}

func main() {
	for _, fileName := range CssFileNames {
		minifyCss(fileName)
	}
}

//go:embed *.css
var cssFileContent embed.FS

func minifyCss(fileName string) {

	// Input file path
	inputFile := fileName + ".css"

	// Output file path
	outputFile := "data/" + fileName + ".min.css"

	// Open the input file for reading
	input, err := cssFileContent.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	inputData := bytes.NewReader(input)

	// Create or open the output file for writing
	output, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	err = m.Minify("text/css", output, inputData)
	if err != nil {
		panic(err)
	}

	fmt.Printf("css file %s minified", outputFile)

}
