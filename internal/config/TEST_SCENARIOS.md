# Test Scenarios for config.go

## Overview
This document provides detailed test scenarios for complete coverage of the config package. Each scenario includes setup, inputs, expected outputs, and edge cases.

---

## 1. Load() Function Test Scenarios

### 1.1 Successful Config Loading

#### TC1.1.1: Load from explicit file path with complete configuration
**Setup**: Create temporary config file with all fields populated
**Input**: `Load("/tmp/test_config.yaml")`
**Expected Output**:
- Returns non-nil Config pointer
- All fields properly unmarshaled
- No error returned
**Verification**:
- Config.CurrentProvider matches YAML value
- All Providers properly loaded
- Nested configs (Context, Display, History) correct
- Can access specific nested values

**Test Data**:
```yaml
currentProvider: anthropic
providers:
  anthropic:
    type: anthropic
    apiKey: sk-test-123
    model: claude-3-opus
    maxTokens: 4096
    temperature: 0.7
context:
  includeFiles: true
  includeHistory: 10
  includeEnvironment: true
  includeGit: true
  maxContextSize: 100000
  excludePatterns:
    - "*.log"
    - "node_modules/*"
display:
  syntaxHighlight: true
  showContext: true
  emoji: true
  color: true
history:
  enabled: true
  maxSize: 1000
  filePath: ~/.how/history
```

#### TC1.1.2: Load from default location
**Setup**: Create config at ~/.config/how/config.yaml with sample data
**Input**: `Load("")`
**Expected Output**:
- Returns non-nil Config pointer
- Loads from default location
- Data matches file content
**Verification**:
- getConfigDir() called internally
- Config loaded successfully
- Matches content of ~/.config/how/config.yaml

**Edge Case**: Test with actual user home directory

#### TC1.1.3: Load minimal configuration
**Setup**: Config file with only required fields
**Input**: `Load("/tmp/minimal_config.yaml")`
**Test Data**:
```yaml
currentProvider: anthropic
providers:
  anthropic:
    type: anthropic
    apiKey: sk-test
    model: claude-3
    maxTokens: 2048
```
**Expected**: Loads successfully, optional fields are zero values

#### TC1.1.4: Load config with only CurrentProvider
**Setup**: Config file with only CurrentProvider set
**Test Data**:
```yaml
currentProvider: anthropic
```
**Expected**: Loads successfully, other fields are empty/nil

#### TC1.1.5: Load empty config file
**Setup**: Valid YAML file with no content or `{}`
**Input**: `Load("/tmp/empty_config.yaml")`
**Expected Output**:
- Returns Config pointer with all zero values
- No error
- CurrentProvider is empty string
- Providers map is nil/empty

### 1.2 Missing Configuration File

#### TC1.2.1: File doesn't exist at default location
**Setup**: ~/.config/how/config.yaml doesn't exist
**Input**: `Load("")`
**Expected Output**:
- Returns non-nil Config (empty/default values)
- No error (ConfigFileNotFoundError is expected)
- Warning message printed to stdout
**Verification**:
- Should NOT return error for missing file
- Config object is still usable
- All fields are zero values

#### TC1.2.2: File doesn't exist at explicit path
**Setup**: /nonexistent/path/config.yaml
**Input**: `Load("/nonexistent/path/config.yaml")`
**Expected Output**:
- Warning issued
- Returns empty Config
- No error returned
**Verification**:
- Graceful handling of missing file

#### TC1.2.3: Parent directory doesn't exist
**Setup**: ~/.config/how/ doesn't exist
**Input**: `Load("")`
**Expected Output**:
- Tries to read from ~/.config/how/config.yaml
- File not found, returns warning
- Returns empty Config
**Verification**:
- Path traversal handled correctly

### 1.3 Invalid YAML

#### TC1.3.1: Malformed YAML syntax
**Setup**: Config file with invalid YAML
**Test Data**:
```yaml
currentProvider: anthropic
providers:
  anthropic
    type: anthropic  # Missing colon
```
**Input**: `Load("/tmp/invalid_config.yaml")`
**Expected Output**:
- Error returned
- Error contains "failed to read config file"
- Config is nil
**Verification**:
- Error wrapping works correctly
- Error message is informative

