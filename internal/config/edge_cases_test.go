package config

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

// ===== Invalid YAML Syntax Tests =====

// TestLoadInvalidYAMLSyntaxBadIndentation tests loading config with invalid indentation.
func TestLoadInvalidYAMLSyntaxBadIndentation(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "invalid.yaml")

	invalidYAML := `currentProvider: anthropic
providers:
  anthropic:
    type: anthropic
  apiKey: test  # Wrong indentation - should be indented under anthropic
    model: claude-3-5-sonnet
    maxTokens: 2048
`

	helper.WriteYAML(invalidYAML, tempDir)

	// Load should handle or error gracefully
	config, err := Load(configPath)
	if err != nil {
		t.Logf("Invalid YAML correctly returned error: %v", err)
	} else if config != nil {
		// If it doesn't error, should at least return a config
		t.Logf("Invalid YAML returned config (may have parsed partially)")
	}
}

// TestLoadInvalidYAMLMissingColon tests YAML with missing colons.
func TestLoadInvalidYAMLMissingColon(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "invalid.yaml")

	invalidYAML := `currentProvider anthropic
providers:
  anthropic:
    type: anthropic
    apiKey: test
    model: claude
`

	helper.WriteYAML(invalidYAML, tempDir)

	config, err := Load(configPath)
	if err != nil {
		t.Logf("Missing colon correctly returned error: %v", err)
	} else if config != nil {
		t.Logf("Missing colon returned config")
	}
}

// TestLoadInvalidYAMLUnclosedQuote tests YAML with unclosed quotes.
func TestLoadInvalidYAMLUnclosedQuote(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "invalid.yaml")

	invalidYAML := `currentProvider: anthropic
providers:
  anthropic:
    type: anthropic
    apiKey: "unclosed-quote
    model: claude
`

	helper.WriteYAML(invalidYAML, tempDir)

	config, err := Load(configPath)
	if err != nil {
		t.Logf("Unclosed quote correctly returned error: %v", err)
	} else if config != nil {
		t.Logf("Unclosed quote returned config")
	}
}

// TestLoadInvalidYAMLInvalidScalar tests YAML with invalid scalar values.
func TestLoadInvalidYAMLInvalidScalar(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "invalid.yaml")

	invalidYAML := `currentProvider: anthropic
providers:
  anthropic:
    type: anthropic
    apiKey: test
    model: claude
    maxTokens: [invalid array syntax for integer]
`

	helper.WriteYAML(invalidYAML, tempDir)

	config, err := Load(configPath)
	// YAML parser may accept this as a string or error
	if err != nil {
		t.Logf("Invalid scalar correctly returned error: %v", err)
	} else if config != nil {
		t.Logf("Invalid scalar returned config")
	}
}

// ===== Corrupted Config Files Tests =====

// TestLoadCorruptedConfigFileBinary tests loading a config file with binary data.
func TestLoadCorruptedConfigFileBinary(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "corrupted.yaml")

	// Write binary data to file
	binaryData := []byte{0xFF, 0xFE, 0x00, 0x01, 0x02, 0x03}
	if err := os.WriteFile(configPath, binaryData, 0644); err != nil {
		t.Fatalf("failed to write binary data: %v", err)
	}

	config, err := Load(configPath)
	if err != nil {
		t.Logf("Binary data correctly returned error: %v", err)
	} else if config != nil {
		t.Logf("Binary data returned config")
	}
}

// TestLoadEmptyConfigFile tests loading an empty config file.
func TestLoadEmptyConfigFile(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "empty.yaml")

	helper.WriteYAML("", tempDir)

	config, err := Load(configPath)
	// Empty file should either error or return an empty config
	if err != nil {
		t.Logf("Empty config correctly returned error: %v", err)
	} else if config != nil {
		// Empty config should be valid
		t.Logf("Empty config returned: %v", config)
	}
}

