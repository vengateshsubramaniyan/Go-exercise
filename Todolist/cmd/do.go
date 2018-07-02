package cmd

import (
	"Go-exercise/Todolist/db"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "To mark the task as completed.",
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int
		for _, arg := range args {
			val, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Printf("Error while parsing %s value", arg)
			} else {
				ids = append(ids, val)
			}
		}
		tasks := taskbucket.ListTask()
		for _, id := range ids {
			if id < 0 || id > len(tasks) {
				fmt.Printf("Task id %d is invalid\n", id)
				continue
			}
			taskbucket.DoTask(tasks[id-1].Key)
			taskbucket.AddToCompletedList(tasks[id-1].Value)
		}
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
