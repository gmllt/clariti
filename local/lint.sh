#!/bin/bash
# lint.sh - Run golangci-lint on the project

echo "Running golangci-lint..."
golangci-lint run

if [ $? -eq 0 ]; then
    echo "✅ All linting checks passed!"
else
    echo "❌ Linting issues found. Please fix them before committing."
    exit 1
fi