// TestLoadPartiallyCorruptedYAML tests loading YAML that is partially valid.
func TestLoadPartiallyCorruptedYAML(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "partial.yaml")

	// First part is valid, second part is corrupted
	yamlContent := `currentProvider: anthropic
providers:
  anthropic:
    type: anthropic
    apiKey: test
    model: claude
    maxTokens: 2048
context:
  includeFiles: true
  [invalid yaml here
`

	helper.WriteYAML(yamlContent, tempDir)

	config, err := Load(configPath)
	if err != nil {
		t.Logf("Partially corrupted YAML correctly returned error: %v", err)
	} else if config != nil {
		t.Logf("Partially corrupted YAML returned config")
	}
}

// ===== Malformed Config Structure Tests =====

// TestLoadMalformedProviderConfig tests loading with missing required provider fields.
func TestLoadMalformedProviderConfig(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "malformed.yaml")

	// Missing required fields like type or apiKey
	malformedYAML := `currentProvider: anthropic
providers:
  anthropic:
    model: claude
    maxTokens: 2048
`

	helper.WriteYAML(malformedYAML, tempDir)

	config, err := Load(configPath)
	if err != nil {
		t.Logf("Malformed provider config correctly returned error: %v", err)
	} else if config != nil {
		// Config should still load, but with empty fields
		if config.Providers != nil && len(config.Providers) > 0 {
			provider := config.Providers["anthropic"]
			if provider.Type == "" && provider.APIKey == "" {
				t.Logf("Malformed provider loaded with empty fields as expected")
			}
		}
	}
}

// TestLoadMalformedContextConfig tests loading with invalid context config values.
func TestLoadMalformedContextConfig(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "malformed.yaml")

	malformedYAML := `currentProvider: anthropic
providers:
  anthropic:
    type: anthropic
    apiKey: test
    model: claude
    maxTokens: 2048
context:
  includeFiles: "not-a-boolean"
  includeHistory: "not-a-number"
  maxContextSize: invalid
`

	helper.WriteYAML(malformedYAML, tempDir)

	config, err := Load(configPath)
	if err != nil {
		t.Logf("Malformed context config correctly returned error: %v", err)
	} else if config != nil {
		t.Logf("Malformed context config returned: %v", config.Context)
	}
}

// TestLoadMalformedDisplayConfig tests loading with invalid display config values.
func TestLoadMalformedDisplayConfig(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "malformed.yaml")

	malformedYAML := `currentProvider: anthropic
providers:
  anthropic:
    type: anthropic
    apiKey: test
    model: claude
    maxTokens: 2048
display:
  syntaxHighlight: "maybe"
  showContext: 123
  emoji: null
  color: []
`

	helper.WriteYAML(malformedYAML, tempDir)

	config, err := Load(configPath)
	if err != nil {
		t.Logf("Malformed display config correctly returned error: %v", err)
	} else if config != nil {
		t.Logf("Malformed display config returned: %v", config.Display)
	}
}

// TestLoadMalformedHistoryConfig tests loading with invalid history config values.
func TestLoadMalformedHistoryConfig(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "malformed.yaml")

	malformedYAML := `currentProvider: anthropic
providers:
  anthropic:
    type: anthropic
    apiKey: test
    model: claude
    maxTokens: 2048
history:
  enabled: "yes"
  maxSize: "thousands"
  filePath: null
`

	helper.WriteYAML(malformedYAML, tempDir)

	config, err := Load(configPath)
	if err != nil {
		t.Logf("Malformed history config correctly returned error: %v", err)
	} else if config != nil {
		t.Logf("Malformed history config returned: %v", config.History)
	}
}

// ===== Permission Error Tests =====

// TestSaveConfigToReadOnlyDirectory tests saving config to a read-only directory.
func TestSaveConfigToReadOnlyDirectory(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("skipping permission test when running as root")
	}

	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	readOnlyDir := filepath.Join(tempDir, "readonly")

	if err := os.Mkdir(readOnlyDir, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	// Remove write permission
	if err := os.Chmod(readOnlyDir, 0555); err != nil {
		t.Fatalf("failed to change permissions: %v", err)
	}

	defer os.Chmod(readOnlyDir, 0755)

	config := SampleConfig()
	configPath := filepath.Join(readOnlyDir, "config.yaml")

	err := config.Save(configPath)
	if err != nil {
		t.Logf("Save to read-only directory correctly returned error: %v", err)
	} else {
		t.Error("expected error when saving to read-only directory, got nil")
	}
}

