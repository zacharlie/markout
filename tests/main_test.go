// main_test.go

package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestFileInputOutput(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Create a temporary markdown file
	markdownContent := "# Test\nThis is a test."
	markdownPath := tempDir + "/test.md"
	err := os.WriteFile(markdownPath, []byte(markdownContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Run the markout command with file input/output
	wd, err := os.Getwd() // current working directory
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "run", wd+"/../cmd/markout/main.go", "--outdir", tempDir, markdownPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Error running markout command: %v\nOutput:\n%s", err, output)
	}

	// Check if the output HTML file was created
	htmlPath := tempDir + "/test.html"
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		t.Fatalf("Expected output file %s not found", htmlPath)
	}

	// Cleanup: remove temporary files
	os.Remove(markdownPath)
	os.Remove(htmlPath)
}

// func TestStdinStdout(t *testing.T) {
// 	// Run markout with stdin-out
// 	wd, err := os.Getwd() // current working directory
// 	cmd := exec.Command("go", "run", wd+"/../cmd/markout/main.go", "--stdin", "--stdout")
// 	cmd.Stdin = bytes.NewBufferString("# Test\nThis is a test.")
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		t.Fatalf("Error running markout command: %v\nOutput:\n%s", err, output)
// 	}

// 	// Check if the output contains the expected HTML content
// 	expectedOutput := "<!DOCTYPE html>\n<html>\n  <head>\n    <meta charset=\"UTF-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n    <title>MarkOut</title>\n  </head>\n  <body>\n<h1>Test</h1>\n\n<p>This is a test.</p>\n  </body>\n</html>\n"
// 	if string(output) != expectedOutput {
// 		t.Fatalf("Expected output:\n%s\nActual output:\n%s", expectedOutput, output)
// 	}
// }
