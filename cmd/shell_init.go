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
          # Handle switch to previous workspace
          if [ -n "$WRKT_OLDPWD" ]; then
            local current_pwd="$PWD"
            cd "$WRKT_OLDPWD"
            export WRKT_OLDPWD="$current_pwd"
          else
            echo "wt: no previous workspace" >&2
            return 1
          fi
        else
          # Handle switch to named workspace
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

# Detect VCS type (git or jj)
function _wt_detect_vcs() {
  if [ -d ".jj" ]; then
    echo "jj"
  elif [ -d ".git" ] || [ -f ".git" ]; then
    echo "git"
  else
    # Walk up to find repo root
    local dir="$PWD"
    while [ "$dir" != "/" ]; do
      if [ -d "$dir/.jj" ]; then
        echo "jj"
        return
      elif [ -d "$dir/.git" ] || [ -f "$dir/.git" ]; then
        echo "git"
        return
      fi
      dir="$(dirname "$dir")"
    done
    echo "unknown"
  fi
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
        'add[Add a new workspace]' \
        'list[List all workspaces]' \
        'ls[List all workspaces (alias)]' \
        'switch[Switch to a workspace]' \
        'sw[Switch to a workspace (alias)]' \
        'remove[Remove a workspace]' \
        'rm[Remove a workspace (alias)]' \
        'clean[Clean up stale workspaces]' \
        'shell-init[Generate shell integration]' \
        'help[Help about any command]'
      ;;
    args)
      case $words[2] in
        switch|sw)
          local workspaces
          workspaces=($(command wt list --names-only 2>/dev/null))
          # Add the previous workspace option
          workspaces+=("-")
          _values 'workspaces' $workspaces
          ;;
        remove|rm)
          local workspaces
          workspaces=($(command wt list --names-only 2>/dev/null | grep -v -E '^(main|master|default)$'))
          _values 'workspaces' $workspaces
          ;;
        add)
          local vcs_type
          vcs_type=$(_wt_detect_vcs)
          if [ "$vcs_type" = "jj" ]; then
            # For jj, suggest existing workspace names or let user type new name
            # Also provide -r flag completion
            if [[ "$words[-1]" == "-r" ]] || [[ "$words[-2]" == "-r" ]]; then
              # Complete with jj bookmarks
              local bookmarks
              bookmarks=($(jj bookmark list 2>/dev/null | awk '{print $1}'))
              _values 'revisions' $bookmarks 'trunk()' '@'
            else
              _arguments \
                '-r[Base revision]:revision:' \
                '--revision[Base revision]:revision:'
            fi
          else
            # For git, complete with branch names
            local branches
            branches=($(git branch -r 2>/dev/null | sed 's|.*origin/||' | grep -v HEAD))
            _values 'branches' $branches
          fi
          ;;
      esac
      ;;
  esac
}

compdef _wt_completion wt
`
}