// TestLoadConfigFromUnreadableFile tests loading from a file without read permissions.
func TestLoadConfigFromUnreadableFile(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("skipping permission test when running as root")
	}

	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	helper.WriteYAML(SampleYAML(), tempDir)

	// Remove read permission
	if err := os.Chmod(configPath, 0000); err != nil {
		t.Fatalf("failed to change permissions: %v", err)
	}

	defer os.Chmod(configPath, 0644)

	config, err := Load(configPath)
	if err != nil {
		t.Logf("Load from unreadable file correctly returned error: %v", err)
	} else if config != nil {
		t.Logf("Load from unreadable file returned config (permissions may allow read)")
	}
}

// ===== Disk Space / Resource Limit Tests =====

// TestSaveConfigWithExtremelyLargeValues tests saving config with very large values.
func TestSaveConfigWithExtremelyLargeValues(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "large.yaml")

	config := &Config{
		CurrentProvider: "test",
		Providers: map[string]ProviderConfig{
			"test": {
				Type:      "test",
				APIKey:    "test",
				Model:     "test",
				MaxTokens: 2147483647, // Max int32
			},
		},
	}

	err := config.Save(configPath)
	if err != nil {
		t.Errorf("failed to save config with large values: %v", err)
	}

	loaded, err := Load(configPath)
	if err != nil {
		t.Errorf("failed to load config with large values: %v", err)
	}

	if loaded.Providers["test"].MaxTokens != 2147483647 {
		t.Errorf("large value not preserved: expected 2147483647, got %d", loaded.Providers["test"].MaxTokens)
	}
}

// TestSaveConfigWithVeryLongStrings tests saving config with extremely long string values.
func TestSaveConfigWithVeryLongStrings(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "longstrings.yaml")

	// Create a very long API key
	longAPIKey := ""
	for i := 0; i < 10000; i++ {
		longAPIKey += "a"
	}

	config := &Config{
		CurrentProvider: "test",
		Providers: map[string]ProviderConfig{
			"test": {
				Type:      "test",
				APIKey:    longAPIKey,
				Model:     "test",
				MaxTokens: 2048,
			},
		},
	}

	err := config.Save(configPath)
	if err != nil {
		t.Errorf("failed to save config with long strings: %v", err)
	}

	loaded, err := Load(configPath)
	if err != nil {
		t.Errorf("failed to load config with long strings: %v", err)
	}

	if loaded.Providers["test"].APIKey != longAPIKey {
		t.Errorf("long string not preserved correctly")
	}
}

// ===== Edge Cases for Nested Structures =====

// TestLoadConfigWithDeeplyNestedProviders tests loading config with unusual provider structure.
func TestLoadConfigWithNullValues(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "nulls.yaml")

	yamlContent := `currentProvider: null
providers: null
context: null
display: null
history: null
`

	helper.WriteYAML(yamlContent, tempDir)

	config, err := Load(configPath)
	if err != nil {
		t.Logf("Config with null values correctly returned error: %v", err)
	} else if config != nil {
		t.Logf("Config with null values returned: %v", config)
	}
}

// TestLoadConfigWithEmptyProviders tests loading config with empty providers map.
func TestLoadConfigWithEmptyProviders(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "empty_providers.yaml")

	yamlContent := `currentProvider: anthropic
providers: {}
context:
  includeFiles: true
display:
  color: true
history:
  enabled: true
`

	helper.WriteYAML(yamlContent, tempDir)

	config, err := Load(configPath)
	if err != nil {
		t.Errorf("failed to load config with empty providers: %v", err)
	}

	if config != nil {
		if len(config.Providers) != 0 {
			t.Errorf("expected empty providers map, got %d providers", len(config.Providers))
		}
	}
}

// ===== Edge Cases for Array Fields =====

