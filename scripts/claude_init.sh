#!/bin/bash

# Claude Session Initialization Script
# This script helps Claude validate the development environment and prevent dogfooding violations

set -e

echo "🚀 Claude Session Initialization"
echo "================================="

# Check if we're in the correct repository
if [[ ! -f "CLAUDE.md" ]]; then
    echo "❌ ERROR: CLAUDE.md not found"
    echo "   You must run this script from the wrkt repository root"
    exit 1
fi

# Display current state
echo ""
echo "📍 Current Location: $(pwd)"
echo "🌿 Current Branch: $(git branch --show-current 2>/dev/null || echo 'unknown')"

# Check for dogfooding violations
CURRENT_DIR=$(pwd)
CURRENT_BRANCH=$(git branch --show-current 2>/dev/null || echo "unknown")

if [[ "$CURRENT_DIR" == */wrkt ]] && [[ "$CURRENT_DIR" != */worktrees/* ]]; then
    if [[ "$CURRENT_BRANCH" != "main" ]]; then
        echo ""
        echo "🚨 DOGFOODING VIOLATION DETECTED!"
        echo "   Main directory is on branch: $CURRENT_BRANCH"
        echo "   Expected: main"
        echo ""
        echo "🔧 RECOMMENDED ACTION:"
        echo "   git checkout main"
        echo ""
        echo "⚠️  WARNING: This will switch to main branch in the main directory"
        echo "   Make sure to commit any important changes first"
        
        read -p "Fix violation automatically? (y/N): " auto_fix
        if [[ "$auto_fix" =~ ^[Yy]$ ]]; then
            git checkout main
            echo "✅ Fixed: Switched to main branch"
        else
            echo "⚠️  Violation not fixed - proceed with caution"
        fi
    else
        echo "✅ Dogfooding Status: CLEAN"
    fi
fi

# Display session context if available
echo ""
if [[ -f ".claude/CURRENT_STATE.md" ]]; then
    echo "📋 Session Context:"
    echo "-------------------"
    cat .claude/CURRENT_STATE.md
else
    echo "⚠️  No session context found"
    echo "   Creating new session state..."
    
    mkdir -p .claude
    cat > .claude/CURRENT_STATE.md <<EOF
# Claude Development State

**Last Updated**: $(date '+%Y-%m-%d %H:%M:%S %Z')

## Current Focus
- **Worktree**: main
- **Task**: New session started
- **Branch**: $(git branch --show-current)

## Session Context
- **Objective**: [Specify your objective]
- **Phase**: [Specify current phase]
- **Priority**: [High/Medium/Low]

## Recent Violations
- None recorded

## Progress Today
- Session initialized

## Next Actions
1. Define your current objective
2. Update this file with your task
3. Follow sequential task management

## Active Worktrees Status
$(./wrkt list --verbose 2>/dev/null || echo "Run './wrkt list --verbose' to see worktree status")

## Notes for Next Session
- Update this file before ending session
- Record any violations or lessons learned
EOF
    
    echo "✅ Created new session state file"
fi

# List worktrees
echo ""
echo "🌲 Worktree Status:"
echo "-------------------"
./wrkt list --verbose 2>/dev/null || echo "Error: wrkt command not available"

echo ""
echo "✅ Session initialization complete!"
echo ""
echo "🎯 Next Steps:"
echo "1. Review the session context above"
echo "2. Update .claude/CURRENT_STATE.md with your current task"
echo "3. Follow the SESSION START CHECKLIST in CLAUDE.md"
echo "4. Use sequential task management (one worktree at a time)"
echo ""
echo "Remember: Always validate your location before any git operations!"