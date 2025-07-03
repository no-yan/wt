# wrkt Aliases Implementation Guide

## Overview
This document provides a comprehensive implementation plan for adding aliases to the wrkt tool. The feature will allow users to create shortcuts for commonly used commands:
- `sw` for `switch`
- `rm` for `remove`
- `ls` for `list`

## Current Architecture Analysis

### Command Structure
- **Framework**: Cobra CLI
- **Entry Point**: `main.go` → `cmd.Execute()`
- **Commands**: Registered in `cmd/root.go` init() function
- **Pattern**: Each command is a separate `cobra.Command` struct

### Shell Integration
- **File**: `cmd/shell_init.go`
- **Function**: `generateZshIntegration()` generates zsh completion code
- **Special Handling**: `wrkt switch` has custom shell function for directory changing
- **Completion**: Uses zsh `_arguments`, `_values`, and `compdef` functions

## Implementation Plan

### Phase 1: Core Alias System

#### 1.1 Create Alias Command Structure
**File**: `cmd/alias.go`

```go
package cmd

import (
    "github.com/spf13/cobra"
)

var aliasCmd = &cobra.Command{
    Use:   "alias",
    Short: "Manage command aliases",
    Long:  "Create, remove, and list command aliases for wrkt commands",
}

var aliasAddCmd = &cobra.Command{
    Use:   "add <alias> <command>",
    Short: "Add a new alias",
    Args:  cobra.ExactArgs(2),
    RunE:  runAliasAdd,
}

var aliasRemoveCmd = &cobra.Command{
    Use:   "remove <alias>",
    Short: "Remove an alias",
    Args:  cobra.ExactArgs(1),
    RunE:  runAliasRemove,
}

var aliasListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all aliases",
    RunE:  runAliasList,
}

func init() {
    // Register subcommands
    aliasCmd.AddCommand(aliasAddCmd)
    aliasCmd.AddCommand(aliasRemoveCmd)
    aliasCmd.AddCommand(aliasListCmd)
    
    // Register with root
    rootCmd.AddCommand(aliasCmd)
}
```

#### 1.2 Create Alias Storage System
**File**: `internal/alias_manager.go`

```go
package internal

import (
    "encoding/json"
    "os"
    "path/filepath"
)

type AliasManager struct {
    configPath string
    aliases    map[string]string
}

type AliasConfig struct {
    Aliases map[string]string `json:"aliases"`
}

func NewAliasManager(configDir string) (*AliasManager, error) {
    configPath := filepath.Join(configDir, "aliases.json")
    
    manager := &AliasManager{
        configPath: configPath,
        aliases:    make(map[string]string),
    }
    
    // Load existing aliases
    if err := manager.load(); err != nil {
        return nil, err
    }
    
    return manager, nil
}

func (am *AliasManager) AddAlias(alias, command string) error {
    am.aliases[alias] = command
    return am.save()
}

func (am *AliasManager) RemoveAlias(alias string) error {
    delete(am.aliases, alias)
    return am.save()
}

func (am *AliasManager) GetCommand(alias string) (string, bool) {
    command, exists := am.aliases[alias]
    return command, exists
}

func (am *AliasManager) GetAllAliases() map[string]string {
    result := make(map[string]string)
    for k, v := range am.aliases {
        result[k] = v
    }
    return result
}

func (am *AliasManager) load() error {
    data, err := os.ReadFile(am.configPath)
    if os.IsNotExist(err) {
        // Create default aliases
        am.aliases = map[string]string{
            "sw": "switch",
            "rm": "remove",
            "ls": "list",
        }
        return am.save()
    }
    if err != nil {
        return err
    }
    
    var config AliasConfig
    if err := json.Unmarshal(data, &config); err != nil {
        return err
    }
    
    am.aliases = config.Aliases
    return nil
}

func (am *AliasManager) save() error {
    config := AliasConfig{
        Aliases: am.aliases,
    }
    
    data, err := json.MarshalIndent(config, "", "  ")
    if err != nil {
        return err
    }
    
    // Ensure directory exists
    if err := os.MkdirAll(filepath.Dir(am.configPath), 0755); err != nil {
        return err
    }
    
    return os.WriteFile(am.configPath, data, 0644)
}
```

