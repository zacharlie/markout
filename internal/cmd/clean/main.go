//go:build ignore

// # remove generated/ preprocessed files

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	pattern := "*.min.css"
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("error, root path not set")
	}
	fmt.Printf("cleaning artifacts from working directory: %s\n", dir)

	cleanedCount := 0

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		match, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil {
			return err
		}
		if !info.IsDir() && match {
			fmt.Printf("removing file: %s\n", path)
			cleanedCount += 1
			return os.Remove(path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("error removing files: %v", err)
	}

	if cleanedCount == 0 {
		fmt.Printf("no files matching pattern %s removed from %s\n", pattern, dir)
	}
}
