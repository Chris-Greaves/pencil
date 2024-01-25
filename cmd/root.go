/*
Copyright © 2024 Chris Greaves cjgreaves97@hotmail.co.uk
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
		files = make([]string, 0)
		for i := 0; i < len(args); i++ {
			arg := args[i]
			fileInfo, err := os.Stat(arg)
			if err != nil {
				return fmt.Errorf("error getting file or folder %v: %w", arg, err)
			}
			if fileInfo.IsDir() {
				err = filepath.WalkDir(arg, func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						log.Printf("Error accessing path %v, skipping", path)
						return nil
					}

					if !d.IsDir() {
						files = append(files, path)
					}

					return nil
				})

				if err != nil {
					return fmt.Errorf("error while walking through directory: %w", err)
				}

				continue
			}

			files = append(files, arg)
		}

		// Check Variables
		for i := 0; i < len(variables); i++ {
			variable := variables[i]
			splitVar := strings.Split(variable, "=")
			if len(splitVar) != 2 {
				return ErrInvalidVariable
			}
			parsedVariables[splitVar[0]] = splitVar[1]
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("Files: %s", files)
		log.Printf("Variables: %s", parsedVariables)

		proc := processor.NewWithModel(processor.Model{Var: parsedVariables})
		for _, file := range files {
			var buf bytes.Buffer
			if err := proc.ParseAndExecuteFile(file, &buf); err != nil {
				return fmt.Errorf("error while processing %v: %w", file, err)
			}

			if err := os.WriteFile(file, buf.Bytes(), 0666); err != nil {
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
