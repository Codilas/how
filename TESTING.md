# Testing Documentation for "how" CLI

## Overview

The `how` CLI project includes comprehensive test coverage for the `internal/config` package, which handles configuration loading, saving, and management.

## Test Suite Structure

### Test File Location
- **Path:** `internal/config/config_test.go`
- **Size:** 2,871 lines
- **Functions Tested:** 3 (Load, Save, getConfigDir)
- **Structs Tested:** 5 (Config, ProviderConfig, ContextConfig, DisplayConfig, HistoryConfig)
- **Total Test Cases:** 71+

### Package Structure
```
internal/config/
├── config.go          # Source code
├── config_test.go     # Test suite
```

## Functions Under Test

### 1. Load(configFile string) (*Config, error)
Loads configuration from YAML file with support for:
- Explicit file paths
- Default config directory (~/.config/how)
- Environment variable overrides (HOW_ prefix)
- Missing file handling
- YAML parsing with error reporting

**Tests:** 14+ test cases covering success paths, error conditions, and edge cases

### 2. Config.Save(configFile string) error
Saves configuration to YAML file with support for:
- Explicit file paths
- Default config directory
- Parent directory creation
- File overwriting
- Permission handling

**Tests:** 21+ test cases covering directory creation, permissions, and error conditions

### 3. getConfigDir() (string, error)
Internal helper function that returns the config directory path (~/.config/how)

**Tests:** 12+ test cases covering path construction, integration, and edge cases

## Struct Coverage

### Config
Main configuration struct with all application settings

**Test Coverage:**
- Field initialization
- YAML marshaling/unmarshaling
- Nested struct handling
- Round-trip consistency

### ProviderConfig
Configuration for AI providers (e.g., Anthropic)

**Test Coverage:**
- Field initialization with required/optional fields
- Custom headers marshaling
- Float precision handling
- Omitempty field behavior

### ContextConfig
Configuration for context gathering (files, history, git, etc.)

**Test Coverage:**
- Slice marshaling (empty and populated)
- Field initialization
- YAML tag validation

### DisplayConfig
Configuration for display options (syntax highlighting, colors, emoji)

**Test Coverage:**
- Boolean field handling
- Round-trip marshaling

### HistoryConfig
Configuration for command history management

**Test Coverage:**
- Field initialization
- File path handling
- Round-trip marshaling

## Test Execution

### Running Tests

#### All Tests
```bash
make test
```

#### Tests with Race Detection
```bash
make test-race
```

#### Tests with Coverage Report
```bash
make test-coverage
```

#### Tests with HTML Coverage Report
```bash
make test-coverage-html
```

#### Specific Package Tests
```bash
go test -v ./internal/config
```

#### Specific Test Function
```bash
go test -v -run TestLoad_SuccessfulLoadFromExplicitFile ./internal/config
```

## Go Testing Conventions

The test suite follows Go testing conventions:

### 1. **File Naming**
- Test file: `config_test.go` (follows `_test.go` suffix)
- Located in same package as source code

### 2. **Function Naming**
- Format: `Test<FunctionName>_<TestCase>`
- Examples:
  - `TestLoad_SuccessfulLoadFromExplicitFile`
  - `TestSave_RoundTripConsistency`
  - `TestGetConfigDir_CorrectPathConstruction`

### 3. **Test Isolation**
- Each test uses `t.TempDir()` for isolated temporary directories
- Tests don't depend on execution order
- No global state modifications

### 4. **Error Handling**
- `t.Fatalf()` for critical errors that stop the test
- `t.Errorf()` for assertion failures
- Descriptive error messages for debugging

### 5. **Documentation**
- Each test has a comment describing its purpose
- Example:
  ```go
  // TestLoad_SuccessfulLoadFromExplicitFile tests loading config from an explicit file path
  func TestLoad_SuccessfulLoadFromExplicitFile(t *testing.T) {
  ```

### 6. **Test Independence**
- Tests can run in any order
- Each test sets up its own fixtures
- Automatic cleanup via `t.TempDir()`

## Edge Cases Covered

### File System Operations
- ✅ Missing config files
- ✅ Invalid file paths
- ✅ Permission denied errors
- ✅ Parent directory creation
- ✅ File overwriting
- ✅ Special characters in paths
- ✅ Unicode in file paths

