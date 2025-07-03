package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wrkt",
	Short: "Git worktree management made simple",
	Long: `wrkt organizes git worktrees in a predictable structure and provides
seamless navigation with zsh integration.

All worktrees are organized in the worktrees/ subdirectory for easy management.

Built-in aliases:
  sw  - alias for switch
  rm  - alias for remove
  ls  - alias for list

Use "wrkt [command] --help" for more information about a command.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(switchCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(shellInitCmd)
}
