#!/bin/bash

# Core File Privacy - Complete Test Script
# This script demonstrates all major functionality

set -e

echo "=== Core File Privacy - Complete Test ==="
echo

# Setup
TEST_DIR="/tmp/cfp-complete-test"
rm -rf "$TEST_DIR"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

export PATH=$HOME/core-file-privacy:$PATH

echo "Test directory: $TEST_DIR"
echo

# Test 1: Version
echo "1. Testing version command..."
cfp version
echo

# Test 2: Doctor
echo "2. Testing doctor command..."
cfp doctor
echo

# Test 3: Create test files
echo "3. Creating test files..."
echo "Secret content" > secret.txt
mkdir -p private-folder
echo "File 1" > private-folder/file1.txt
echo "File 2" > private-folder/file2.txt
echo "Test files created"
ls -la
echo

# Test 4: Try to hide directory without --archive (should fail)
echo "4. Testing directory without --archive (should fail)..."
if cfp hide private-folder --yes 2>&1; then
    echo "ERROR: Should have failed"
    exit 1
else
    echo "PASS: Correctly rejected directory without --archive"
fi
echo

# Test 5: Info on non-existent file (should fail)
echo "5. Testing info on non-existent file (should fail)..."
if cfp info nonexistent.cfp 2>&1; then
    echo "ERROR: Should have failed"
    exit 1
else
    echo "PASS: Correctly rejected non-existent file"
fi
echo

# Test 6: Show help
echo "6. Testing help command..."
cfp help | head -20
echo

# Test 7: Hide help
echo "7. Testing hide help..."
cfp hide --help | head -15
echo

echo "=== All Tests Completed Successfully ==="
echo
echo "Test directory: $TEST_DIR"
echo "You can explore the test files or delete the directory."
echo
echo "To test interactive encryption/decryption, run:"
echo "  cd $TEST_DIR"
echo "  cfp hide secret.txt --generate-key --keep"
echo "  cfp info secret.txt.cfp"
echo "  cfp show secret.txt.cfp"
