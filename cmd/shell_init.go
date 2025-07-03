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
  eval "$(wrkt shell-init)"

This enables the 'wrkt switch' command to actually change directories.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(generateZshIntegration())
	},
}

func generateZshIntegration() string {
	return `# wrkt shell integration for zsh
function wrkt() {
  case "$1" in
    switch)
      if [ $# -eq 2 ]; then
        if [ "$2" = "-" ]; then
          # Handle switch to previous worktree
          if [ -n "$WRKT_OLDPWD" ]; then
            local current_pwd="$PWD"
            cd "$WRKT_OLDPWD"
            export WRKT_OLDPWD="$current_pwd"
          else
            echo "wrkt: no previous worktree" >&2
            return 1
          fi
        else
          # Handle switch to named worktree
          local target_path
          target_path=$(command wrkt switch "$2" 2>/dev/null)
          if [ $? -eq 0 ] && [ -n "$target_path" ]; then
            # Save current location as previous
            export WRKT_OLDPWD="$PWD"
            cd "$target_path"
          else
            command wrkt switch "$2"
          fi
        fi
      else
        echo "Usage: wrkt switch <name>" >&2
        return 1
      fi
      ;;
    *)
      command wrkt "$@"
      ;;
  esac
}

# Tab completion for wrkt
function _wrkt_completion() {
  local state
  _arguments \
    '1: :->commands' \
    '*: :->args'
  
  case $state in
    commands)
      _values 'commands' \
        'add[Add a new worktree]' \
        'list[List all worktrees]' \
        'switch[Switch to a worktree]' \
        'shell-init[Generate shell integration]' \
        'help[Help about any command]'
      ;;
    args)
      case $words[2] in
        switch)
          local worktrees
          worktrees=($(command wrkt list 2>/dev/null | cut -f1))
          # Add the previous worktree option
          worktrees+=("-")
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

compdef _wrkt_completion wrkt
`
}
