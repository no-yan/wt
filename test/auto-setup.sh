#!/bin/bash

# Test script to verify auto-setup feature functionality
set -e

echo "Testing auto-setup feature..."

# Create a temporary test directory
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

# Initialize a git repository
git init
git config user.email "test@example.com"
git config user.name "Test User"
echo "# Test repo" > README.md
git add README.md
git commit -m "Initial commit"

# Build wt binary
WRKT_DIR="/Users/noyan/ghq/github.com/no-yan/wt/worktrees/feature-auto-setup"
cd "$WRKT_DIR"
go build -o wt .

# Test auto-setup functionality
cd "$TEMP_DIR"
echo "Testing auto-setup functionality..."

# Use the built wt binary to add a worktree
"$WRKT_DIR/wt" add feature-test

# Check if worktrees directory was created
if [ -d "worktrees" ]; then
    echo "✅ worktrees directory was created"
else
    echo "❌ worktrees directory was NOT created"
    exit 1
fi

# Check if .gitignore entry was added
if grep -q "^worktrees/$" .gitignore; then
    echo "✅ .gitignore entry was added"
else
    echo "❌ .gitignore entry was NOT added"
    exit 1
fi

# Check if the worktree was actually created
if [ -d "worktrees/feature-test" ]; then
    echo "✅ worktree was created successfully"
else
    echo "❌ worktree was NOT created"
    exit 1
fi

# Test that running add again doesn't duplicate the .gitignore entry
"$WRKT_DIR/wt" add feature-test2

GITIGNORE_COUNT=$(grep -c "^worktrees/$" .gitignore)
if [ "$GITIGNORE_COUNT" -eq 1 ]; then
    echo "✅ .gitignore entry was not duplicated"
else
    echo "❌ .gitignore entry was duplicated (count: $GITIGNORE_COUNT)"
    exit 1
fi

# Clean up
cd /
rm -rf "$TEMP_DIR"

echo "🎉 All auto-setup tests passed!"