#### TC1.3.2: Incorrect indentation
**Test Data**:
```yaml
currentProvider: anthropic
providers:
  anthropic:
   type: anthropic  # Wrong indentation
   apiKey: test
```
**Expected Output**: Error returned for invalid YAML

#### TC1.3.3: Tab characters in YAML
**Test Data**: YAML with tabs instead of spaces
**Expected Output**: YAML parser error returned

### 1.4 Invalid Data Types

#### TC1.4.1: String value for integer field
**Test Data**:
```yaml
providers:
  anthropic:
    type: anthropic
    maxTokens: "not a number"
```
**Expected Output**: Unmarshal error returned

#### TC1.4.2: String value for boolean field
**Test Data**:
```yaml
display:
  color: "maybe"
```
**Expected Output**: Unmarshal error or defaults to false

#### TC1.4.3: Invalid float for Temperature
**Test Data**:
```yaml
providers:
  anthropic:
    temperature: "hot"
```
**Expected Output**: Unmarshal error returned

#### TC1.4.4: List instead of map for Providers
**Test Data**:
```yaml
providers:
  - type: anthropic
```
**Expected Output**: Unmarshal error returned

### 1.5 Environment Variable Overrides

#### TC1.5.1: Override CurrentProvider via env var
**Setup**: HOW_CURRENTPROVIDER=openai
**File Content**: currentProvider: anthropic
**Expected Output**:
- Config.CurrentProvider = "openai"
- File value is overridden
**Verification**:
- Environment variables take precedence

#### TC1.5.2: Multiple environment variable overrides
**Setup**:
- HOW_CURRENTPROVIDER=custom
- HOW_CONTEXT_INCLUDEFILES=true (if viper supports this)
**Expected**: Both overrides applied

#### TC1.5.3: Nested field overrides (if viper supports)
**Setup**: HOW_DISPLAY_COLOR=false (override display.color)
**Expected**: Nested field updated

#### TC1.5.4: Invalid env var value type
**Setup**: HOW_PROVIDERS_ANTHROPIC_MAXTOKENS=notanumber
**Expected**: Handled gracefully (ignored or error)

### 1.6 Error Cases

#### TC1.6.1: getConfigDir() fails (no home directory)
**Setup**: Mock os.UserHomeDir() to return error
**Input**: `Load("")`
**Expected Output**:
- Error returned
- Error contains "failed to get config directory"
- Config is nil

#### TC1.6.2: Config file not readable (permission denied)
**Setup**: Create file with 000 permissions
**Input**: `Load("/tmp/unreadable_config.yaml")`
**Expected Output**:
- Error returned
- Config is nil

#### TC1.6.3: Unmarshal fails after read succeeds
**Setup**: Valid YAML syntax but invalid structure for unmarshaling
**Input**: `Load("/tmp/complex_config.yaml")`
**Expected Output**:
- Error returned
- Error contains "failed to unmarshal config"
- Config is nil

### 1.7 Special Characters and Edge Cases

#### TC1.7.1: Config values with special characters
**Test Data**:
```yaml
providers:
  anthropic:
    apiKey: "sk-\n\t\r"
    baseUrl: "https://api.example.com/path?query=value&other=123"
```
**Expected**: Special characters preserved

#### TC1.7.2: Very long config values
**Test Data**:
- APIKey: 10000+ character string
- SystemPrompt: 100000+ character string
**Expected**: Loads successfully

#### TC1.7.3: Unicode characters in config
**Test Data**:
```yaml
providers:
  anthropic:
    systemPrompt: "你好世界 🚀 Привет мир"
```
**Expected**: Unicode properly handled

#### TC1.7.4: Null/nil values in YAML
**Test Data**:
```yaml
providers:
  anthropic:
    baseUrl: null
    customHeaders: null
```
**Expected**: Fields handled as empty/nil

#### TC1.7.5: Duplicate keys in YAML
**Test Data**:
```yaml
currentProvider: anthropic
currentProvider: openai
```
**Expected**: YAML parser behavior (typically last value wins)

