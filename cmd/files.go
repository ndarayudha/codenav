package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yofan2408/codenav/files"
)

var Path string
var Pattern string

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Show list of path file and line of code",
	Run: func(cmd *cobra.Command, args []string) {
		files.ReadDir(Path, Pattern)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVarP(&Path, "folder", "f", "", "Limit the number of files returned")
	err := searchCmd.MarkFlagRequired("folder")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	searchCmd.Flags().StringVarP(&Pattern, "pattern", "p", "", "Minimum size for files in search in MB.")
	err = searchCmd.MarkFlagRequired("pattern")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
