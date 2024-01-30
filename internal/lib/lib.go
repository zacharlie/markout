package lib

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/spf13/cobra"
)

//go:generate go run -tags=process_data data/preprocess.go
//go:embed data/*.min.css
var minifiedCssContent embed.FS

func getEmbeddedCssFileContent() (map[string]string, error) {

	var CssStyleFiles = map[string]string{
		"pandoc": "",
		"retro":  "",
	}

	for key := range CssStyleFiles {
		cssContent, err := minifiedCssContent.ReadFile("data/" + key + ".min.css")
		CssStyleFiles[key] = string(cssContent)
		if err != nil {
			return nil, fmt.Errorf("error reading css file %s", key+".min.css")
		}
	}
	return CssStyleFiles, nil
}

func WriteOutput(cmd *cobra.Command, result []byte, inputFile string) error {
	useStdin, _ := cmd.Flags().GetBool("stdin")
	defaultExtension, _ := cmd.Flags().GetString("extension")
	outputDir, _ := cmd.Flags().GetString("outdir")
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	useStdout, _ := cmd.Flags().GetBool("stdout")

	if useStdout {
		// Write to stdout
		fmt.Println(string(result))
		return nil
	} else {
		// Save to a file
		outputFileName := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile)) + defaultExtension

		// check if output directory exists, else create it
		if _, err := os.Stat(outputDir); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(outputDir, os.ModePerm)
			if err != nil {
				log.Printf("error creating output directory %s", outputDir)
			}
		}

		outputPath := filepath.Join(outputDir, outputFileName)
		if !overwrite {
			if _, err := os.Stat(outputPath); err == nil {
				return fmt.Errorf("output file %s already exists, use -w or --overwrite to replace", outputPath)
			}
		}

		err := os.WriteFile(outputPath, []byte(string(result)), 0644)
		if err != nil {
			return fmt.Errorf("error writing output file %s: %v", outputPath, err)
		}

		if useStdin {
			fmt.Printf("Successfully converted stdin to %s\n", outputPath)
		} else {
			fmt.Printf("Successfully converted %s to %s\n", inputFile, outputPath)
		}

		return nil
	}
}

func ProcessContent(content []byte, cssContent []byte, title string, useFullHtml bool) ([]byte, error) {
	var htmlContent strings.Builder
	if useFullHtml {
		htmlContent.WriteString("<!DOCTYPE html>\n<html>\n  <head>\n    <meta charset=\"UTF-8\">\n    <meta name=\"viewport\"")
		htmlContent.WriteString(" content=\"width=device-width, initial-scale=1.0\">\n    <title>")

		htmlContent.WriteString(string(title))

		htmlContent.WriteString("</title>\n")
		htmlContent.WriteString(string(cssContent))
		htmlContent.WriteString("  </head>\n  <body>\n")
	}
	htmlContent.Write(markdown.ToHTML([]byte(content), nil, nil))

	if useFullHtml {
		htmlContent.WriteString("  </body>\n</html>\n")
	}
	return []byte(htmlContent.String()), nil
}

func ReadInput(inputFile string) ([]byte, error) {
	// Read from file
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return []byte(""), fmt.Errorf("error reading input file %s: %v", inputFile, err)
	}

	return content, nil
}

func FindMarkdownFiles(root string, recurse bool) ([]string, error) {
	var matches []string
	patterns := []string{"*.md", "*.markdown"}

	walkFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			for _, pattern := range patterns {
				if matched, err := filepath.Match(pattern, d.Name()); err == nil && matched {
					matches = append(matches, path)
					break // Stop checking other patterns once a match is found
				}
			}
		}

		return nil
	}

	if recurse {
		err := filepath.WalkDir(root, walkFunc)
		if err != nil {
			return nil, err
		}
	} else {
		entries, err := os.ReadDir(root)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			err := walkFunc(filepath.Join(root, entry.Name()), entry, nil)
			if err != nil {
				return nil, err
			}
		}
	}

	return matches, nil
}

