#!/bin/bash

echo "🔍 Running firmware quality checks..."

# Static analysis with cppcheck
echo "=== Static Analysis ==="
cppcheck --enable=all --error-exitcode=1 \
  -I include/ -I board/ \
  src/ drivers/ startup/ || exit 1

# Check for duplicate code
echo "=== Duplicate Code Detection ==="
jscpd --min-lines 5 --min-tokens 50 \
  --reporters console \
  --threshold 10 \
  src/ drivers/ || exit 1

# MISRA C compliance (if available)
echo "=== MISRA C Compliance ==="
if command -v pc-lint-plus &> /dev/null; then
    pc-lint-plus -i include/ src/*.c drivers/*.c || exit 1
else
    echo "⚠️  PC-lint not available, skipping MISRA checks"
fi

# Check for dangerous functions
echo "=== Dangerous Function Check ==="
DANGEROUS_FUNCS="malloc|free|printf|sprintf|strcpy|strcat|gets"
if grep -r -E "$DANGEROUS_FUNCS" src/ drivers/ 2>/dev/null; then
    echo "❌ Dangerous functions found in code"
    exit 1
fi

echo "✅ All quality checks passed!"
