/*
Copyright © 2024 Chris Greaves cjgreaves97@hotmail.co.uk

See the file COPYING in the root of this repository for details.
*/
package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Chris-Greaves/pencil/processor"
	"github.com/spf13/cobra"
)

var (
	files           []string
	variables       []string
	parsedVariables = make(map[string]string)

	ErrInvalidVariable           = errors.New("variable was invalid, please follow 'NAME=value'")
	ErrNoFilesOrFoldersSpecified = errors.New("must specify at least one file or folder")
)

var rootCmd = &cobra.Command{
	Use:   "pencil",
	Short: "A Tool to fill in the blanks in your files.",
	Long: `This tool uses a series of flags and the Go Templating library to
fill in secrets and values in files and paths.

Point it to a file and it'll treat the file like a Go Template and execute it, 
pulling in Environment Variables and any other variables passed into the tool.
For example: "pencil -v SECRET_KEY=something-secret app/config.yml"`,

	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return ErrNoFilesOrFoldersSpecified
		}

		// Check Args
		if validatedPaths, err := validatePathArgs(args); err != nil {
			return err
		} else {
			files = validatedPaths
		}

		// Check Variables
		if validatedVars, err := validateDirectVars(variables); err != nil {
			return err
		} else {
			parsedVariables = validatedVars
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("Files: %s", files)
		log.Printf("Variables: %s", parsedVariables)

		model := processor.BuildModel(parsedVariables)
		proc := processor.NewWithModel(model)
		for _, file := range files {
			var buf bytes.Buffer
			if err := proc.ParseAndExecuteFile(file, &buf); err != nil {
				return fmt.Errorf("error while processing %v: %w", file, err)
			}

			if err := writeToFile(file, buf); err != nil {
				return fmt.Errorf("error while writing over existing file %v, may be left in a partial state: %w", file, err)
			}
		}

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringArrayVarP(&variables, "variable", "v", make([]string, 0), "A variable to be used in the Templates. e.g. SECRET_KEY=something-secret")
}

func validatePathArgs(args []string) ([]string, error) {
	returnPaths := make([]string, 0)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		fileInfo, err := os.Stat(arg)
		if err != nil {
			return nil, fmt.Errorf("error getting file or folder %v: %w", arg, err)
		}
		if fileInfo.IsDir() {
			err = filepath.WalkDir(arg, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return fmt.Errorf("error accessing path %v: %w", path, err)
				}

				if !d.IsDir() {
					returnPaths = append(returnPaths, path)
				}

				return nil
			})

			if err != nil {
				return nil, fmt.Errorf("error while walking through directory: %w", err)
			}

			continue
		}

		returnPaths = append(returnPaths, arg)
	}
	return returnPaths, nil
}

func validateDirectVars(vars []string) (map[string]string, error) {
	returnVars := make(map[string]string)
	for i := 0; i < len(vars); i++ {
		variable := vars[i]
		index := strings.Index(variable, "=")
		if index == -1 {
			return nil, ErrInvalidVariable
		}
		key := variable[:index]
		value := variable[index+1:]
		returnVars[strings.TrimSpace(key)] = value
	}

	return returnVars, nil
}

func writeToFile(filePath string, buf bytes.Buffer) error {
	inputPath := filePath
	if filepath.Ext(filePath) == ".gotmpl" {
		inputPath = strings.TrimSuffix(filePath, ".gotmpl")
	}
	return os.WriteFile(inputPath, buf.Bytes(), 0666)
}
