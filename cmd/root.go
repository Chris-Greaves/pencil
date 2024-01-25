/*
Copyright © 2024 Chris Greaves cjgreaves97@hotmail.co.uk
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	files           []string
	folders         []string
	variables       []string
	parsedVariables = make(map[string]string)

	ErrInvalidVariable          = errors.New("variable was invalid, please follow 'NAME=value'")
	ErrFolderGivenAsFile        = errors.New("folder given as file, please us -F for folders")
	ErrFileGivenAsFolder        = errors.New("file given as folder, please us -f for files")
	ErrNoFilesOrFolderSpecified = errors.New("must specify at least one file or folder")
)

var rootCmd = &cobra.Command{
	Use:   "pencil",
	Short: "A Tool to fill in the blanks.",
	Long: `This tool uses a series of flags and the Go Templating library to
fill in secrets and values in files and paths.

Point it to a file and it'll treat the file like a Go Template and execute it, 
pulling in Environment Variables and any other secrets passed into the tool.
For example: "pencil -f app/config.yml -v SECRET_KEY=something-secret"`,

	Args: func(cmd *cobra.Command, args []string) error {
		if len(files)+len(folders) == 0 {
			return ErrNoFilesOrFolderSpecified
		}

		// Check Files
		for i := 0; i < len(files); i++ {
			file := files[i]
			fileInfo, err := os.Stat(file)
			if err != nil {
				return errors.Join(fmt.Errorf("error getting file at %v", file), err)
			}
			if fileInfo.IsDir() {
				return ErrFolderGivenAsFile
			}
		}

		// Check Folders
		for i := 0; i < len(folders); i++ {
			folder := folders[i]
			fileInfo, err := os.Stat(folder)
			if err != nil {
				return errors.Join(fmt.Errorf("error getting folder at %v", folder), err)
			}
			if !fileInfo.IsDir() {
				return ErrFileGivenAsFolder
			}
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

	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Files: %s", files)
		log.Printf("Folders: %s", folders)
		log.Printf("Variables: %s", parsedVariables)

	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringArrayVarP(&files, "file", "f", make([]string, 0), "Path to a file you want run through the Template Engine")
	rootCmd.Flags().StringArrayVarP(&folders, "folder", "F", make([]string, 0), "Path to a directory you want all files underneath run through the Template Engine")
	rootCmd.Flags().StringArrayVarP(&variables, "variable", "v", make([]string, 0), "A variable to be used in the Templates. e.g. SECRET_KEY=something-secret")
}
