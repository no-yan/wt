# Implementation Guide: `wrkt switch -` (Previous Worktree) Functionality

## 🎯 Acceptance Criteria

The implementation is complete when:

1. **Command Usage**: `wrkt switch -` switches to the previously visited worktree
2. **Error Handling**: Shows "wrkt: no previous worktree" when no previous worktree exists
3. **Tab Completion**: `-` appears in tab completion for `wrkt switch`
4. **Terminal Isolation**: Each terminal session maintains its own previous worktree history
5. **Session Persistence**: Previous worktree is remembered until terminal session ends (same as `cd -`)
6. **Behavioral Consistency**: Works exactly like `cd -` but for worktrees

## 📋 Implementation Checklist

### Required Changes
- [ ] Modify zsh function in `cmd/shell_init.go` to handle `-` argument
- [ ] Add `-` to tab completion in `cmd/shell_init.go`
- [ ] Test all functionality manually
- [ ] Test error cases

### Files to Modify
- `cmd/shell_init.go` - Only file that needs changes

### Implementation Time
- **Estimated**: 30 minutes
- **Complexity**: Low (single file, ~20 lines of changes)

## 🏗️ Implementation Details

### Architecture Decision

**Approach**: Environment variable-based state management
**Rationale**: 
- Same mechanism as `cd -` (uses `OLDPWD`)
- Automatic terminal isolation
- No persistent state files needed
- Simple and reliable

### State Management

**Previous Worktree Storage**: `WRKT_OLDPWD` environment variable
**Scope**: Per terminal session (same as `cd -`)
**Lifetime**: Until terminal session ends
**Isolation**: Automatic (environment variables are process-specific)

## 📝 Detailed Implementation

### 1. Modify Zsh Function

**File**: `cmd/shell_init.go`
**Function**: `generateZshIntegration()`

**Current Code**:
```bash
function wrkt() {
  case "$1" in
    switch)
      if [ $# -eq 2 ]; then
        local target_path
        target_path=$(command wrkt switch "$2" 2>/dev/null)
        if [ $? -eq 0 ] && [ -n "$target_path" ]; then
          cd "$target_path"
        else
          command wrkt switch "$2"
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
```

**Updated Code**:
```bash
function wrkt() {
  case "$1" in
    switch)
      if [ $# -eq 2 ]; then
        if [ "$2" = "-" ]; then
          # Handle switch to previous worktree
          if [ -n "$WRKT_OLDPWD" ]; then
            cd "$WRKT_OLDPWD"
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
```

### 2. Update Tab Completion

**File**: `cmd/shell_init.go`
**Function**: `_wrkt_completion()`

**Current Code**:
```bash
case $words[2] in
  switch)
    local worktrees
    worktrees=($(command wrkt list 2>/dev/null | cut -f1))
    _values 'worktrees' $worktrees
    ;;
```

**Updated Code**:
```bash
case $words[2] in
  switch)
    local worktrees
    worktrees=($(command wrkt list 2>/dev/null | cut -f1))
    # Add the previous worktree option
    worktrees+=("-")
    _values 'worktrees' $worktrees
    ;;
```

## 🧪 Testing Procedures

### Manual Testing Steps

1. **Setup Test Environment**:
   ```bash
   # Ensure shell integration is loaded
   eval "$(wrkt shell-init)"
   
   # Create test worktrees
   wrkt add feature-auth
   wrkt add feature-login
   ```

2. **Test Basic Functionality**:
   ```bash
   # Switch to feature-auth
   wrkt switch feature-auth
   pwd  # Should show: .../worktrees/feature-auth
   
   # Switch to feature-login
   wrkt switch feature-login
   pwd  # Should show: .../worktrees/feature-login
   
   # Switch back to previous (should be feature-auth)
   wrkt switch -
   pwd  # Should show: .../worktrees/feature-auth
   
   # Switch back again (should be feature-login)
   wrkt switch -
   pwd  # Should show: .../worktrees/feature-login
   ```

3. **Test Error Handling**:
   ```bash
   # In a fresh terminal session
   wrkt switch -
   # Expected output: "wrkt: no previous worktree"
   # Expected exit code: 1
   ```

4. **Test Tab Completion**:
   ```bash
   # Type and press Tab
   wrkt switch <TAB>
   # Should show: feature-auth feature-login -
   ```

5. **Test Terminal Isolation**:
   ```bash
   # Terminal 1
   wrkt switch feature-auth
   wrkt switch feature-login
   
   # Terminal 2 (new session)
   wrkt switch -
   # Expected: "wrkt: no previous worktree"
   
   # Terminal 1
   wrkt switch -
   # Expected: switch to feature-auth
   ```

### Expected Behavior

| Action | Expected Result |
|--------|----------------|
| `wrkt switch -` (no previous) | Error: "wrkt: no previous worktree" |
| `wrkt switch feature-auth` → `wrkt switch feature-login` → `wrkt switch -` | Returns to feature-auth |
| `wrkt switch -` → `wrkt switch -` | Alternates between two worktrees |
| Tab completion on `wrkt switch ` | Shows worktree names + `-` |
| New terminal session | No previous worktree available |

## 🔍 Code Review Checklist

- [ ] Environment variable `WRKT_OLDPWD` is properly set before `cd`
- [ ] Error message matches expected format: "wrkt: no previous worktree"
- [ ] Tab completion includes `-` option
- [ ] Function handles both `-` and named worktree cases
- [ ] Error handling returns proper exit codes
- [ ] No persistent state files are created
- [ ] Code follows existing patterns in the file

## 🐛 Common Issues and Solutions

### Issue: Previous worktree not saved
**Cause**: `WRKT_OLDPWD` not set before `cd`
**Solution**: Ensure `export WRKT_OLDPWD="$PWD"` is called before `cd`

### Issue: Tab completion not working
**Cause**: Shell integration not reloaded
**Solution**: Run `eval "$(wrkt shell-init)"` after making changes

### Issue: Error message format inconsistent
**Cause**: Wrong error message format
**Solution**: Use exact format: "wrkt: no previous worktree"

## 🔄 Rollback Plan

If issues are found, simply revert `cmd/shell_init.go` to the previous version:
```bash
git checkout HEAD~1 -- cmd/shell_init.go
```

This implementation requires no database migrations, state file cleanup, or complex rollback procedures.

## 📖 Usage Examples

After implementation, the feature works as follows:

```bash
# Switch to a worktree
wrkt switch feature-auth

# Switch to another worktree  
wrkt switch main

# Switch back to previous worktree (feature-auth)
wrkt switch -

# Switch back again (main)
wrkt switch -

# In a new terminal session
wrkt switch -
# Output: wrkt: no previous worktree
```

## 🎯 Success Metrics

The implementation is successful when:
1. All manual tests pass
2. Tab completion includes `-` option
3. Error messages match expected format
4. Terminal isolation works correctly
5. Behavior matches `cd -` exactly