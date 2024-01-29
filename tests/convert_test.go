package main

import (
	"os"
	"testing"

	"github.com/gomarkdown/markdown"
)

// Test against existing local files with default MD->HTML conversion
func TestConvertMarkdown(t *testing.T) {
	tests := []struct {
		inputFile  string
		outputFile string
		expected   string
	}{
		{
			inputFile:  "../tests/test.md",
			outputFile: "../tests/test.html",
			expected:   "<h1>Test</h1>\n\n<p>Test Content</p>\n",
		},
	}

	for _, test := range tests {
		t.Run(test.inputFile, func(t *testing.T) {
			// Read the input file.
			input, err := os.ReadFile(test.inputFile)
			if err != nil {
				t.Fatalf("error reading input file: %v", err)
			}

			// Convert the Markdown to HTML.
			html := markdown.ToHTML(input, nil, nil)
			if err != nil {
				t.Fatalf("error converting Markdown to HTML: %v", err)
			}

			// Write the HTML to the output file.
			err = os.WriteFile(test.outputFile, html, 0644)
			if err != nil {
				t.Fatalf("error writing output file: %v", err)
			}

			// Read the output file.
			output, err := os.ReadFile(test.outputFile)
			if err != nil {
				t.Fatalf("error reading output file: %v", err)
			}

			// Compare the output to the expected output.
			if string(output) != test.expected {
				t.Errorf("expected output to be %q, but got %q", test.expected, string(output))
			}
		})
	}
}
