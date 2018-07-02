package cmd

import (
	"Go-exercise/Todolist/db"
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all tasks on your todolist",
	Run: func(cmd *cobra.Command, args []string) {
		tasks := taskbucket.ListTask()
		for i, t := range tasks {
			fmt.Printf("%d: \"%s\"\n", i+1, t.Value)
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