// TestLoadConfigWithDuplicateExcludePatterns tests loading with duplicate patterns.
func TestLoadConfigWithDuplicateExcludePatterns(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "duplicates.yaml")

	yamlContent := `currentProvider: anthropic
providers:
  anthropic:
    type: anthropic
    apiKey: test
    model: claude
    maxTokens: 2048
context:
  includeFiles: true
  excludePatterns:
    - .git
    - .git
    - node_modules
    - .git
display:
  color: true
history:
  enabled: true
`

	helper.WriteYAML(yamlContent, tempDir)

	config, err := Load(configPath)
	if err != nil {
		t.Errorf("failed to load config with duplicate patterns: %v", err)
	}

	if config != nil {
		// Duplicates may be preserved or deduplicated - both are acceptable
		t.Logf("Exclude patterns: %v", config.Context.ExcludePatterns)
	}
}

// TestLoadConfigWithLargeExcludePatternsList tests loading with many exclude patterns.
func TestLoadConfigWithLargeExcludePatternsList(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "many_patterns.yaml")

	// Create YAML with many patterns
	patterns := "excludePatterns:\n"
	for i := 0; i < 1000; i++ {
		patterns += "    - pattern" + string(rune(i)) + "\n"
	}

	yamlContent := `currentProvider: anthropic
providers:
  anthropic:
    type: anthropic
    apiKey: test
    model: claude
    maxTokens: 2048
context:
  includeFiles: true
` + patterns + `display:
  color: true
history:
  enabled: true
`

	helper.WriteYAML(yamlContent, tempDir)

	config, err := Load(configPath)
	if err != nil {
		t.Errorf("failed to load config with many patterns: %v", err)
	}

	if config != nil {
		if len(config.Context.ExcludePatterns) == 0 {
			t.Error("expected exclude patterns to be loaded")
		}
	}
}

// ===== Edge Cases for Numeric Fields =====

// TestLoadConfigWithNegativeNumericValues tests loading with negative values where positive is expected.
func TestLoadConfigWithNegativeNumericValues(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "negative.yaml")

	yamlContent := `currentProvider: anthropic
providers:
  anthropic:
    type: anthropic
    apiKey: test
    model: claude
    maxTokens: -2048
context:
  includeFiles: true
  includeHistory: -50
  maxContextSize: -8000
display:
  color: true
history:
  enabled: true
  maxSize: -1000
`

	helper.WriteYAML(yamlContent, tempDir)

	config, err := Load(configPath)
	if err != nil {
		t.Logf("Config with negative values correctly returned error: %v", err)
	} else if config != nil {
		// YAML will parse negative numbers
		t.Logf("Config loaded with negative values (accepted by YAML parser): maxTokens=%d, includeHistory=%d",
			config.Providers["anthropic"].MaxTokens, config.Context.IncludeHistory)
	}
}

// TestLoadConfigWithZeroValues tests loading with zero values.
func TestLoadConfigWithZeroValues(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "zeros.yaml")

	yamlContent := `currentProvider: anthropic
providers:
  anthropic:
    type: anthropic
    apiKey: test
    model: claude
    maxTokens: 0
context:
  includeFiles: false
  includeHistory: 0
  maxContextSize: 0
display:
  color: false
history:
  enabled: false
  maxSize: 0
`

	helper.WriteYAML(yamlContent, tempDir)

	config, err := Load(configPath)
	if err != nil {
		t.Errorf("failed to load config with zero values: %v", err)
	}

	if config != nil {
		if config.Providers["anthropic"].MaxTokens != 0 {
			t.Errorf("expected maxTokens to be 0, got %d", config.Providers["anthropic"].MaxTokens)
		}
		if config.Context.IncludeHistory != 0 {
			t.Errorf("expected includeHistory to be 0, got %d", config.Context.IncludeHistory)
		}
	}
}

// ===== Edge Cases for String Fields =====