#### TC1.7.6: Large config file
**Setup**: Config with 100+ providers, 1000+ exclude patterns
**Expected**: Loads successfully, no memory issues

---

## 2. Save() Function Test Scenarios

### 2.1 Successful Save Operations

#### TC2.1.1: Save to explicit file path
**Setup**: Create Config struct with test data
**Input**: `config.Save("/tmp/test_save_config.yaml")`
**Expected Output**:
- No error returned
- File created at /tmp/test_save_config.yaml
- File is readable
- File contains valid YAML

**Verification**:
- File exists after call
- File permissions are 0644
- Parent directory has 0755 permissions
- Content can be parsed as YAML

#### TC2.1.2: Save to default location
**Setup**: Create Config struct
**Input**: `config.Save("")`
**Expected Output**:
- File created at ~/.config/how/config.yaml
- No error returned
- Content is valid YAML

#### TC2.1.3: Overwrite existing file
**Setup**: Config file already exists
**Input**: `config.Save("/tmp/existing_config.yaml")`
**Expected Output**:
- File overwritten
- New content replaces old content
- No error returned

#### TC2.1.4: Save newly created Config
**Setup**: `config := &Config{CurrentProvider: "anthropic"}`
**Expected**: Saves successfully with minimal data

#### TC2.1.5: Save modified Config
**Setup**: Load config, modify fields, save
**Expected**: All modifications persisted

### 2.2 Directory Creation

#### TC2.2.1: Create parent directory if missing
**Setup**: Parent directory doesn't exist
**Input**: `config.Save("/tmp/new_dir/config.yaml")`
**Expected Output**:
- /tmp/new_dir created
- File written successfully
- Directory has 0755 permissions

#### TC2.2.2: Create multiple missing parent directories
**Setup**: /tmp/a/b/c/d/ doesn't exist
**Input**: `config.Save("/tmp/a/b/c/d/config.yaml")`
**Expected Output**:
- All parent directories created
- Each with 0755 permissions
- File written successfully

#### TC2.2.3: Parent directory already exists
**Setup**: ~/.config/how already exists
**Input**: `config.Save("")`
**Expected Output**:
- No error during directory creation
- File written successfully

### 2.3 File Permissions

#### TC2.3.1: Verify created file permissions
**Setup**: Create and save config
**Expected**:
- File permissions are exactly 0644
- Readable by owner and others
- Writable only by owner

**Verification**:
```bash
stat -c '%A' file  # Should be -rw-r--r--
```

#### TC2.3.2: Verify created directory permissions
**Setup**: Save creates new directory
**Expected**:
- Directory permissions are exactly 0755
- Readable/executable by all
- Writable only by owner

#### TC2.3.3: File remains readable and writable
**Setup**: Save config, then read file
**Expected**: Can immediately read file back

### 2.4 YAML Output

#### TC2.4.1: Marshaling produces valid YAML
**Setup**: Config with various data types
**Expected Output**:
- Output is valid YAML syntax
- Can be parsed by YAML parser
- No data loss during marshaling

#### TC2.4.2: Round-trip consistency (Load → Save → Load)
**Setup**: Load config from file
**Step 1**: Save config to new file
**Step 2**: Load from new file
**Expected**: Loaded config matches original config

#### TC2.4.3: Empty/nil fields handling
**Setup**: Config with nil Providers, empty strings
**Expected**:
- Marshals correctly to YAML
- Fields with nil/empty values handled appropriately
- Can be unmarshaled back

#### TC2.4.4: Sensitive data is written (APIKey)
**Setup**: Config with populated APIKey
**Expected**: APIKey is written to file in plaintext
**Note**: Test that security implications are understood

#### TC2.4.5: Special characters preserved
**Setup**: Config with special chars (newlines, quotes, etc)
**Expected**: Special characters properly escaped in YAML

#### TC2.4.6: YAML formatting consistency
**Setup**: Save same config multiple times
**Expected**: Output formatting is consistent across saves

### 2.5 Complex Data Structures

#### TC2.5.1: Save with multiple providers
**Setup**: Config with 5 different providers
**Expected**: All providers written correctly, structure preserved

