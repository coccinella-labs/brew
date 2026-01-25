#!/bin/bash

echo "🚀 Setting up coffee maker firmware development environment..."

# Install pre-commit if not available
if ! command -v pre-commit &> /dev/null; then
    echo "📦 Installing pre-commit..."
    pip3 install pre-commit
fi

# Install pre-commit hooks
echo "🔧 Installing pre-commit hooks..."
pre-commit install

# Install ARM toolchain if not available
if ! command -v arm-none-eabi-gcc &> /dev/null; then
    echo "🔧 Installing ARM toolchain..."
    if [[ "$OSTYPE" == "darwin"* ]]; then
        brew install arm-none-eabi-gcc
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        sudo apt-get update
        sudo apt-get install -y gcc-arm-none-eabi
    fi
fi

# Install quality tools
echo "🔍 Installing quality analysis tools..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    brew install cppcheck
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    sudo apt-get install -y cppcheck
fi

# Install Node.js tools for duplicate detection
if command -v npm &> /dev/null; then
    npm install -g jscpd lizard
fi

# Make scripts executable
chmod +x scripts/*.sh
chmod +x scripts/*.py

# Initial build test
echo "🏗️ Testing initial build..."
make clean && make

echo "✅ Development environment setup complete!"
echo ""
echo "Next steps:"
echo "1. Connect STM32F4 hardware"
echo "2. Run 'make flash' to deploy firmware"
echo "3. Run 'make debug' for debugging session"
echo "4. All commits will automatically run quality checks"
