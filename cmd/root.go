package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wrkt",
	Short: "Git worktree management made simple",
	Long: `wrkt organizes git worktrees in a predictable structure and provides
seamless navigation with zsh integration.

All worktrees are organized in the worktrees/ subdirectory for easy management.`,
}

func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(listCmd)
}
