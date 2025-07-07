#!/bin/bash

# Test script for wt switch - functionality
# Tests the implementation of previous worktree switching

echo "======================================"
echo "wt switch - Functionality Test"
echo "======================================"
echo

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

TESTS_PASSED=0
TESTS_FAILED=0

pass_test() {
    echo -e "${GREEN}✓ $1${NC}"
    ((TESTS_PASSED++))
}

fail_test() {
    echo -e "${RED}✗ $1${NC}"
    ((TESTS_FAILED++))
}

info() {
    echo -e "${YELLOW}→ $1${NC}"
}

# Test 1: Build wt
echo "Test 1: Build wt binary"
if go build -o wt 2>/dev/null; then
    pass_test "Built wt binary successfully"
else
    fail_test "Failed to build wt binary"
    exit 1
fi

# Test 2: Check shell integration includes '-' in tab completion
echo
echo "Test 2: Tab completion includes '-'"
if ./wt shell-init | grep -q 'worktrees+=("-")'; then
    pass_test "Tab completion includes '-' option"
else
    fail_test "Tab completion missing '-' option"
fi

# Test 3: Error case - no previous worktree
echo
echo "Test 3: Error handling - no previous worktree"
# Test in clean zsh process
current_dir=$(pwd)
output=$(zsh -c "export PATH='$current_dir:\$PATH'; eval \"\$(wt shell-init)\" 2>/dev/null; unset WRKT_OLDPWD; wt switch -" 2>&1)
exit_code=$?

if [[ $exit_code -eq 1 ]] && [[ "$output" == "wt: no previous worktree" ]]; then
    pass_test "Correct error message and exit code"
else
    fail_test "Expected 'wt: no previous worktree' with exit code 1, got: '$output' (exit code: $exit_code)"
fi

# Test 4: Success case - switch to previous worktree
echo
echo "Test 4: Switch to previous worktree"
current_dir=$(pwd)
info "Testing switch to previous worktree (/tmp)"

# Test in clean zsh process
result=$(zsh -c "export PATH='$current_dir:\$PATH'; eval \"\$(wt shell-init)\" 2>/dev/null; export WRKT_OLDPWD='/tmp'; wt switch - >/dev/null 2>&1; echo \"\$(pwd)\"")

if [[ "$result" == "/tmp" ]]; then
    pass_test "Successfully switched to previous worktree"
else
    fail_test "Failed to switch to previous worktree (result: $result, expected: /tmp)"
fi

# Test 5: Verify shell integration function structure
echo
echo "Test 5: Shell integration function structure"
shell_code=$(./wt shell-init)

if echo "$shell_code" | grep -q 'if \[ "\$2" = "-" \]; then'; then
    pass_test "Shell function handles '-' argument"
else
    fail_test "Shell function missing '-' argument handling"
fi

if echo "$shell_code" | grep -q 'export WRKT_OLDPWD="\$PWD"'; then
    pass_test "Shell function sets WRKT_OLDPWD"
else
    fail_test "Shell function missing WRKT_OLDPWD setting"
fi

# Test 6: Usage message
echo
echo "Test 6: Usage message"
output=$(wt switch 2>&1)
exit_code=$?

if [[ $exit_code -eq 1 ]] && [[ "$output" == *"Usage:"* ]] && [[ "$output" == *"wt switch"* ]]; then
    pass_test "Correct usage message format"
else
    fail_test "Incorrect usage message: '$output'"
fi

# Test 7: Verify implementation matches cd - behavior
echo
echo "Test 7: Behavior consistency with cd -"
# Test that the logic follows the same pattern as cd -
if echo "$shell_code" | grep -q 'if \[ -n "\$WRKT_OLDPWD" \]; then'; then
    pass_test "Implementation checks WRKT_OLDPWD like cd - checks OLDPWD"
else
    fail_test "Implementation doesn't properly check WRKT_OLDPWD"
fi

# Test 8: Consecutive wt switch - behavior (A → B → A)
echo
echo "Test 8: Consecutive switch - behavior"
current_dir=$(pwd)
info "Testing consecutive wt switch - commands"

# Test A → B → A pattern
result=$(zsh -c "
export PATH='$current_dir:\$PATH'
eval \"\$(wt shell-init)\" 2>/dev/null
cd /tmp
export WRKT_OLDPWD='/var'
wt switch - >/dev/null 2>&1
wt switch - >/dev/null 2>&1
echo \"\$(pwd)\"
")

if [[ "$result" == "/tmp" ]]; then
    pass_test "Consecutive switch - correctly toggles between locations"
else
    fail_test "Consecutive switch - failed to toggle (result: $result, expected: /tmp)"
fi

# Final results
echo
echo "======================================"
echo "TEST RESULTS"
echo "======================================"
echo -e "Tests passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Tests failed: ${RED}$TESTS_FAILED${NC}"
echo -e "Total tests: $((TESTS_PASSED + TESTS_FAILED))"

if [[ $TESTS_FAILED -eq 0 ]]; then
    echo
    echo -e "${GREEN}🎉 All tests passed!${NC}"
    echo -e "${GREEN}The wt switch - functionality is working correctly.${NC}"
    echo
    echo "Usage examples:"
    echo "  wt switch feature-branch    # Switch to feature-branch"
    echo "  wt switch main             # Switch to main"
    echo "  wt switch -                # Switch back to feature-branch"
    echo "  wt switch -                # Switch back to main"
    exit 0
else
    echo
    echo -e "${RED}❌ Some tests failed. Please check the implementation.${NC}"
    exit 1
fi