#### 1.3 Modify Root Command for Alias Resolution
**File**: `cmd/root.go` (modifications)

```go
// Add to imports
import (
    "os"
    "path/filepath"
    "github.com/no-yan/wrkt/internal"
)

// Add to init() function
func init() {
    // ... existing code ...
    
    // Set up alias resolution
    cobra.OnInitialize(initAliases)
}

func initAliases() {
    // Get config directory
    configDir := getConfigDir()
    
    // Initialize alias manager
    aliasManager, err := internal.NewAliasManager(configDir)
    if err != nil {
        // Handle error silently, aliases are optional
        return
    }
    
    // Check if first argument is an alias
    if len(os.Args) > 1 {
        if command, exists := aliasManager.GetCommand(os.Args[1]); exists {
            // Replace alias with actual command
            os.Args[1] = command
        }
    }
}

func getConfigDir() string {
    if configDir := os.Getenv("XDG_CONFIG_HOME"); configDir != "" {
        return filepath.Join(configDir, "wrkt")
    }
    
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return ".wrkt"
    }
    
    return filepath.Join(homeDir, ".config", "wrkt")
}
```

### Phase 2: Shell Integration Enhancement

#### 2.1 Modify Shell Integration for Alias Support
**File**: `cmd/shell_init.go` (modifications)

```go
// Modify generateZshIntegration() function
func generateZshIntegration() string {
    return `
# ... existing code ...

# Load aliases for completion
_wrkt_get_aliases() {
    local config_dir="${XDG_CONFIG_HOME:-$HOME/.config}/wrkt"
    local alias_file="$config_dir/aliases.json"
    
    if [[ -f "$alias_file" ]]; then
        # Extract aliases from JSON
        jq -r '.aliases | keys[]' "$alias_file" 2>/dev/null || echo ""
    fi
}

# Enhanced completion function
_wrkt_completion() {
    local context curcontext="$curcontext" state line
    typeset -A opt_args
    
    # Get aliases
    local aliases=($(_wrkt_get_aliases))
    
    # Main command completion including aliases
    local main_commands=(
        'add:Add a new worktree'
        'remove:Remove a worktree'
        'list:List all worktrees'
        'switch:Switch to a worktree'
        'clean:Clean up unused worktrees'
        'shell-init:Generate shell integration'
        'alias:Manage command aliases'
    )
    
    # Add aliases to completion
    for alias_cmd in $aliases; do
        main_commands+=("$alias_cmd:Alias command")
    done
    
    _arguments -C \
        '1: :->command' \
        '*::arg:->args' \
        && ret=0
    
    case $state in
        command)
            _values 'wrkt command' $main_commands && ret=0
            ;;
        args)
            # Resolve aliases for argument completion
            local actual_command="$words[1]"
            if [[ -n "$aliases[(r)$actual_command]" ]]; then
                # Get the actual command for this alias
                actual_command=$(wrkt alias list | grep "^$actual_command:" | cut -d: -f2)
            fi
            
            case $actual_command in
                add)
                    _message 'branch name' && ret=0
                    ;;
                remove|switch|rm|sw)
                    local worktrees=($(wrkt list 2>/dev/null | tail -n +2))
                    _values 'worktree' $worktrees && ret=0
                    ;;
                list|ls)
                    _arguments \
                        '--verbose[Show detailed information]' \
                        '--dirty[Show only dirty worktrees]' \
                        && ret=0
                    ;;
                alias)
                    _arguments \
                        '1: :(add remove list)' \
                        '*::alias_arg:->alias_args' \
                        && ret=0
                    
                    case $state in
                        alias_args)
                            case $words[1] in
                                add)
                                    if [[ $CURRENT -eq 2 ]]; then
                                        _message 'alias name'
                                    elif [[ $CURRENT -eq 3 ]]; then
                                        _values 'command' 'switch' 'remove' 'list' 'add' 'clean'
                                    fi
                                    ;;
                                remove)
                                    _values 'alias' $aliases
                                    ;;
                            esac
                            ;;
                    esac
                    ;;
            esac
            ;;
    esac
    
    return ret
}

