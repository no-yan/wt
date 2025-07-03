#!/bin/bash

# Test script for wrkt switch - functionality
# Tests the implementation of previous worktree switching

echo "======================================"
echo "wrkt switch - Functionality Test"
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

# Test 1: Build wrkt
echo "Test 1: Build wrkt binary"
if go build -o wrkt 2>/dev/null; then
    pass_test "Built wrkt binary successfully"
else
    fail_test "Failed to build wrkt binary"
    exit 1
fi

# Test 2: Check shell integration includes '-' in tab completion
echo
echo "Test 2: Tab completion includes '-'"
if ./wrkt shell-init | grep -q 'worktrees+=("-")'; then
    pass_test "Tab completion includes '-' option"
else
    fail_test "Tab completion missing '-' option"
fi

# Test 3: Error case - no previous worktree
echo
echo "Test 3: Error handling - no previous worktree"
# Test in clean zsh process
current_dir=$(pwd)
output=$(zsh -c "export PATH='$current_dir:\$PATH'; eval \"\$(wrkt shell-init)\" 2>/dev/null; unset WRKT_OLDPWD; wrkt switch -" 2>&1)
exit_code=$?

if [[ $exit_code -eq 1 ]] && [[ "$output" == "wrkt: no previous worktree" ]]; then
    pass_test "Correct error message and exit code"
else
    fail_test "Expected 'wrkt: no previous worktree' with exit code 1, got: '$output' (exit code: $exit_code)"
fi

# Test 4: Success case - switch to previous worktree
echo
echo "Test 4: Switch to previous worktree"
current_dir=$(pwd)
info "Testing switch to previous worktree (/tmp)"

# Test in clean zsh process
result=$(zsh -c "export PATH='$current_dir:\$PATH'; eval \"\$(wrkt shell-init)\" 2>/dev/null; export WRKT_OLDPWD='/tmp'; wrkt switch - >/dev/null 2>&1; echo \"\$(pwd)\"")

if [[ "$result" == "/tmp" ]]; then
    pass_test "Successfully switched to previous worktree"
else
    fail_test "Failed to switch to previous worktree (result: $result, expected: /tmp)"
fi

# Test 5: Verify shell integration function structure
echo
echo "Test 5: Shell integration function structure"
shell_code=$(./wrkt shell-init)

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
output=$(wrkt switch 2>&1)
exit_code=$?

if [[ $exit_code -eq 1 ]] && [[ "$output" == *"Usage:"* ]] && [[ "$output" == *"wrkt switch"* ]]; then
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
    echo -e "${GREEN}The wrkt switch - functionality is working correctly.${NC}"
    echo
    echo "Usage examples:"
    echo "  wrkt switch feature-branch    # Switch to feature-branch"
    echo "  wrkt switch main             # Switch to main"
    echo "  wrkt switch -                # Switch back to feature-branch"
    echo "  wrkt switch -                # Switch back to main"
    exit 0
else
    echo
    echo -e "${RED}❌ Some tests failed. Please check the implementation.${NC}"
    exit 1
fi