#### TC2.5.2: Save with large CustomHeaders map
**Setup**: ProviderConfig with 100+ custom headers
**Expected**: All headers written and readable

#### TC2.5.3: Save with long ExcludePatterns list
**Setup**: ContextConfig with 1000 exclude patterns
**Expected**: All patterns written correctly

#### TC2.5.4: Nested null/empty values
**Setup**: Config with null maps, empty lists, zero integers
**Expected**: All null/empty values handled correctly

### 2.6 Error Cases

#### TC2.6.1: getConfigDir() fails when saving to default location
**Setup**: Mock os.UserHomeDir() to fail
**Input**: `config.Save("")`
**Expected Output**:
- Error returned
- Error contains "failed to get config directory"
- File not created

#### TC2.6.2: Directory creation fails (permission denied)
**Setup**: Parent directory exists but no write permission
**Input**: `config.Save("/restricted/dir/config.yaml")`
**Expected Output**:
- Error returned
- Error indicates directory creation failed
- File not created

#### TC2.6.3: File write fails (disk full)
**Setup**: Mock os.WriteFile() to fail
**Expected Output**:
- Error returned
- Error contains "failed to write config file"
- Partial file may exist

#### TC2.6.4: Marshaling fails (circular reference)
**Setup**: Config with circular references (if possible)
**Expected**: Marshaling error handled

#### TC2.6.5: Invalid file path
**Setup**: `config.Save("\x00invalid")`
**Expected**: File write error returned

### 2.7 Edge Cases

#### TC2.7.1: Save with empty Providers map
**Setup**: `Config{CurrentProvider: "test", Providers: map[string]ProviderConfig{}}`
**Expected**: Saves successfully, empty map in YAML

#### TC2.7.2: Save with nil Providers map
**Setup**: `Config{CurrentProvider: "test", Providers: nil}`
**Expected**: Saves as null or empty in YAML

#### TC2.7.3: Save to same location multiple times
**Setup**: Save, modify, save again (same path)
**Expected**:
- Both saves succeed
- Final file contains latest data
- No conflicts or race conditions

#### TC2.7.4: Very large Config
**Setup**: Config with 1000 providers, large SystemPrompt values
**Expected**: Saves successfully without memory issues

#### TC2.7.5: FilePath with spaces
**Setup**: `config.Save("/tmp/path with spaces/config.yaml")`
**Expected**: Saves successfully

#### TC2.7.6: FilePath with special characters
**Setup**: `config.Save("/tmp/config@2024.yaml")`
**Expected**: Saves successfully if path is valid

---

## 3. getConfigDir() Function Test Scenarios

### 3.1 Successful Calls

#### TC3.1.1: Normal user with home directory
**Setup**: Normal user environment with HOME set
**Input**: `getConfigDir()`
**Expected Output**:
- Returns path like "/home/username/.config/how"
- Path is absolute
- No error

#### TC3.1.2: Root user
**Setup**: Running as root with HOME=/root
**Expected Output**:
- Returns "/root/.config/how"
- No error

#### TC3.1.3: Path format verification
**Input**: `getConfigDir()`
**Expected**:
- Ends with ".config/how"
- Contains home directory prefix
- Uses filepath.Join (OS-appropriate separators)

#### TC3.1.4: Multiple calls return same value
**Setup**: Call getConfigDir() multiple times
**Expected**: Returns consistent path

### 3.2 Error Cases

#### TC3.2.1: HOME environment variable not set
**Setup**: HOME env var unset
**Input**: `getConfigDir()`
**Expected Output**:
- Error returned from os.UserHomeDir()
- Function propagates error
- No path returned

#### TC3.2.2: os.UserHomeDir() returns error
**Setup**: Mock os.UserHomeDir() to return error
**Expected**:
- getConfigDir() returns empty string and error
- Error is propagated

#### TC3.2.3: Running in container without home
**Setup**: Container environment without proper home directory
**Expected**: os.UserHomeDir() error is returned

#### TC3.2.4: User is system user without home
**Setup**: System user (www-data, etc) without home directory
**Expected**: Error from os.UserHomeDir()

### 3.3 Path Edge Cases

