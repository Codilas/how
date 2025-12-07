# config.go Analysis: Methods and Test Coverage

## Project Context
- **Package**: `internal/config`
- **Language**: Go 1.22.4
- **Dependencies**: github.com/spf13/viper, gopkg.in/yaml.v3
- **Purpose**: Configuration loading and management from YAML files
- **Testing Framework**: Go's standard `testing` package

---

## Exported Types

### 1. **Config struct**
**Location**: `config.go:12-18`
```go
type Config struct {
    CurrentProvider string                    `yaml:"currentProvider"`
    Providers       map[string]ProviderConfig `yaml:"providers"`
    Context         ContextConfig             `yaml:"context"`
    Display         DisplayConfig             `yaml:"display"`
    History         HistoryConfig             `yaml:"history"`
}
```

**Fields**:
- `CurrentProvider`: string - Name of the currently active provider
- `Providers`: map of provider configurations indexed by name
- `Context`: nested configuration for context gathering
- `Display`: nested configuration for display settings
- `History`: nested configuration for history management

**Test Scenarios**:
- Verify struct tags for YAML marshaling/unmarshaling
- Test with valid provider configurations
- Test with empty providers map
- Test with missing optional fields

---

### 2. **ProviderConfig struct**
**Location**: `config.go:20-32`
```go
type ProviderConfig struct {
    Type           string            `yaml:"type"`
    APIKey         string            `yaml:"apiKey"`
    Model          string            `yaml:"model"`
    BaseURL        string            `yaml:"baseUrl,omitempty"`
    MaxTokens      int               `yaml:"maxTokens"`
    Temperature    float32           `yaml:"temperature,omitempty"`
    TopP           float32           `yaml:"topP,omitempty"`
    SystemPrompt   string            `yaml:"systemPrompt,omitempty"`
    CustomHeaders  map[string]string `yaml:"customHeaders,omitempty"`
}
```

**Fields**:
- `Type`: Provider type identifier (e.g., "anthropic")
- `APIKey`: Authentication key for the provider
- `Model`: Model identifier (e.g., "claude-3-opus-20240229")
- `BaseURL`: Optional custom API endpoint
- `MaxTokens`: Maximum tokens for responses
- `Temperature`: Optional sampling temperature (0.0-1.0)
- `TopP`: Optional nucleus sampling parameter
- `SystemPrompt`: Optional system message for prompts
- `CustomHeaders`: Optional custom HTTP headers

**Test Scenarios**:
- Minimal valid config (Type, APIKey, Model, MaxTokens)
- Config with all optional fields set
- Invalid temperature values (negative, > 1.0)
- Invalid TopP values (negative, > 1.0)
- Zero or negative MaxTokens
- Empty APIKey
- Empty Model
- Special characters in CustomHeaders keys/values
- Large CustomHeaders map

---

### 3. **ContextConfig struct**
**Location**: `config.go:34-41`
```go
type ContextConfig struct {
    IncludeFiles       bool     `yaml:"includeFiles"`
    IncludeHistory     int      `yaml:"includeHistory"`
    IncludeEnvironment bool     `yaml:"includeEnvironment"`
    IncludeGit         bool     `yaml:"includeGit"`
    MaxContextSize     int      `yaml:"maxContextSize"`
    ExcludePatterns    []string `yaml:"excludePatterns"`
}
```

**Fields**:
- `IncludeFiles`: Whether to include file context
- `IncludeHistory`: Number of history items to include
- `IncludeEnvironment`: Whether to include environment variables
- `IncludeGit`: Whether to include git information
- `MaxContextSize`: Maximum size of context in bytes
- `ExcludePatterns`: Glob patterns to exclude from context

**Test Scenarios**:
- All inclusion flags true
- All inclusion flags false
- Zero IncludeHistory
- Negative IncludeHistory values
- Zero MaxContextSize
- Negative MaxContextSize
- Empty ExcludePatterns
- ExcludePatterns with various glob patterns (*, ?, [])
- ExcludePatterns with special characters
- Large ExcludePatterns array

---

### 4. **DisplayConfig struct**
**Location**: `config.go:43-48`
```go
type DisplayConfig struct {
    SyntaxHighlight bool `yaml:"syntaxHighlight"`
    ShowContext     bool `yaml:"showContext"`
    Emoji           bool `yaml:"emoji"`
    Color           bool `yaml:"color"`
}
```

**Fields**:
- `SyntaxHighlight`: Enable syntax highlighting
- `ShowContext`: Show context in output
- `Emoji`: Enable emoji in output
- `Color`: Enable colored output

**Test Scenarios**:
- All flags true
- All flags false
- Mixed flag combinations (2^4 = 16 combinations total)
- Default values (zero values)

---