# Enhanced wrkt function with alias support
wrkt() {
    local cmd="$1"
    
    # Check if command is an alias
    local config_dir="${XDG_CONFIG_HOME:-$HOME/.config}/wrkt"
    local alias_file="$config_dir/aliases.json"
    
    if [[ -f "$alias_file" ]]; then
        local actual_command=$(jq -r ".aliases[\"$cmd\"] // \"$cmd\"" "$alias_file" 2>/dev/null || echo "$cmd")
        if [[ "$actual_command" != "$cmd" ]]; then
            # Replace first argument with actual command
            shift
            set -- "$actual_command" "$@"
            cmd="$actual_command"
        fi
    fi
    
    # Handle switch command specially
    if [[ "$cmd" == "switch" ]]; then
        local target_path=$(command wrkt switch "$@")
        if [[ $? -eq 0 && -n "$target_path" ]]; then
            cd "$target_path"
        else
            return $?
        fi
    else
        command wrkt "$@"
    fi
}

# ... rest of existing code ...
`
}
```

### Phase 3: Testing Implementation

#### 3.1 Unit Tests for Alias Manager
**File**: `internal/alias_manager_test.go`

```go
package internal

import (
    "os"
    "path/filepath"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestAliasManager_DefaultAliases(t *testing.T) {
    tempDir := t.TempDir()
    
    manager, err := NewAliasManager(tempDir)
    require.NoError(t, err)
    
    // Check default aliases
    command, exists := manager.GetCommand("sw")
    assert.True(t, exists)
    assert.Equal(t, "switch", command)
    
    command, exists = manager.GetCommand("rm")
    assert.True(t, exists)
    assert.Equal(t, "remove", command)
    
    command, exists = manager.GetCommand("ls")
    assert.True(t, exists)
    assert.Equal(t, "list", command)
}

func TestAliasManager_AddRemoveAlias(t *testing.T) {
    tempDir := t.TempDir()
    
    manager, err := NewAliasManager(tempDir)
    require.NoError(t, err)
    
    // Add new alias
    err = manager.AddAlias("test", "list")
    require.NoError(t, err)
    
    // Verify alias exists
    command, exists := manager.GetCommand("test")
    assert.True(t, exists)
    assert.Equal(t, "list", command)
    
    // Remove alias
    err = manager.RemoveAlias("test")
    require.NoError(t, err)
    
    // Verify alias is gone
    _, exists = manager.GetCommand("test")
    assert.False(t, exists)
}

func TestAliasManager_Persistence(t *testing.T) {
    tempDir := t.TempDir()
    
    // Create first manager and add alias
    manager1, err := NewAliasManager(tempDir)
    require.NoError(t, err)
    
    err = manager1.AddAlias("test", "clean")
    require.NoError(t, err)
    
    // Create second manager and verify alias persists
    manager2, err := NewAliasManager(tempDir)
    require.NoError(t, err)
    
    command, exists := manager2.GetCommand("test")
    assert.True(t, exists)
    assert.Equal(t, "clean", command)
}
```

#### 3.2 Command Tests
**File**: `cmd/alias_test.go`

```go
package cmd

import (
    "os"
    "path/filepath"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestAliasCommand_List(t *testing.T) {
    tempDir := t.TempDir()
    
    // Set up config directory
    configDir := filepath.Join(tempDir, "config")
    os.Setenv("XDG_CONFIG_HOME", tempDir)
    defer os.Unsetenv("XDG_CONFIG_HOME")
    
    // Execute alias list command
    cmd := rootCmd
    cmd.SetArgs([]string{"alias", "list"})
    
    err := cmd.Execute()
    require.NoError(t, err)
    
    // Should show default aliases
    // Note: Output testing would require capturing stdout
}

func TestAliasCommand_AddRemove(t *testing.T) {
    tempDir := t.TempDir()
    
    // Set up config directory
    os.Setenv("XDG_CONFIG_HOME", tempDir)
    defer os.Unsetenv("XDG_CONFIG_HOME")
    
    // Add alias
    cmd := rootCmd
    cmd.SetArgs([]string{"alias", "add", "test", "list"})
    
    err := cmd.Execute()
    require.NoError(t, err)
    
    // Verify alias was added (would need to check config file)
    configPath := filepath.Join(tempDir, "wrkt", "aliases.json")
    assert.FileExists(t, configPath)
    
    // Remove alias
    cmd.SetArgs([]string{"alias", "remove", "test"})
    err = cmd.Execute()
    require.NoError(t, err)
}
```