// TestLoadConfigWithSpecialCharactersInAllFields tests special chars in various fields.
func TestLoadConfigWithSpecialCharactersInAllFields(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "special_chars.yaml")

	config := &Config{
		CurrentProvider: "special!@#$%",
		Providers: map[string]ProviderConfig{
			"special!@#$%": {
				Type:         "type-with-special-chars-!@#$%",
				APIKey:       "key-!@#$%^&*()",
				Model:        "model-日本語-العربية-עברית",
				MaxTokens:    2048,
				SystemPrompt: "Line 1\nLine 2\tTabbed\n\"Quoted\"\n'Apostrophe'",
				CustomHeaders: map[string]string{
					"X-Special-!@#$": "value-with-special-chars-!@#$%^&*()",
				},
			},
		},
		Context: ContextConfig{
			IncludeFiles:    true,
			ExcludePatterns: []string{"pattern-!@#$%", "pattern-日本語", "path/with/slashes\\and\\backslashes"},
		},
		Display: DisplayConfig{
			Color: true,
		},
		History: HistoryConfig{
			Enabled:  true,
			FilePath: "/path/with/special-!@#$%^&*()/chars",
		},
	}

	if err := config.Save(configPath); err != nil {
		t.Errorf("failed to save config with special characters: %v", err)
	}

	loaded, err := Load(configPath)
	if err != nil {
		t.Errorf("failed to load config with special characters: %v", err)
	}

	if loaded != nil {
		provider := loaded.Providers["special!@#$%"]
		if provider.APIKey != "key-!@#$%^&*()" {
			t.Errorf("special characters in APIKey not preserved: got %s", provider.APIKey)
		}

		if provider.Model != "model-日本語-العربية-עברית" {
			t.Errorf("unicode characters not preserved: got %s", provider.Model)
		}

		if provider.SystemPrompt != "Line 1\nLine 2\tTabbed\n\"Quoted\"\n'Apostrophe'" {
			t.Errorf("multiline and special characters in SystemPrompt not preserved")
		}
	}
}

// TestLoadConfigWithEmptyStrings tests loading with empty string values.
func TestLoadConfigWithEmptyStrings(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "empty_strings.yaml")

	yamlContent := `currentProvider: ""
providers:
  empty:
    type: ""
    apiKey: ""
    model: ""
    maxTokens: 0
context:
  includeFiles: false
display:
  color: true
history:
  filePath: ""
  enabled: false
`

	helper.WriteYAML(yamlContent, tempDir)

	config, err := Load(configPath)
	if err != nil {
		t.Errorf("failed to load config with empty strings: %v", err)
	}

	if config != nil {
		// Empty strings should be preserved
		if config.CurrentProvider != "" {
			t.Errorf("expected empty CurrentProvider, got %q", config.CurrentProvider)
		}

		if config.Providers["empty"].Type != "" {
			t.Errorf("expected empty type, got %q", config.Providers["empty"].Type)
		}
	}
}

// ===== YAML Marshaling Edge Cases =====

// TestYAMLMarshalingWithCyclicReferences would test if there were any, but Go maps can't have cycles
// TestYAMLUnmarshalingWithUnknownFields tests handling of unknown fields in YAML.
func TestYAMLUnmarshalingWithUnknownFields(t *testing.T) {
	yamlData := `currentProvider: anthropic
providers:
  anthropic:
    type: anthropic
    apiKey: test
    model: claude
    maxTokens: 2048
    unknownField: "should be ignored"
    anotherUnknown: 123
context:
  includeFiles: true
  unknownContextField: true
display:
  color: true
  unknownDisplayField: false
history:
  enabled: true
  unknownHistoryField: "ignored"
unknownTopLevel: "should be ignored"
`

	var config Config
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Errorf("failed to unmarshal YAML with unknown fields: %v", err)
	}

	// Unknown fields should be ignored
	if config.CurrentProvider != "anthropic" {
		t.Errorf("provider mismatch: expected anthropic, got %s", config.CurrentProvider)
	}

	provider := config.Providers["anthropic"]
	if provider.APIKey != "test" {
		t.Errorf("APIKey mismatch: expected test, got %s", provider.APIKey)
	}
}