### 5. **HistoryConfig struct**
**Location**: `config.go:50-54`
```go
type HistoryConfig struct {
    Enabled  bool   `yaml:"enabled"`
    MaxSize  int    `yaml:"maxSize"`
    FilePath string `yaml:"filePath"`
}
```

**Fields**:
- `Enabled`: Whether history tracking is enabled
- `MaxSize`: Maximum number of history items
- `FilePath`: Path to history file

**Test Scenarios**:
- History enabled with valid FilePath
- History disabled with empty FilePath
- Zero or negative MaxSize
- FilePath with special characters
- FilePath with relative vs absolute paths
- FilePath with non-existent parent directories
- Very long FilePath values

---

## Exported Functions

### 1. **Load(configFile string) (*Config, error)**
**Location**: `config.go:56-92`

**Parameters**:
- `configFile`: Path to config file (empty string uses default location)

**Return Values**:
- `*Config`: Pointer to loaded configuration
- `error`: Error if loading or parsing fails

**Behavior**:
1. Creates a new Viper instance
2. If `configFile` is provided, uses it directly; otherwise looks in default config directory
3. Reads YAML configuration file
4. If file not found, issues warning but continues
5. Sets up environment variable overrides with "HOW" prefix
6. Unmarshal YAML into Config struct
7. Returns pointer to Config or error

**Error Cases**:
- `getConfigDir()` fails
- Config file reading fails (except ConfigFileNotFoundError)
- YAML unmarshaling fails

**Edge Cases**:
- Empty configFile path
- configFile path doesn't exist (returns warning, empty config)
- Invalid YAML syntax
- Missing required fields
- Extra unknown fields
- Environment variables override config values

**Test Scenarios**:
1. **Valid Config Loading**:
   - Load from explicit file path
   - Load from default location (~/.config/how/config.yaml)
   - Load with all fields populated
   - Load with minimal required fields
   - Load with only empty providers

2. **Missing File Handling**:
   - Config file doesn't exist (returns warning, empty config)
   - Empty config file (valid YAML but no content)
   - Parent directory doesn't exist

3. **Invalid YAML**:
   - Malformed YAML syntax
   - Incorrect indentation
   - Invalid data types for fields (string for int field, etc.)
   - Invalid float values for Temperature/TopP

4. **Environment Variable Overrides**:
   - HOW_CURRENTPROVIDER env var
   - HOW_PROVIDERS_* nested env vars
   - HOW_CONTEXT_* nested env vars
   - Env vars override file values

5. **Error Conditions**:
   - getConfigDir() fails (home directory not available)
   - Config file exists but not readable (permission denied)
   - YAML unmarshaling fails

6. **Edge Cases**:
   - Empty Providers map
   - Null/nil values in YAML
   - Duplicate keys in YAML
   - Very large config file
   - Config with only CurrentProvider
   - Config with circular references (if applicable)

---

### 2. **Save(configFile string) error**
**Location**: `config.go:94-120`

**Receiver**: `*Config` (method on Config struct)

**Parameters**:
- `configFile`: Path where config should be saved (empty string uses default location)

**Return Values**:
- `error`: Error if save fails, nil on success

**Behavior**:
1. If `configFile` is empty, determines default location
2. Creates directory structure if it doesn't exist (with 0755 permissions)
3. Marshals Config struct to YAML
4. Writes YAML to file with 0644 permissions
5. Returns error or nil

**Error Cases**:
- `getConfigDir()` fails
- Directory creation fails
- YAML marshaling fails
- File write fails

**Permissions**:
- Directory: 0755 (rwxr-xr-x)
- File: 0644 (rw-r--r--)

**Test Scenarios**:
1. **Valid Save Operations**:
   - Save to explicit file path
   - Save to default location
   - Save newly created Config
   - Save modified Config
   - Overwrite existing config file

2. **Directory Creation**:
   - Parent directory doesn't exist (creates it)
   - Multiple levels of missing directories
   - Directory already exists

3. **File Permissions**:
   - Verify created file has 0644 permissions
   - Verify created directory has 0755 permissions
   - Verify file is readable after write
   - Verify file is writable after write

4. **YAML Output**:
   - Verify YAML marshaling produces valid YAML
   - Verify can load saved config back
   - Verify sensitive data is written (APIKey, etc.)
   - Verify null/empty fields are handled correctly
   - Verify formatting is consistent

5. **Error Conditions**:
   - getConfigDir() fails
   - Parent directory can't be created (permission denied)
   - File write permission denied
   - Disk full
   - Invalid file path
   - Path contains non-existent parent directories (with restricted perms)

6. **Edge Cases**:
   - Config with empty Providers map
   - Config with nil fields
   - Saving to same location multiple times
   - Saving with very large config
   - Empty ConfigFile path when home directory unavailable

---

## Private Functions

### 1. **getConfigDir() (string, error)**
**Location**: `config.go:122-128`

**Parameters**: None

