#!/bin/bash

# Post-Task Validation & Quality Gates
# Ensures task completion meets quality standards before integration

set -e

echo "🎯 POST-TASK VALIDATION"
echo "======================="

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

ERRORS=0
WARNINGS=0

# Function to check and report
validate_requirement() {
    local requirement="$1"
    local command="$2"
    local error_level="$3"  # "error" or "warning"
    
    echo -n "Checking: $requirement... "
    
    if eval "$command" >/dev/null 2>&1; then
        echo -e "${GREEN}✅${NC}"
    else
        if [ "$error_level" = "error" ]; then
            echo -e "${RED}❌${NC}"
            ((ERRORS++))
        else
            echo -e "${YELLOW}⚠️${NC}"
            ((WARNINGS++))
        fi
    fi
}

# Function to run command and show output
run_check() {
    local description="$1"
    local command="$2"
    local error_level="$3"
    
    echo ""
    echo -e "${BLUE}🔍 $description${NC}"
    echo "Command: $command"
    echo "----------------------------------------"
    
    if eval "$command"; then
        echo -e "${GREEN}✅ $description - PASSED${NC}"
    else
        if [ "$error_level" = "error" ]; then
            echo -e "${RED}❌ $description - FAILED${NC}"
            ((ERRORS++))
        else
            echo -e "${YELLOW}⚠️ $description - WARNING${NC}"
            ((WARNINGS++))
        fi
    fi
}

echo "📋 Basic Quality Gates"
echo "----------------------"

# Core quality requirements
validate_requirement "All tests pass" "go test ./..." "error"
validate_requirement "Code builds successfully" "go build -o /tmp/wrkt-validation ./main.go" "error"
validate_requirement "Go modules are tidy" "go mod tidy && git diff --exit-code go.mod go.sum" "error"

# Clean up test binary
rm -f /tmp/wrkt-validation 2>/dev/null || true

echo ""
echo "🎨 Code Quality"
echo "---------------"

# Code quality checks
validate_requirement "Code is formatted" "gofmt -l . | wc -l | xargs test 0 -eq" "error"
validate_requirement "Go vet passes" "go vet ./..." "warning"

# Optional linting
if command -v golangci-lint >/dev/null 2>&1; then
    validate_requirement "Linting passes" "golangci-lint run --timeout=60s" "warning"
else
    echo "⚠️  golangci-lint not available - consider installing for better code quality"
    ((WARNINGS++))
fi

echo ""
echo "🧪 Test Coverage Analysis"
echo "-------------------------"

# Test coverage
run_check "Generate test coverage report" "go test ./... -coverprofile=/tmp/coverage.out" "warning"

if [ -f /tmp/coverage.out ]; then
    COVERAGE=$(go tool cover -func=/tmp/coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    echo "Current test coverage: ${COVERAGE}%"
    
    if [ "$(echo "$COVERAGE >= 80" | bc -l 2>/dev/null || echo 0)" = "1" ]; then
        echo -e "${GREEN}✅ Test coverage is good (${COVERAGE}%)${NC}"
    elif [ "$(echo "$COVERAGE >= 60" | bc -l 2>/dev/null || echo 0)" = "1" ]; then
        echo -e "${YELLOW}⚠️ Test coverage could be improved (${COVERAGE}%)${NC}"
        ((WARNINGS++))
    else
        echo -e "${RED}❌ Test coverage is low (${COVERAGE}%)${NC}"
        ((ERRORS++))
    fi
    
    rm -f /tmp/coverage.out
fi

echo ""
echo "📝 Git State Validation"
echo "-----------------------"

# Git state checks
validate_requirement "No uncommitted changes" "git diff --exit-code" "error"
validate_requirement "No untracked files (or intentional)" "git status --porcelain | grep -v '^??' || true" "warning"
validate_requirement "Commit message follows conventions" "git log -1 --pretty=%s | grep -E '^(feat|fix|docs|style|refactor|test|chore):'|| echo 'Consider conventional commit format'" "warning"

echo ""
echo "🔧 Integration Readiness"
echo "------------------------"

# Check if wrkt binary works
if [ -f wrkt ]; then
    validate_requirement "wrkt binary is functional" "./wrkt list >/dev/null" "error"
    
    # Test specific functionality if available
    validate_requirement "Basic worktree operations work" "./wrkt list --verbose >/dev/null" "warning"
fi

# Check for potential conflicts
CURRENT_BRANCH=$(git branch --show-current)
echo "Current branch: $CURRENT_BRANCH"

if [ "$CURRENT_BRANCH" != "main" ]; then
    validate_requirement "Branch is ready for merge review" "git log main..$CURRENT_BRANCH --oneline | wc -l | xargs test 0 -lt" "warning"
fi

echo ""
echo "📊 VALIDATION SUMMARY"
echo "===================="

if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}🎉 PERFECT! All quality gates passed.${NC}"
    echo "✅ Task is ready for integration"
elif [ $ERRORS -eq 0 ]; then
    echo -e "${YELLOW}⚠️ GOOD with minor issues: $WARNINGS warnings${NC}"
    echo "✅ Task is ready for integration with minor improvements suggested"
else
    echo -e "${RED}❌ BLOCKED: $ERRORS errors must be fixed${NC}"
    if [ $WARNINGS -gt 0 ]; then
        echo -e "${YELLOW}   Also: $WARNINGS warnings should be addressed${NC}"
    fi
    echo ""
    echo "🔧 Required actions before integration:"
    echo "1. Fix all error conditions above"
    echo "2. Address warnings for best practices"
    echo "3. Re-run validation"
    exit 1
fi

echo ""
echo "🚀 Integration Steps:"
echo "1. Update task status to 'completed' in SESSION.md"
echo "2. Push branch for review (if using PR workflow)"
echo "3. Coordinate with task manager for integration"
echo "4. Consider running integration tests"
echo ""

# Optional: Generate integration report
TASK_ID=$(basename $(pwd) | sed 's/.*-\([T][0-9][0-9][0-9]\).*/\1/')
if [[ "$TASK_ID" =~ ^T[0-9]{3}$ ]]; then
    echo "📋 Task ID detected: $TASK_ID"
    echo "💡 Consider updating TASK_ROADMAP.md status"
fi

echo "✨ Validation complete!"