### YAML Parsing
- ✅ Empty files
- ✅ Invalid YAML syntax
- ✅ Invalid indentation
- ✅ Invalid data types
- ✅ Null values
- ✅ Duplicate keys
- ✅ Special characters and Unicode
- ✅ Large configuration files

### Configuration Handling
- ✅ Minimal configuration (required fields only)
- ✅ Full configuration (all fields)
- ✅ Empty providers map
- ✅ Nil providers map
- ✅ Environment variable overrides
- ✅ Load-modify-save consistency

### Data Validation
- ✅ Float field precision
- ✅ Custom headers marshaling
- ✅ Slice marshaling (empty and populated)
- ✅ Nested structure marshaling
- ✅ Round-trip marshaling consistency

## Coverage Metrics

### Function Coverage
- Load() - ✅ 100%
- Config.Save() - ✅ 100%
- getConfigDir() - ✅ 100%

### Struct Coverage
- Config - ✅ 100%
- ProviderConfig - ✅ 100%
- ContextConfig - ✅ 100%
- DisplayConfig - ✅ 100%
- HistoryConfig - ✅ 100%

### Expected Coverage
- **Goal:** 100% for public API
- **Current:** Comprehensive coverage of all functions and edge cases

## Test Verification

### Automated Verification Script
A verification script is available to validate test coverage:

```bash
./scripts/verify_test_coverage.sh
```

This script:
1. Verifies test file exists
2. Counts test functions
3. Runs all tests
4. Generates coverage report
5. Validates Go testing conventions

## Integration with CI/CD

The test suite is integrated into the project's standard test command:

```bash
make test
```

This command runs all tests in the project:
```bash
go test -v ./...
```

For continuous integration, use:
```bash
make test-coverage-html
```

This generates `coverage.html` with detailed coverage visualization.

## Troubleshooting

### Tests Fail Due to Permissions
Some tests verify permission error handling. If running as root or with unusual permissions:
```bash
make test-race  # Run with race detection for other issues
```

### Large Config File Tests
Tests that use large configuration files may take longer to run:
```bash
go test -v -timeout 30s ./internal/config
```

### Coverage Report Generation
If coverage report generation fails:
```bash
# Generate coverage data
go test -coverprofile=coverage.out ./internal/config

# View in terminal
go tool cover -func=coverage.out

# Generate HTML (requires graphviz on some systems)
go tool cover -html=coverage.out -o coverage.html
```

## Best Practices for New Tests

When adding new tests, follow these patterns:

### 1. Use Table-Driven Tests for Multiple Cases
```go
tests := []struct {
    name    string
    input   string
    want    string
    wantErr bool
}{
    {"case1", "input1", "output1", false},
    {"case2", "input2", "output2", false},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // test logic
    })
}
```

### 2. Use t.TempDir() for File Operations
```go
tmpDir := t.TempDir()  // Automatically cleaned up
configFile := filepath.Join(tmpDir, "config.yaml")
```

### 3. Clear Error Messages
```go
if got != want {
    t.Errorf("Got %v, want %v", got, want)
}
```

### 4. Test Both Success and Failure Cases
```go
// Success case
if err != nil {
    t.Fatalf("unexpected error: %v", err)
}

// Failure case
if err == nil {
    t.Fatalf("expected error, got nil")
}
```

## References

- [Go Testing Package Documentation](https://golang.org/pkg/testing/)
- [Table-Driven Tests in Go](https://golang.org/doc/tutorial/add-a-test)
- [Go Code Review Comments - Testing](https://golang.org/doc/effective_go#testing)

## Future Improvements

1. **Benchmarking:** Add performance benchmarks for Load/Save operations
2. **Concurrency Tests:** Add tests for concurrent Load/Save scenarios
3. **Property-Based Testing:** Consider using property-based testing libraries
4. **Fuzz Testing:** Add fuzz testing for YAML parsing
5. **Integration Tests:** Additional integration tests with CLI

## Contact

For questions about testing or to report issues:
- File an issue on [GitHub](https://github.com/Codilas/how)
- Check the main [README.md](./README.md) for general project information
