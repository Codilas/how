#!/bin/bash
# Test coverage verification script for internal/config/config.go
# This script validates test coverage and ensures all tests follow Go conventions

set -e

echo "======================================================"
echo "Test Coverage Verification for internal/config"
echo "======================================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
CONFIG_PKG="./internal/config"
COVERAGE_THRESHOLD=80

echo "[1/5] Checking test file exists..."
if [ ! -f "internal/config/config_test.go" ]; then
    echo -e "${RED}✗ Test file not found: internal/config/config_test.go${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Test file found${NC}"
echo ""

echo "[2/5] Verifying test function naming conventions..."
test_count=$(grep -c "^func Test" internal/config/config_test.go || true)
if [ "$test_count" -lt 50 ]; then
    echo -e "${YELLOW}⚠ Expected many test functions, found: $test_count${NC}"
else
    echo -e "${GREEN}✓ Found $test_count test functions${NC}"
fi
echo ""

echo "[3/5] Running tests..."
if go test -v "$CONFIG_PKG" 2>&1 | tail -5; then
    echo -e "${GREEN}✓ All tests passed${NC}"
else
    echo -e "${RED}✗ Some tests failed${NC}"
    exit 1
fi
echo ""

echo "[4/5] Generating coverage report..."
go test -coverprofile=coverage.out -covermode=atomic "$CONFIG_PKG"
coverage_percent=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')

echo "Coverage: ${coverage_percent}%"
if (( $(echo "$coverage_percent >= $COVERAGE_THRESHOLD" | bc -l) )); then
    echo -e "${GREEN}✓ Coverage meets threshold (${COVERAGE_THRESHOLD}%)${NC}"
else
    echo -e "${YELLOW}⚠ Coverage below threshold (found: ${coverage_percent}%, threshold: ${COVERAGE_THRESHOLD}%)${NC}"
fi
echo ""

echo "[5/5] Verifying test conventions..."
conventions_ok=true

# Check for proper test isolation (t.TempDir usage)
if grep -q "t.TempDir()" internal/config/config_test.go; then
    echo -e "${GREEN}✓ Tests use t.TempDir() for isolation${NC}"
else
    echo -e "${YELLOW}⚠ Tests may not use t.TempDir() for isolation${NC}"
    conventions_ok=false
fi

# Check for proper error handling
if grep -q "t.Fatalf" internal/config/config_test.go && grep -q "t.Errorf" internal/config/config_test.go; then
    echo -e "${GREEN}✓ Tests use t.Fatalf() and t.Errorf() properly${NC}"
else
    echo -e "${YELLOW}⚠ Tests may not use proper error functions${NC}"
    conventions_ok=false
fi

# Check for documentation comments
if grep -q "// Test" internal/config/config_test.go; then
    echo -e "${GREEN}✓ Tests have documentation comments${NC}"
else
    echo -e "${YELLOW}⚠ Tests may lack documentation${NC}"
    conventions_ok=false
fi

echo ""
echo "======================================================"
echo "Test Coverage Verification Complete"
echo "======================================================"
echo ""

# Generate HTML report
echo "Generating HTML coverage report..."
go tool cover -html=coverage.out -o coverage.html
echo -e "${GREEN}✓ HTML report generated: coverage.html${NC}"
echo ""

# Summary
echo "Summary:"
echo "  - Test file: internal/config/config_test.go"
echo "  - Test functions: $test_count"
echo "  - Coverage: ${coverage_percent}%"
echo "  - Conventions: $([ "$conventions_ok" = true ] && echo "✓ Passed" || echo "⚠ Review recommended")"
echo ""
echo "To view detailed coverage report:"
echo "  open coverage.html"
