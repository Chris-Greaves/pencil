package processor

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestConvertVarsArrayToMap(t *testing.T) {
	testCases := []struct {
		desc   string
		input  []string
		output map[string]string
	}{
		{
			desc:   "will convert valid vars",
			input:  []string{"foo=bar"},
			output: map[string]string{"foo": "bar"},
		},
		{
			desc:   "will convert with more than one =",
			input:  []string{"foo=bar=="},
			output: map[string]string{"foo": "bar=="},
		},
		{
			desc:   "will ignore vars without an =",
			input:  []string{"foo=bar", "boo"},
			output: map[string]string{"foo": "bar"},
		},
	}
	for _, c := range testCases {
		t.Run(c.desc, func(t *testing.T) {
			result := convertVarsArrayToMap(c.input)

			if eq := reflect.DeepEqual(result, c.output); !eq {
				t.Errorf("got '%v', expected '%v'", result, c.output)
			}
		})
	}
}

func TestParseAndExecuteFile(t *testing.T) {
	testCases := []struct {
		desc         string
		path         string
		vars         map[string]string
		envs         map[string]string
		output       string
		expect_err   bool
		err_contains string
	}{
		{
			desc:       "can handle valid file",
			path:       copyToTempForTesting("../test/config.yml", t),
			vars:       map[string]string{"Name": "Chris"},
			envs:       map[string]string{"AppName": "Pencil"},
			output:     "Hello Chris, it's Pencil",
			expect_err: false,
		},
		{
			desc:         "error when template is invalid",
			path:         copyToTempForTesting("../test/invalid.yml", t),
			vars:         map[string]string{"Name": "Chris"},
			envs:         map[string]string{"AppName": "Pencil"},
			output:       "",
			expect_err:   true,
			err_contains: "unexpected",
		},
		{
			desc:         "error when file doesn't exist",
			path:         "../test/does_not_exist.yml",
			vars:         map[string]string{"Name": "Chris"},
			envs:         map[string]string{"AppName": "Pencil"},
			output:       "",
			expect_err:   true,
			err_contains: "The system cannot find the file specified",
		},
		{
			desc:         "error when a function that doesn't exist is used",
			path:         copyToTempForTesting("../test/bad_function.yml", t),
			vars:         map[string]string{"Name": "Chris"},
			envs:         map[string]string{"AppName": "Pencil"},
			output:       "",
			expect_err:   true,
			err_contains: "function \"bang\" not defined",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			model := Model{Var: tC.vars, Env: tC.envs}
			proc := NewWithModel(model)
			var buf strings.Builder

			err := proc.ParseAndExecuteFile(tC.path, &buf)
			if err != nil && !tC.expect_err {
				t.Errorf("error wasn't expected but got '%s'", err)
			} else if err == nil && tC.expect_err {
				t.Errorf("error was expected, but got '%s'", err)
			} else if err != nil && !strings.Contains(err.Error(), tC.err_contains) {
				t.Errorf("error was expected to contain '%v', but got '%s'", tC.err_contains, err)
			}

			output := buf.String()
			if len(output) != 0 || len(tC.output) != 0 {
				if output != tC.output {
					t.Errorf("expected output to be '%s', but got '%s'", tC.output, output)
				}
			}
		})
	}
}

func copyToTempForTesting(src string, t *testing.T) string {
	tempDir := t.TempDir()

	fileName := filepath.Base(src)
	dst := filepath.Join(tempDir, fileName)

	sourceFile, err := os.Open(src)
	if err != nil {
		t.Fatal(err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		t.Fatal(err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		t.Fatal(err)
	}

	return dst
}
