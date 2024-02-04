package cmd

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Chris-Greaves/pencil/test"
)

func TestValidateVars(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  error
	}{
		{"no errors when valid var", "foo=bar", nil},
		{"no error when value contains =", "foo=bar==", nil},
		{"error when input doesn't contain =", "foo", ErrInvalidVariable},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vars := []string{
				tt.input,
			}

			output, err := validateDirectVars(vars)
			t.Logf("The output was %v", output)

			if err != tt.want {
				t.Errorf("expected error to be '%s', but got '%s'", tt.want, err)
			}
		})
	}
}

func TestValidatePathArgs(t *testing.T) {
	testCases := []struct {
		desc       string
		input      []string
		output_err error
		output     []string
	}{
		{
			desc:       "can accept valid file path",
			input:      []string{"../test/config.yml"},
			output_err: nil,
			output:     []string{"../test/config.yml"},
		},
		{
			desc:       "can accept valid folder path",
			input:      []string{"../test/folder/"},
			output_err: nil,
			output:     []string{"../test/folder/file1.ini", "../test/folder/file2.txt"},
		},
		{
			desc:       "will error when file doesn't exist",
			input:      []string{"../test/does_not_exist.yaml"},
			output_err: fs.ErrNotExist,
			output:     []string{},
		},
		{
			desc:       "will error when folder doesn't exist",
			input:      []string{"../test/does_not_exist/"},
			output_err: fs.ErrNotExist,
			output:     []string{},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			output, err := validatePathArgs(tC.input)
			if (err != nil && tC.output_err != nil && !errors.Is(err, tC.output_err)) || (err != nil && tC.output_err == nil) { // if error wasn't the expected error, or error wasn't expected.
				t.Errorf("expected err to be '%s', but got '%s'", tC.output_err, err)
			}
			if len(output) != 0 || len(tC.output) != 0 {
				if eq := reflect.DeepEqual(output, tC.output); !eq {
					t.Errorf("expected output to be '%s', but got '%s'", tC.output, output)
				}
			}
		})
	}
}

func TestWriteToFile(t *testing.T) {
	testCases := []struct {
		desc           string
		path           string
		input_filename string
		created_file   string
		err_expected   bool
		err_contains   string
	}{
		{
			desc:           "will write to the file",
			path:           t.TempDir(),
			input_filename: "file.txt",
			created_file:   "file.txt",
			err_expected:   false,
			err_contains:   "",
		},
		{
			desc:           "error when path doesn't exist",
			path:           "/does/not/exist",
			input_filename: "file.txt",
			created_file:   "file.txt",
			err_expected:   true,
			err_contains:   test.PathNotFoundText(),
		},
		{
			desc:           "will remove .gotmpl from file",
			path:           t.TempDir(),
			input_filename: "file.txt.gotmpl",
			created_file:   "file.txt",
			err_expected:   false,
			err_contains:   "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			fileContent := bytes.NewBufferString("file contents")
			fullInputPath := filepath.Join(tC.path, tC.input_filename)
			fullOutputPath := filepath.Join(tC.path, tC.created_file)

			err := writeToFile(fullInputPath, *fileContent)
			if (err != nil) != tC.err_expected {
				t.Errorf("resulting error was not what was expected, '%s'", err)
			}
			if err != nil && !strings.Contains(err.Error(), tC.err_contains) {
				t.Errorf("expected error to contain '%s', but got '%s'", tC.err_contains, err.Error())
			}
			if err == nil && !tC.err_expected {
				_, statErr := os.Stat(fullOutputPath)
				if statErr != nil {
					t.Error("expected file wasn't created")
				}
			}
		})
	}
}
