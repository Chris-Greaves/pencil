/*
Copyright © 2024 Chris Greaves cjgreaves97@hotmail.co.uk
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pencil",
	Short: "A Tool to fill in the blanks.",
	Long: `This tool uses a series of flags and the Go Templating library to
fill in secrets and values in files and paths.

Point it to a file and it'll treat the file like a Go Template and execute it, 
pulling in Environment Variables and any other secrets passed into the tool.
For example: "pencil -f app/config.yml -v SECRET_KEY=something-secret"`,

	Run: func(cmd *cobra.Command, args []string) {
		fileFlag := cmd.Flag("file").Value
		log.Printf("Files: %s", fileFlag)
		folderFlag := cmd.Flag("folder").Value
		log.Printf("Folders: %s", folderFlag)
		variableFlag := cmd.Flag("variable").Value
		log.Printf("Variables: %s", variableFlag)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pencil.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringArrayP("file", "f", make([]string, 0), "Path to a file you want run through the Template Engine")
	rootCmd.Flags().StringArrayP("folder", "F", make([]string, 0), "Path to a directory you want all files underneath run through the Template Engine")
	rootCmd.Flags().StringArrayP("variable", "v", make([]string, 0), "A variable to be used in the Templates. e.g. SECRET_KEY=something-secret")
}