func GetCssContent(useStyleTheme string, useStyleFile string, useStyleLink string) ([]byte, error) {
	var cssContent strings.Builder

	minCss, err := getEmbeddedCssFileContent()
	if err != nil {
		fmt.Printf("error reading embedded css files: %v", err)
	}

	if useStyleTheme == "none" &&
		useStyleFile == "none" &&
		useStyleLink == "none" {
		// apply default
		styleContent := minCss["pandoc"]
		cssContent.WriteString("    <style>\n    ")
		cssContent.WriteString(string(styleContent))
		cssContent.WriteString("\n    </style>\n")
	} else if useStyleTheme == "pandoc" {
		styleContent := minCss["pandoc"]
		cssContent.WriteString("    <style>\n    ")
		cssContent.WriteString(string(styleContent))
		cssContent.WriteString("\n    </style>\n")
	} else if useStyleTheme == "retro" {
		styleContent := minCss["retro"]
		cssContent.WriteString("    <style>\n    ")
		cssContent.WriteString(string(styleContent))
		cssContent.WriteString("\n    </style>\n")
	} else if useStyleTheme == "blank" || useStyleTheme == "none" {
		cssContent.WriteString("")
	} else {
		cssContent.WriteString("")
		fmt.Printf("invalid theme selection - no theme data applied.")
	}

	if useStyleFile != "none" && useStyleFile != "" {

		styleFileContent, err := os.ReadFile(useStyleFile)
		if err != nil {
			fmt.Printf("error reading css file %s: %v", useStyleFile, err)
		}

		cssContent.WriteString("    <style>\n")
		cssContent.WriteString(string(styleFileContent))
		cssContent.WriteString("\n    </style>\n")
	}

	if useStyleLink == "bulma" {
		cssContent.WriteString(`    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@0.9.4/css/bulma.min.css" crossorigin="anonymous" referrerpolicy="no-referrer" />`)
		cssContent.WriteString("\n")
	} else if useStyleLink == "bootstrap" {
		cssContent.WriteString(`    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css" crossorigin="anonymous" referrerpolicy="no-referrer" />`)
		cssContent.WriteString("\n")
	} else if useStyleLink == "tachyons" {
		cssContent.WriteString(`    <link rel="stylesheet" href="https://unpkg.com/tachyons@4.12.0/css/tachyons.min.css" crossorigin="anonymous" referrerpolicy="no-referrer" />`)
		cssContent.WriteString("\n")
	} else if useStyleLink == "milligram" {
		cssContent.WriteString(`    <!-- Google Fonts -->`)
		cssContent.WriteString("\n")
		cssContent.WriteString(`    <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Roboto:300,300italic,700,700italic">`)
		cssContent.WriteString("\n")
		cssContent.WriteString(`    <!-- CSS Reset -->`)
		cssContent.WriteString("\n")
		cssContent.WriteString(`    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/normalize/8.0.1/normalize.css">`)
		cssContent.WriteString("\n")
		cssContent.WriteString(`    <!-- Milligram CSS -->`)
		cssContent.WriteString("\n")
		cssContent.WriteString(`    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/milligram/1.4.1/milligram.css">`)
		cssContent.WriteString("\n")
	} else if useStyleLink == "pure" {
		cssContent.WriteString(`    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/purecss@3.0.0/build/pure-min.css" crossorigin="anonymous" referrerpolicy="no-referrer" />`)
		cssContent.WriteString("\n")
	} else if useStyleLink == "wing" {
		cssContent.WriteString(`    <link rel="stylesheet" href="https://unpkg.com/wingcss" crossorigin="anonymous" referrerpolicy="no-referrer" />`)
		cssContent.WriteString("\n")
	} else if useStyleLink == "pico" {
		cssContent.WriteString(`    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@1/css/pico.min.css" crossorigin="anonymous" referrerpolicy="no-referrer" />`)
		cssContent.WriteString("\n")
	} else if useStyleLink != "none" && useStyleLink != "" {
		cssContent.WriteString(`    <link rel="stylesheet" href="`)
		cssContent.WriteString(useStyleLink)
		cssContent.WriteString(`" crossorigin="anonymous" referrerpolicy="no-referrer" />`)
		cssContent.WriteString("\n")
	}

	return []byte(cssContent.String()), nil
}
