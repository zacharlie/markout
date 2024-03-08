package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zacharlie/markout/internal/lib"
)

var rootCmd = &cobra.Command{
	Use:   "markout",
	Short: "Convert Markdown files to HTML",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Calculate and set the default value of the flag here
		preFlightChecks(cmd, args)
	},
	Run: convertMarkdown,
}

var (
	outputDir        string
	overwriteOutput  bool
	defaultExtension string
	useStdin         bool
	useStdout        bool
	runRecursive     bool
	useFullHtml      bool
	useStyleTheme    string
	useStyleFile     string
	useStyleLink     string
	sanitizeOutput   bool
	printToPdf       bool
)

func preFlightChecks(cmd *cobra.Command, args []string) {
	validateFlags(cmd, args)
	setDefaultFlagValues(cmd)
}

func validateFlags(cmd *cobra.Command, args []string) {
	useStdout, _ := cmd.Flags().GetBool("stdout")
	useStdin, _ := cmd.Flags().GetBool("stdin")
	outputDir, _ := cmd.Flags().GetBool("directory")
	runRecursive, _ := cmd.Flags().GetBool("recurse")
	if useStdout && outputDir {
		log.Fatalf("Cannot use stdout and directory flags together.")
	}
	if len(args) > 0 && useStdin {
		log.Fatalf("Cannot use stdin and file inputs simultaneously.")
	}
	if useStdin && runRecursive {
		log.Fatalf("Cannot use recursive with stdin.")
	}
}

func setDefaultFlagValues(cmd *cobra.Command) {
	stdout, _ := cmd.Flags().GetBool("stdout")
	if !stdout && !cmd.Flags().Changed("full") {
		cmd.Flags().Set("full", "true")
	}
	if !cmd.Flags().Changed("theme") {
		cmd.Flags().Set("theme", "undefined")
	}
	if !cmd.Flags().Changed("style") {
		cmd.Flags().Set("style", "undefined")
	}
	if !cmd.Flags().Changed("link") {
		cmd.Flags().Set("link", "undefined")
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&useStdin, "stdin", "i", false, `Read input from stdin`)
	rootCmd.Flags().BoolVarP(&useStdout, "stdout", "o", false, `Print output to stdout`)
	rootCmd.Flags().BoolVarP(&printToPdf, "pdf", "p", false, `Save a pdf copy of the output (applies when writing files only. Requires chrome executable.)`)
	rootCmd.Flags().BoolVarP(&overwriteOutput, "overwrite", "w", false, `Overwrite existing output files`)
	rootCmd.Flags().BoolVarP(&runRecursive, "recurse", "r", false, `Run recursively on subdirectory contents`)
	rootCmd.Flags().BoolVarP(&sanitizeOutput, "sanitize", "c", false, `Sanitize the HTML output with bluemonday`)
	rootCmd.Flags().StringVarP(&defaultExtension, "extension", "e", ".html", `Output file extension`)
	rootCmd.Flags().StringVarP(&outputDir, "directory", "d", "./markoutput", `Output directory`)
	rootCmd.Flags().BoolVarP(&useFullHtml, "full", "f", false, `Write complete HTML page (including head, with md content in body)`)
	rootCmd.Flags().StringVarP(&useStyleFile, "style", "s", "none", `Path to a css file. Contents are injected into a <style> block`)
	rootCmd.Flags().StringVarP(&useStyleLink, "link", "l", "none", `Text value to insert into the href attribute of <link rel="stylesheet" />`)
	rootCmd.Flags().StringVarP(&useStyleTheme, "theme", "t", "light", `A predefined css to embed. Options include "none", "light", and "dark"`)
	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false
}

func convertMarkdown(cmd *cobra.Command, args []string) {
	cssContent, err := lib.GetCssContent(
		strings.ToLower(useStyleTheme),
		strings.ToLower(useStyleFile),
		strings.ToLower(useStyleLink),
	)
	if err != nil {
		log.Fatalf("error getting css content: %v", err)
	}

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

		result, err := lib.ProcessContent(content, cssContent, "MarkOut", useFullHtml)
		if err != nil {
			log.Fatalf("error processing stdin: %v", err)
		}

		err = lib.WriteOutput(cmd, result, filepath.Join(outputDir, "MarkOut"+defaultExtension))
		if err != nil {
			log.Fatalf("error writing output: %v", err)
		}
	} else {
		for _, inputFile := range args {
			content, err := lib.ReadInput(inputFile)
			if err != nil {
				log.Fatalf("error reading from file: %v", err)
			}

			result, err := lib.ProcessContent(content, cssContent,
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
