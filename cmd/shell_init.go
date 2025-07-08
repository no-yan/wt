package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var shellInitCmd = &cobra.Command{
	Use:   "shell-init",
	Short: "Generate shell integration code",
	Long: `Generate shell integration code for zsh.

Add this to your ~/.zshrc:
  eval "$(wt shell-init)"

This enables the 'wt switch' command to actually change directories.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(generateZshIntegration())
	},
}

func generateZshIntegration() string {
	return `# wt shell integration for zsh
function wt() {
  case "$1" in
    switch|sw)
      if [ $# -eq 2 ]; then
        if [ "$2" = "-" ]; then
          # Handle switch to previous worktree
          if [ -n "$WRKT_OLDPWD" ]; then
            local current_pwd="$PWD"
            cd "$WRKT_OLDPWD"
            export WRKT_OLDPWD="$current_pwd"
          else
            echo "wt: no previous worktree" >&2
            return 1
          fi
        else
          # Handle switch to named worktree
          local target_path
          target_path=$(command wt switch "$2" 2>/dev/null)
          if [ $? -eq 0 ] && [ -n "$target_path" ]; then
            # Save current location as previous
            export WRKT_OLDPWD="$PWD"
            cd "$target_path"
          else
            command wt switch "$2"
          fi
        fi
      else
        echo "Usage: wt switch <name>" >&2
        return 1
      fi
      ;;
    *)
      command wt "$@"
      ;;
  esac
}

# Tab completion for wt
function _wt_completion() {
  local state
  _arguments \
    '1: :->commands' \
    '*: :->args'
  
  case $state in
    commands)
      _values 'commands' \
        'add[Add a new worktree]' \
        'list[List all worktrees]' \
        'ls[List all worktrees (alias)]' \
        'switch[Switch to a worktree]' \
        'sw[Switch to a worktree (alias)]' \
        'remove[Remove a worktree]' \
        'rm[Remove a worktree (alias)]' \
        'shell-init[Generate shell integration]' \
        'help[Help about any command]'
      ;;
    args)
      case $words[2] in
        switch|sw)
          local worktrees
          worktrees=($(command wt list --names-only 2>/dev/null))
          # Add the previous worktree option
          worktrees+=("-")
          _values 'worktrees' $worktrees
          ;;
        remove|rm)
          local worktrees
          worktrees=($(command wt list --names-only 2>/dev/null))
          _values 'worktrees' $worktrees
          ;;
        add)
          # Complete with branch names
          local branches
          branches=($(git branch -r 2>/dev/null | sed 's|.*origin/||' | grep -v HEAD))
          _values 'branches' $branches
          ;;
      esac
      ;;
  esac
}

compdef _wt_completion wt
`
}