#### TC3.3.1: Home directory with spaces
**Setup**: Home directory is "/home/user name"
**Expected Output**:
- Returns "/home/user name/.config/how"
- Spaces properly handled
- No error

#### TC3.3.2: Home directory with special characters
**Setup**: Home directory contains: ~, `, ', ", etc
**Expected**:
- Path constructed correctly
- os.UserHomeDir() returns sanitized path
- filepath.Join handles correctly

#### TC3.3.3: Symlinked home directory
**Setup**: $HOME points to symlink
**Expected**:
- os.UserHomeDir() resolves symlink
- Path returned correctly

#### TC3.3.4: Very long home directory path
**Setup**: Home directory with many nested levels
**Expected**: Path constructed correctly

#### TC3.3.5: Home directory with trailing slash
**Setup**: $HOME="/home/user/"
**Expected**: filepath.Join handles trailing slash correctly

---

## 4. Type Definition Tests

### 4.1 Config Struct

#### TC4.1.1: YAML tag mappings
**Test**: Create Config, marshal to YAML, verify field names
**Expected**:
- `CurrentProvider` maps to "currentProvider"
- `Providers` maps to "providers"
- `Context` maps to "context"
- `Display` maps to "display"
- `History` maps to "history"

#### TC4.1.2: Nested struct embedding
**Test**: Verify nested structs unmarshal correctly
**Expected**: All nested fields are accessible

#### TC4.1.3: Map type handling
**Test**: Marshal/unmarshal Config with Providers map
**Expected**: Map keys and values preserved

### 4.2 ProviderConfig Struct

#### TC4.2.1: All fields present and correct types
**Test**: Create ProviderConfig with all fields, verify marshaling
**Expected**: All 9 fields correctly marshaled

#### TC4.2.2: Optional fields (omitempty)
**Test**: Marshal ProviderConfig without optional fields
**Expected**:
- Required fields present in YAML
- Optional fields omitted if empty
- baseUrl, temperature, topP, systemPrompt, customHeaders affected

#### TC4.2.3: Float type handling (Temperature, TopP)
**Test**: Various float values
**Expected**: Decimal values correctly preserved

#### TC4.2.4: CustomHeaders map handling
**Test**: CustomHeaders with multiple entries
**Expected**: Map correctly marshaled/unmarshaled

### 4.3 ContextConfig Struct

#### TC4.3.1: Boolean and integer fields
**Test**: All combinations of boolean and integer values
**Expected**: Types correctly handled

#### TC4.3.2: ExcludePatterns array handling
**Test**: Array with various string values
**Expected**: Array correctly marshaled/unmarshaled

### 4.4 DisplayConfig Struct

#### TC4.4.1: All boolean combinations
**Test**: 16 combinations of 4 boolean fields
**Expected**: All combinations work (2^4 = 16)

### 4.5 HistoryConfig Struct

#### TC4.5.1: Mixed types (bool, int, string)
**Test**: All fields with various values
**Expected**: All types handled correctly

---

## 5. Integration Tests

### 5.1 Complete Workflow

#### TC5.1.1: Load → Modify → Save → Load
**Scenario**:
1. Load config from file
2. Modify CurrentProvider and add new provider
3. Save to same location
4. Load again
5. Verify modifications persisted

**Expected**:
- All modifications preserved
- No data loss or corruption
- Can be repeated multiple times

#### TC5.1.2: Load → Modify → Save to different location → Load
**Scenario**:
1. Load config from ~/config1.yaml
2. Modify config
3. Save to ~/config2.yaml
4. Load from ~/config2.yaml
5. Verify data matches

**Expected**: Config successfully copied with modifications

#### TC5.1.3: Environment overrides → Save → Reload
**Scenario**:
1. Load config file
2. Set environment variable (HOW_CURRENTPROVIDER=custom)
3. Load again (env override applied)
4. Save config
5. Load config
6. Verify saved config doesn't include override (value from file is saved)

**Expected**: Environment overrides are NOT persisted to file

#### TC5.1.4: Create empty config → Save → Load → Verify structure
**Scenario**:
1. Create Config{} (all zero values)
2. Save to file
3. Load from file
4. Verify empty Config structure preserved

**Expected**: Empty structures handled gracefully

### 5.2 Error Recovery

#### TC5.2.1: Recover from missing config
**Scenario**:
1. Load non-existent config (returns empty)
2. Programmatically set required values
3. Save configuration
4. Load and verify

**Expected**: Can build config from scratch despite missing file

#### TC5.2.2: Handle corrupted config gracefully
**Scenario**:
1. Load config with invalid values
2. Catch errors
3. Provide fallback/default values
4. Save corrected config

**Expected**: Errors are recoverable

---

## 6. Data Validation Edge Cases

### 6.1 ProviderConfig Validation

#### TC6.1.1: Temperature boundary values
**Test Cases**:
- Temperature: -0.1 (invalid)
- Temperature: 0.0 (valid, min)
- Temperature: 0.5 (valid, mid)
- Temperature: 1.0 (valid, max)
- Temperature: 1.1 (invalid)
**Note**: Current code doesn't validate; tests document this behavior

#### TC6.1.2: TopP boundary values
**Test Cases**: -0.1, 0.0, 0.5, 1.0, 1.1
**Note**: Current code doesn't validate

#### TC6.1.3: MaxTokens negative/zero
**Test Cases**: -1, 0, 1, 999999
**Note**: Current code doesn't validate negative values

#### TC6.1.4: Empty APIKey
**Setup**: `ProviderConfig{APIKey: ""}`
**Expected**: Loads but may cause issues at runtime (no validation)

### 6.2 ContextConfig Validation

#### TC6.2.1: Negative IncludeHistory
**Setup**: `IncludeHistory: -1`
**Expected**: Loads (no validation)

#### TC6.2.2: Negative MaxContextSize
**Setup**: `MaxContextSize: -1`
**Expected**: Loads (no validation)

#### TC6.2.3: ExcludePatterns with glob patterns
**Test Cases**:
- `*.log`
- `node_modules/*`
- `**/*.tmp`
- `[abc]*`
- `test?.go`
**Expected**: Patterns stored as-is (no validation)

### 6.3 HistoryConfig Validation

#### TC6.3.1: Zero MaxSize
**Setup**: `MaxSize: 0`
**Expected**: Loads (no limit validation)

#### TC6.3.2: FilePath variations
**Test Cases**:
- Absolute: `/home/user/.how/history`
- Relative: `history`
- With tilde: `~/.how/history`
- With spaces: `~/my documents/history`
- Very long: 1000+ character path
**Expected**: Paths stored as-is (no validation/expansion)

---

## Test Execution Summary

### Recommended Test Count by Category
| Category | Scenario Count | Total Tests |
|----------|---|---|
| Load() | 30 | 30 |
| Save() | 35 | 35 |
| getConfigDir() | 20 | 20 |
| Types | 15 | 15 |
| Integration | 10 | 10 |
| Data Validation | 20 | 20 |
| **TOTAL** | **~130** | **~130** |

### Estimated Coverage
- **Happy path**: 20-25 tests (basic load/save functionality)
- **Error cases**: 40-50 tests (various failure modes)
- **Edge cases**: 40-50 tests (boundary conditions, special data)
- **Integration**: 10-15 tests (multi-function workflows)

### Test File Organization
```
config_test.go
├── Exports test (uses exported API)
│   ├── TestLoad_*
│   ├── TestSave_*
│   └── TestIntegration_*
├── Private test (uses unexported functions with helper setup)
│   └── TestGetConfigDir_*
└── Type tests
    ├── TestConfig_YAML
    ├── TestProviderConfig_YAML
    └── ...
```

---

## Notes for Test Implementation

1. **Use t.TempDir()** for temporary test files (cleaned up automatically)
2. **Use testdata/** directory** for fixture files if needed
3. **Mock os.UserHomeDir()** for testing error cases without modifying environment
4. **Use table-driven tests** for multiple similar test cases
5. **Verify file operations** with actual file system checks (stat, etc)
6. **Test with real YAML files** in addition to generated test data
7. **Check error message content** not just error presence
8. **Test environment variables** carefully (isolate from test environment)
9. **Use assert libraries** (testify/assert) for cleaner test code
10. **Document why each test exists** in comments