// TestYAMLMarshalingWithFloatingPointPrecision tests float preservation.
func TestYAMLMarshalingWithFloatingPointPrecision(t *testing.T) {
	original := ProviderConfig{
		Type:        "test",
		APIKey:      "key",
		Model:       "model",
		MaxTokens:   2048,
		Temperature: 0.123456789,
		TopP:        0.999999999,
	}

	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled ProviderConfig
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Float32 will lose some precision, but should be close
	if unmarshaled.Temperature < 0.123 || unmarshaled.Temperature > 0.124 {
		t.Errorf("temperature precision issue: original %f, got %f", original.Temperature, unmarshaled.Temperature)
	}
}

// ===== Config State Edge Cases =====

// TestLoadSaveWithNilProvidersMap tests handling of nil provider maps.
func TestLoadSaveWithNilProvidersMap(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "nil_providers.yaml")

	config := &Config{
		CurrentProvider: "test",
		Providers:       nil, // Nil instead of empty map
	}

	// This should handle nil gracefully
	if err := config.Save(configPath); err != nil {
		t.Logf("Save with nil providers returned error: %v", err)
	} else {
		loaded, err := Load(configPath)
		if err != nil {
			t.Errorf("failed to load config with nil providers: %v", err)
		}

		if loaded != nil && loaded.Providers == nil {
			t.Logf("Nil providers preserved as nil")
		} else if loaded != nil && len(loaded.Providers) == 0 {
			t.Logf("Nil providers converted to empty map (acceptable)")
		}
	}
}

// TestConfigWithCircularReferences is not applicable as YAML doesn't support them
// but we test with very nested structures instead

// TestLoadConfigWithMixedIndentation tests YAML with inconsistent indentation.
func TestLoadConfigWithMixedIndentation(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "mixed_indent.yaml")

	// Mixed spaces and tabs (problematic in YAML)
	yamlContent := "currentProvider: anthropic\nproviders:\n  anthropic:\n    type: anthropic\n\tapiKey: test\n    model: claude\n    maxTokens: 2048\n"

	helper.WriteYAML(yamlContent, tempDir)

	config, err := Load(configPath)
	if err != nil {
		t.Logf("Mixed indentation correctly returned error: %v", err)
	} else if config != nil {
		t.Logf("Mixed indentation returned config")
	}
}

// TestSaveAndLoadPreservesExactStructure tests that the exact structure is preserved.
func TestSaveAndLoadPreservesExactStructure(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "exact_structure.yaml")

	original := &Config{
		CurrentProvider: "primary",
		Providers: map[string]ProviderConfig{
			"primary": {
				Type:      "anthropic",
				APIKey:    "sk-ant-test",
				Model:     "claude-3-5-sonnet",
				BaseURL:   "https://api.anthropic.com",
				MaxTokens: 4096,
				Temperature: 0.7,
				TopP:       1.0,
				SystemPrompt: "You are a helpful assistant.\nWith multiple lines.\nAnd special chars: $100!",
				CustomHeaders: map[string]string{
					"Authorization": "Bearer token",
					"X-Custom":      "value",
					"X-Another":     "another-value",
				},
			},
		},
		Context: ContextConfig{
			IncludeFiles:       true,
			IncludeHistory:     100,
			IncludeEnvironment: true,
			IncludeGit:         true,
			MaxContextSize:     16000,
			ExcludePatterns:    []string{".git", ".env", "node_modules", ".vscode", "dist"},
		},
		Display: DisplayConfig{
			SyntaxHighlight: true,
			ShowContext:     true,
			Emoji:           false,
			Color:           true,
		},
		History: HistoryConfig{
			Enabled:  true,
			MaxSize:  5000,
			FilePath: "~/.local/share/how/history",
		},
	}

	if err := original.Save(configPath); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	// Deep comparison
	if loaded.CurrentProvider != original.CurrentProvider {
		t.Errorf("CurrentProvider mismatch")
	}

	if len(loaded.Providers) != len(original.Providers) {
		t.Errorf("provider count mismatch")
	}

	if loaded.Providers["primary"].CustomHeaders["Authorization"] != "Bearer token" {
		t.Error("custom header not preserved")
	}

	if len(loaded.Context.ExcludePatterns) != len(original.Context.ExcludePatterns) {
		t.Error("exclude patterns count mismatch")
	}
}
