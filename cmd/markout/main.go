package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	lib "github.com/zacharlie/markout/internal/lib"
)

var rootCmd = &cobra.Command{
	Use:   "markout",
	Short: "Convert Markdown files to HTML",
	Run:   convertMarkdown,
}

var (
	outputDir        string
	overwriteOutput  bool
	defaultExtension string = ".html"
	useStdin         bool
	useStdout        bool
	runRecursive     bool
	useFullHtml      bool
)

func init() {
	rootCmd.Flags().StringVarP(&outputDir, "outdir", "d", "./markoutput", "Output directory")
	rootCmd.Flags().StringVarP(&defaultExtension, "extension", "e", ".html", "Output file extension")
	rootCmd.Flags().BoolVarP(&overwriteOutput, "overwrite", "w", false, "Overwrite existing output files")
	rootCmd.Flags().BoolVarP(&useStdin, "stdin", "s", false, "Read input from stdin")
	rootCmd.Flags().BoolVarP(&useStdout, "stdout", "o", false, "Print output to stdout")
	rootCmd.Flags().BoolVarP(&runRecursive, "recurse", "r", false, "Run recursively on subdirectory contents")
	rootCmd.Flags().BoolVarP(&useFullHtml, "full", "f", false, "Write complete HTML page (including head, with md content in body)")
}

func convertMarkdown(cmd *cobra.Command, args []string) {
	if len(args) == 0 && !useStdin {
		// get current working directory
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("error getting current working directory: %v", err)
		}

		mdFiles, err := lib.FindMarkdownFiles(cwd, runRecursive)
		if err != nil {
			log.Fatalf("error getting files: %v", err)
		} else if len(mdFiles) == 0 {
			log.Fatalf("no markdown files found")
		}

		args = append(args, mdFiles...)
	}

	if useStdin {
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("error reading from stdin: %v", err)
		}

		result, err := lib.ProcessContent(content, "MarkOut", useFullHtml)
		if err != nil {
			log.Fatalf("error processing stdin: %v", err)
		}

		err = lib.WriteOutput(cmd, result, filepath.Join(outputDir, "MarkOut.html"))
		if err != nil {
			log.Fatalf("error writing output: %v", err)
		}
	} else {
		for _, inputFile := range args {
			content, err := lib.ReadInput(inputFile)
			if err != nil {
				log.Fatalf("error reading from file: %v", err)
			}

			result, err := lib.ProcessContent(content,
				strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile)),
				useFullHtml)
			if err != nil {
				log.Fatalf("Error processing file %s: %v\n", inputFile, err)
			}
			err = lib.WriteOutput(cmd, result, inputFile)
			if err != nil {
				log.Fatalf("error writing output: %v", err)
			}
		}
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
}
