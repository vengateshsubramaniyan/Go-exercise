package cmd

import (
	"Go-exercise/Todolist/db"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a new entry to the todolist",
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")
		if task == "" {
			fmt.Println("Task should not be a empty string")
			os.Exit(1)
		}
		taskbucket.AddTask(strings.TrimSpace(task))
		fmt.Println("Task added successfully to the Todolist")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
