package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all worktrees",
	Long:  "List all git worktrees with their status and branch information.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("List command - to be implemented")
	},
}
