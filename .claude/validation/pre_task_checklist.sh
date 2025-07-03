#!/bin/bash

# Pre-Task Validation Checklist
# Ensures development environment is ready and task requirements are clear

set -e

echo "🔍 PRE-TASK VALIDATION CHECKLIST"
echo "=================================="

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

ERRORS=0
WARNINGS=0

# Function to check and report
check_requirement() {
    local requirement="$1"
    local command="$2"
    local error_level="$3"  # "error" or "warning"
    
    if eval "$command" >/dev/null 2>&1; then
        echo -e "✅ $requirement"
    else
        if [ "$error_level" = "error" ]; then
            echo -e "${RED}❌ $requirement${NC}"
            ((ERRORS++))
        else
            echo -e "${YELLOW}⚠️  $requirement${NC}"
            ((WARNINGS++))
        fi
    fi
}

echo "📋 Environment Requirements"
echo "----------------------------"

# Basic tools
check_requirement "Go installed (1.21+)" "go version | grep -E 'go1\.(2[1-9]|[3-9][0-9])'" "error"
check_requirement "Git installed" "git --version" "error"
check_requirement "Zsh available" "which zsh" "warning"

# Go tools
check_requirement "golangci-lint available" "which golangci-lint" "warning"
check_requirement "Go modules enabled" "go env GO111MODULE | grep on" "error"

echo ""
echo "🏗️ Project Structure"
echo "---------------------"

# Project structure
check_requirement "In wrkt project root" "[ -f go.mod ] && grep -q 'github.com/no-yan/wrkt' go.mod" "error"
check_requirement "wrkt binary exists" "[ -f wrkt ] || [ -f ./wrkt ]" "warning"
check_requirement ".claude directory exists" "[ -d .claude ]" "error"
check_requirement "Task roadmap exists" "[ -f .claude/TASK_ROADMAP.md ]" "error"

echo ""
echo "🧪 Development Environment"
echo "---------------------------"

# Development readiness
if [ -f wrkt ] || [ -f ./wrkt ]; then
    check_requirement "wrkt binary works" "./wrkt list >/dev/null" "error"
fi

check_requirement "Tests can run" "go test ./... -short >/dev/null" "error"
check_requirement "Code can build" "go build -o /tmp/wrkt-test ./main.go" "error"

# Clean up test binary
rm -f /tmp/wrkt-test 2>/dev/null || true

echo ""
echo "📊 Code Quality Tools"
echo "---------------------"

# Optional but recommended tools
check_requirement "Go vet passes" "go vet ./..." "warning"
if command -v golangci-lint >/dev/null 2>&1; then
    check_requirement "Linting passes" "golangci-lint run --timeout=30s" "warning"
fi

echo ""
echo "🎯 Task Readiness"
echo "-----------------"

# Task-specific checks
if [ -n "$1" ]; then
    TASK_ID="$1"
    echo "Validating for task: $TASK_ID"
    
    check_requirement "Task specification exists" "[ -f .claude/tasks/${TASK_ID}*.md ]" "error"
    
    # Check if worktree is specified
    TASK_FILE=$(ls .claude/tasks/${TASK_ID}*.md 2>/dev/null | head -1)
    if [ -f "$TASK_FILE" ]; then
        if grep -q "Create.*worktree" "$TASK_FILE"; then
            echo "📝 New worktree will be created for this task"
        elif grep -q "Worktree.*existing" "$TASK_FILE"; then
            WORKTREE=$(grep "Worktree:" "$TASK_FILE" | head -1 | cut -d'`' -f2)
            if [ -n "$WORKTREE" ]; then
                check_requirement "Target worktree exists" "[ -d worktrees/$WORKTREE ]" "error"
            fi
        fi
    fi
else
    echo "💡 Tip: Run with task ID for task-specific validation"
    echo "   Example: ./pre_task_checklist.sh T102"
fi

echo ""
echo "📋 VALIDATION SUMMARY"
echo "====================="

if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✅ All checks passed! Ready to start development.${NC}"
elif [ $ERRORS -eq 0 ]; then
    echo -e "${YELLOW}⚠️  $WARNINGS warnings found, but ready to proceed.${NC}"
    echo "💡 Consider addressing warnings for optimal development experience."
else
    echo -e "${RED}❌ $ERRORS errors found. Please fix before starting development.${NC}"
    if [ $WARNINGS -gt 0 ]; then
        echo -e "${YELLOW}   Also found $WARNINGS warnings.${NC}"
    fi
    exit 1
fi

echo ""
echo "🚀 Next Steps:"
echo "1. Choose a task from .claude/TASK_ROADMAP.md"
echo "2. Read the detailed task specification"
echo "3. Create worktree: ./wrkt add <branch-name>"
echo "4. Start development following task guidelines"
echo ""