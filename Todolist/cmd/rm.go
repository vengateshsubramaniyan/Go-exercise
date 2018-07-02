package cmd

import (
	"Go-exercise/Todolist/db"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "rm is used to delete a task from the Todolist.",
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int
		for _, arg := range args {
			val, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Printf("Error while parsing %s value\n", arg)
			}
			ids = append(ids, val)
		}
		tasks := taskbucket.ListTask()
		for _, id := range ids {
			if id < 0 || id > len(tasks) {
				fmt.Printf("%d is a Invalid task number.\n", id)
				continue
			}
			taskbucket.DoTask(tasks[id-1].Key)
		}
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
