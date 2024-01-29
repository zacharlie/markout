package main

import (
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/zacharlie/markout/internal/lib"
)

func TestFindMarkdownFiles(t *testing.T) {
	tests := []struct {
		root    string
		recurse bool
		want    []string
	}{
		{
			root:    "../tests/data",
			recurse: false,
			want: []string{
				filepath.Join("..", "tests", "data", "test.md"),
				filepath.Join("..", "tests", "data", "file1.md"),
				filepath.Join("..", "tests", "data", "file2.markdown"),
			},
		},
		{
			root:    "../tests/data",
			recurse: true,
			want: []string{
				filepath.Join("..", "tests", "data", "test.md"),
				filepath.Join("..", "tests", "data", "file1.md"),
				filepath.Join("..", "tests", "data", "file2.markdown"),
				filepath.Join("..", "tests", "data", "dir1", "file3.md"),
				filepath.Join("..", "tests", "data", "dir2", "file4.markdown"),
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s_%t", test.root, test.recurse), func(t *testing.T) {
			got, err := lib.FindMarkdownFiles(test.root, test.recurse)
			if err != nil {
				t.Fatalf("FindMarkdownFiles(%q, %t) error: %v", test.root, test.recurse, err)
			}

			sort.Strings(got)
			sort.Strings(test.want)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("FindMarkdownFiles(%q, %t) got %v, want %v", test.root, test.recurse, got, test.want)
			}
		})
	}
}
