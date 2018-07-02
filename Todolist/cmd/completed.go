package cmd

import (
	"Go-exercise/Todolist/db"
	"fmt"

	"github.com/spf13/cobra"
)

var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "List all the completed task",
	Run: func(cmd *cobra.Command, args []string) {
		tasks := taskbucket.ListCompletedTask()
		for i, task := range tasks {
			fmt.Printf("%d: \"%s\"\n", i+1, task.Value)
		}
	},
}

func init() {
	rootCmd.AddCommand(completedCmd)
}
