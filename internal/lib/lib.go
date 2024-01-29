package lib

import (
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

func ProcessContent(content []byte, title string, useRawHtml bool) ([]byte, error) {
	var htmlContent strings.Builder
	if !useRawHtml {
		htmlContent.WriteString("<!DOCTYPE html>\n<html>\n  <head>\n    <meta charset=\"UTF-8\">\n    <meta name=\"viewport\"")
		htmlContent.WriteString(" content=\"width=device-width, initial-scale=1.0\">\n    <title>")

		htmlContent.WriteString(string(title))

		htmlContent.WriteString("</title>\n  </head>\n  <body>\n")
	}
	htmlContent.Write(markdown.ToHTML([]byte(content), nil, nil))

	if !useRawHtml {
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
