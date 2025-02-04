package recurparse

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestGetFiles(t *testing.T) {
	type getFilesTestdatum struct {
		Directory   string
		Glob        string
		Expected    []string
		ExpectedErr string
	}

	getFilesTestdata := []getFilesTestdatum{
		{
			Directory: "testdata/getFiles/1_simple_flat",
			Glob:      "*.html",
			Expected: []string{
				"1.html",
				"3.html",
			},
		},
		{
			Directory: "testdata/getFiles/2_simple_dirs",
			Glob:      "*.html",
			Expected: []string{
				"1.html",
				"3.html",
				"first/4.html",
				"second/7.html",
			},
		},
		{
			Directory: "testdata/getFiles/3_symlink",
			Glob:      "*.html",
			Expected: []string{
				"1.html",
				"3.html",
				"first/4.html",
				"second/4.html",
			},
		},
	}

TEST:
	for i, d := range getFilesTestdata {
		resolved, err := filepath.EvalSymlinks(d.Directory)
		if err != nil {
			t.Fatalf("test %d: cannot resolve %q: %+v", i, d.Directory, err)
		}

		fsys := os.DirFS(resolved)

		files, err := getFilesFS(fsys, d.Glob)

		if d.ExpectedErr != "" {
			if err == nil || err.Error() != d.ExpectedErr {
				t.Fatalf("Expected error %q, got %+v", d.ExpectedErr, err)
			}
			continue TEST
		}

		if err != nil {
			t.Fatalf("test %d: err %+v", i, err)
		}

		if len(files) != len(d.Expected) {
			for _, f := range files {
				fmt.Println(f)
			}
			t.Fatalf("test %d: different lengths: %d vs %d", i, len(files), len(d.Expected))
		}

		for j := range files {
			if files[j] != d.Expected[j] {
				t.Fatalf("test %d: %d : %q != %q", i, j, files[j], d.Expected[j])
			}
		}
	}
}
