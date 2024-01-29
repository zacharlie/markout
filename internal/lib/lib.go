package lib

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/spf13/cobra"
)

func ConvertMarkdown(cmd *cobra.Command, args []string) {
	useStdin, _ := cmd.Flags().GetBool("stdin")
	outputDir, _ := cmd.Flags().GetString("outdir")

	if len(args) == 0 && !useStdin {
		log.Fatalf("Please provide at least one input file or use stdin flag (-s)")
	}

	if useStdin {
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("error reading from stdin: %v", err)
		}

		result, err := processContent(content, "MarkOut")
		if err != nil {
			fmt.Printf("Error processing stdin: %v\n", err)
		}

		err = writeOutput(cmd, result, filepath.Join(outputDir, "MarkOut.html"))
		if err != nil {
			log.Fatalf("error writing output: %v", err)
		}
	} else {
		for _, inputFile := range args {
			content, err := readInput(inputFile)
			if err != nil {
				log.Fatalf("error reading from file: %v", err)
			}

			result, err := processContent(content,
				strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile)))
			if err != nil {
				log.Fatalf("Error processing file %s: %v\n", inputFile, err)
			}
			err = writeOutput(cmd, result, inputFile)
			if err != nil {
				log.Fatalf("error writing output: %v", err)
			}
		}
	}
}

func writeOutput(cmd *cobra.Command, result []byte, inputFile string) error {
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

func processContent(content []byte, title string) ([]byte, error) {
	var htmlContent strings.Builder
	htmlContent.WriteString("<!DOCTYPE html>\n<html>\n  <head>\n    <meta charset=\"UTF-8\">\n    <meta name=\"viewport\"")
	htmlContent.WriteString(" content=\"width=device-width, initial-scale=1.0\">\n    <title>")

	htmlContent.WriteString(string(title))

	htmlContent.WriteString("</title>\n  </head>\n  <body>\n")

	htmlContent.Write(markdown.ToHTML([]byte(content), nil, nil))

	htmlContent.WriteString("  </body>\n</html>\n")
	return []byte(htmlContent.String()), nil
}

func readInput(inputFile string) ([]byte, error) {
	// Read from file
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return []byte(""), fmt.Errorf("error reading input file %s: %v", inputFile, err)
	}

	return content, nil
}