**Return Values**:
- `string`: Path to config directory (~/.config/how)
- `error`: Error if unable to determine home directory

**Behavior**:
1. Calls `os.UserHomeDir()` to get home directory
2. Appends `.config/how` to home directory
3. Returns full path or error

**Error Cases**:
- `os.UserHomeDir()` fails (e.g., $HOME not set, in container without home)

**Depends On**:
- Go standard library `os.UserHomeDir()`

**Test Scenarios**:
1. **Successful Calls**:
   - Normal user home directory exists
   - Verify returned path matches expected pattern
   - Verify path is absolute
   - Verify path ends with `.config/how`

2. **Error Handling**:
   - HOME environment variable not set
   - os.UserHomeDir() returns error
   - User is running as system user without home directory

3. **Edge Cases**:
   - Home directory path contains spaces
   - Home directory path contains special characters
   - Symlinked home directory
   - Home directory is relative path (shouldn't happen but test defensive)

---

## Complete Test Coverage Checklist

### Type Definition Tests
- [ ] Config struct YAML tag mappings
- [ ] ProviderConfig struct YAML tag mappings
- [ ] ContextConfig struct YAML tag mappings
- [ ] DisplayConfig struct YAML tag mappings
- [ ] HistoryConfig struct YAML tag mappings
- [ ] Field type validations
- [ ] Embedded struct handling

### Load() Function Tests
- [ ] Load from explicit file path
- [ ] Load from default location
- [ ] Load with complete configuration
- [ ] Load with minimal configuration
- [ ] Missing config file (warning + empty config)
- [ ] Invalid YAML syntax
- [ ] Invalid field data types
- [ ] Environment variable overrides
- [ ] Multiple environment variable overrides
- [ ] Nested field environment overrides
- [ ] getConfigDir() error propagation
- [ ] Config file read error propagation
- [ ] Unmarshal error handling

### Save() Function Tests
- [ ] Save to explicit file path
- [ ] Save to default location
- [ ] Create parent directories
- [ ] Overwrite existing file
- [ ] Verify file permissions (0644)
- [ ] Verify directory permissions (0755)
- [ ] YAML marshaling output validity
- [ ] Round-trip save/load consistency
- [ ] Save with empty config
- [ ] Save with null fields
- [ ] Save with special characters in values
- [ ] getConfigDir() error propagation
- [ ] Directory creation error handling
- [ ] File write error handling

### getConfigDir() Function Tests
- [ ] Successful home directory resolution
- [ ] Path format verification
- [ ] HOME environment variable not set error
- [ ] os.UserHomeDir() error handling
- [ ] Path with spaces
- [ ] Path with special characters

### Integration Tests
- [ ] Load default config → Modify → Save → Load again
- [ ] Load from file → Save to different location → Load from new location
- [ ] Environment variables override file config → Save → Verify saved doesn't include overrides
- [ ] Create empty config → Save → Load → Verify empty structure

### Data Validation Edge Cases
- [ ] Temperature values: -0.1, 0.0, 0.5, 1.0, 1.1
- [ ] TopP values: -0.1, 0.0, 0.5, 1.0, 1.1
- [ ] MaxTokens: -1, 0, 1, very large number
- [ ] IncludeHistory: -1, 0, 1, very large number
- [ ] MaxContextSize: -1, 0, 1, very large number
- [ ] ExcludePatterns with glob patterns
- [ ] CustomHeaders with special characters
- [ ] FilePath with relative and absolute paths

---

## Test File Structure

### Recommended file: `config_test.go`

Suggested test function organization:
```
TestLoad_ValidConfigFile
TestLoad_MissingConfigFile
TestLoad_InvalidYAML
TestLoad_EnvironmentVariableOverrides
TestLoad_ErrorCases

TestSave_ExplicitPath
TestSave_DefaultLocation
TestSave_FilePermissions
TestSave_RoundTrip
TestSave_ErrorCases

TestGetConfigDir_Success
TestGetConfigDir_NoHomeDir
TestGetConfigDir_Errors

TestConfigStruct_YAML
TestProviderConfig_YAML
TestContextConfig_YAML
TestDisplayConfig_YAML
TestHistoryConfig_YAML

TestIntegration_LoadModifySaveLoad
TestIntegration_LoadFromDifferentLocation
```

---

## Summary

**Total Methods Analyzed**: 5
- Exported functions: 2 (Load, Save)
- Private functions: 1 (getConfigDir)
- Exported types: 5 (Config, ProviderConfig, ContextConfig, DisplayConfig, HistoryConfig)

**Key Edge Cases**:
- Missing configuration files and directory setup
- YAML marshaling/unmarshaling with various data types
- Environment variable overrides
- File system permissions and path handling
- Error propagation from nested function calls

**Recommended Test Count**: 50-70 test cases for comprehensive coverage