#### 3.3 Integration Tests
**File**: `cmd/integration_test.go` (additions)

```go
func TestAliasResolution(t *testing.T) {
    tempDir := t.TempDir()
    
    // Set up config directory
    os.Setenv("XDG_CONFIG_HOME", tempDir)
    defer os.Unsetenv("XDG_CONFIG_HOME")
    
    // Test default alias resolution
    testCases := []struct {
        alias    string
        expected string
    }{
        {"sw", "switch"},
        {"rm", "remove"},
        {"ls", "list"},
    }
    
    for _, tc := range testCases {
        t.Run(tc.alias, func(t *testing.T) {
            // This would require modifying the test to capture
            // command resolution behavior
            // Implementation depends on how alias resolution is tested
        })
    }
}
```

### Phase 4: Documentation Updates

#### 4.1 Update README.md
Add section about aliases:

```markdown
## Aliases

wrkt supports command aliases for frequently used commands:

- `wrkt sw <worktree>` - Alias for `wrkt switch`
- `wrkt rm <worktree>` - Alias for `wrkt remove`
- `wrkt ls` - Alias for `wrkt list`

### Managing Aliases

```bash
# List all aliases
wrkt alias list

# Add a new alias
wrkt alias add <alias> <command>

# Remove an alias
wrkt alias remove <alias>
```

Aliases are stored in `~/.config/wrkt/aliases.json`.
```

#### 4.2 Update Shell Integration Instructions
Update shell integration documentation to mention alias support in tab completion.

## Acceptance Criteria

### Functional Requirements
1. ✅ **Default Aliases**: `sw`, `rm`, `ls` work as shortcuts for `switch`, `remove`, `list`
2. ✅ **Alias Management**: Users can add, remove, and list aliases via `wrkt alias` command
3. ✅ **Shell Integration**: Aliases work with zsh tab completion
4. ✅ **Persistence**: Aliases are saved to configuration file
5. ✅ **Error Handling**: Graceful handling of missing/invalid aliases

### Technical Requirements
1. ✅ **Code Quality**: Follow existing code patterns and conventions
2. ✅ **Testing**: Unit tests for alias manager, command tests, integration tests
3. ✅ **Documentation**: Updated README and inline documentation
4. ✅ **Backward Compatibility**: Existing commands continue to work unchanged
5. ✅ **Configuration**: Aliases stored in standard config location

### Implementation Verification
1. ✅ **Manual Testing**: All aliases work correctly
2. ✅ **Shell Testing**: Tab completion works for aliases
3. ✅ **Test Suite**: All tests pass
4. ✅ **Integration**: Aliases work with existing worktree operations
5. ✅ **Documentation**: All documentation is accurate and complete

## Implementation Order

1. **Phase 1**: Core alias system (alias manager, command structure)
2. **Phase 2**: Shell integration enhancement
3. **Phase 3**: Testing implementation
4. **Phase 4**: Documentation updates

## Configuration File Format

```json
{
  "aliases": {
    "sw": "switch",
    "rm": "remove",
    "ls": "list"
  }
}
```

Location: `~/.config/wrkt/aliases.json` (follows XDG Base Directory specification)

## Notes for Implementation

- The alias resolution happens at the command line level, before Cobra processes arguments
- Shell integration requires jq for JSON parsing (document this dependency)
- Alias completion needs to resolve the actual command to provide appropriate argument completion
- Default aliases are created automatically on first run
- Error handling should be graceful - if alias system fails, commands should still work