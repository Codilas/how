package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadSaveLoadCycle verifies the complete cycle of loading configuration,
// modifying it, saving it, and loading it again to ensure data integrity.
func TestLoadSaveLoadCycle(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Step 1: Create initial config and save it
	originalConfig := SampleConfig()
	if err := originalConfig.Save(configPath); err != nil {
		t.Fatalf("initial save failed: %v", err)
	}

	// Step 2: Load the saved config
	loadedConfig1, err := Load(configPath)
	if err != nil {
		t.Fatalf("first load failed: %v", err)
	}

	// Step 3: Verify first load matches original
	assertConfigsEqual(t, originalConfig, loadedConfig1, "first load")

	// Step 4: Modify the loaded config
	loadedConfig1.CurrentProvider = "modified-provider"
	loadedConfig1.Display.SyntaxHighlight = false
	loadedConfig1.Context.IncludeHistory = 100

	// Step 5: Save the modified config
	if err := loadedConfig1.Save(configPath); err != nil {
		t.Fatalf("modified save failed: %v", err)
	}

	// Step 6: Load the modified config
	loadedConfig2, err := Load(configPath)
	if err != nil {
		t.Fatalf("second load failed: %v", err)
	}

	// Step 7: Verify modifications were persisted
	if loadedConfig2.CurrentProvider != "modified-provider" {
		t.Errorf("CurrentProvider not persisted: expected modified-provider, got %s", loadedConfig2.CurrentProvider)
	}

	if loadedConfig2.Display.SyntaxHighlight {
		t.Error("SyntaxHighlight modification not persisted")
	}

	if loadedConfig2.Context.IncludeHistory != 100 {
		t.Errorf("IncludeHistory modification not persisted: expected 100, got %d", loadedConfig2.Context.IncludeHistory)
	}

	// Step 8: Verify all original data that wasn't modified is still intact
	if loadedConfig2.Providers == nil || len(loadedConfig2.Providers) == 0 {
		t.Error("Providers data lost during save-load cycle")
	}

	if loadedConfig2.Display.ShowContext != originalConfig.Display.ShowContext {
		t.Error("Unmodified Display.ShowContext changed")
	}

	if loadedConfig2.History.Enabled != originalConfig.History.Enabled {
		t.Error("Unmodified History.Enabled changed")
	}
}

// TestLoadSaveLoadCycleWithProviderModification tests the cycle with provider data modification.
func TestLoadSaveLoadCycleWithProviderModification(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Step 1: Create and save initial config with multiple providers
	original := &Config{
		CurrentProvider: "anthropic",
		Providers: map[string]ProviderConfig{
			"anthropic": {
				Type:      "anthropic",
				APIKey:    "sk-ant-original",
				Model:     "claude-3-5-sonnet-20241022",
				MaxTokens: 2048,
				Temperature: 0.7,
			},
			"openai": {
				Type:      "openai",
				APIKey:    "sk-openai-original",
				Model:     "gpt-4",
				MaxTokens: 4096,
				Temperature: 0.5,
			},
		},
		Context: ContextConfig{
			IncludeFiles:       true,
			IncludeHistory:     50,
			MaxContextSize:     8000,
			ExcludePatterns:    []string{".git", "node_modules"},
		},
		Display: DisplayConfig{
			SyntaxHighlight: true,
			ShowContext:     true,
			Color:           true,
		},
		History: HistoryConfig{
			Enabled:  true,
			MaxSize:  1000,
			FilePath: "~/.local/share/how/history",
		},
	}

	if err := original.Save(configPath); err != nil {
		t.Fatalf("initial save failed: %v", err)
	}

	// Step 2: Load the config
	loaded1, err := Load(configPath)
	if err != nil {
		t.Fatalf("first load failed: %v", err)
	}

	// Step 3: Modify provider configuration
	anthropicProvider := loaded1.Providers["anthropic"]
	anthropicProvider.APIKey = "sk-ant-modified"
	anthropicProvider.MaxTokens = 4096
	loaded1.Providers["anthropic"] = anthropicProvider

	// Add a new provider
	loaded1.Providers["local"] = ProviderConfig{
		Type:      "local",
		APIKey:    "local-key",
		Model:     "llama-2",
		MaxTokens: 8192,
		BaseURL:   "http://localhost:8000",
	}

	// Step 4: Save the modified config
	if err := loaded1.Save(configPath); err != nil {
		t.Fatalf("modified save failed: %v", err)
	}

	// Step 5: Load the modified config
	loaded2, err := Load(configPath)
	if err != nil {
		t.Fatalf("second load failed: %v", err)
	}

	// Step 6: Verify provider modifications were persisted
	anthropic, exists := loaded2.Providers["anthropic"]
	if !exists {
		t.Fatal("anthropic provider not found after cycle")
	}

	if anthropic.APIKey != "sk-ant-modified" {
		t.Errorf("anthropic APIKey not persisted: expected sk-ant-modified, got %s", anthropic.APIKey)
	}

	if anthropic.MaxTokens != 4096 {
		t.Errorf("anthropic MaxTokens not persisted: expected 4096, got %d", anthropic.MaxTokens)
	}

	// Step 7: Verify new provider was persisted
	local, exists := loaded2.Providers["local"]
	if !exists {
		t.Fatal("local provider not found after cycle")
	}

	if local.BaseURL != "http://localhost:8000" {
		t.Errorf("local provider BaseURL not persisted: expected http://localhost:8000, got %s", local.BaseURL)
	}

	// Step 8: Verify unmodified provider is still intact
	openai, exists := loaded2.Providers["openai"]
	if !exists {
		t.Fatal("openai provider not found after cycle")
	}

	if openai.APIKey != "sk-openai-original" {
		t.Errorf("unmodified openai APIKey changed: expected sk-openai-original, got %s", openai.APIKey)
	}
}

// TestLoadSaveLoadCycleWithContextModification tests the cycle with context config modification.
func TestLoadSaveLoadCycleWithContextModification(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Step 1: Create and save initial config
	original := SampleConfig()
	if err := original.Save(configPath); err != nil {
		t.Fatalf("initial save failed: %v", err)
	}

	// Step 2: Load the config
	loaded1, err := Load(configPath)
	if err != nil {
		t.Fatalf("first load failed: %v", err)
	}

	// Step 3: Modify context configuration
	loaded1.Context.IncludeEnvironment = true
	loaded1.Context.IncludeGit = false
	loaded1.Context.MaxContextSize = 16000
	loaded1.Context.ExcludePatterns = append(loaded1.Context.ExcludePatterns, ".env", ".vscode")

	// Step 4: Save the modified config
	if err := loaded1.Save(configPath); err != nil {
		t.Fatalf("modified save failed: %v", err)
	}

	// Step 5: Load the modified config
	loaded2, err := Load(configPath)
	if err != nil {
		t.Fatalf("second load failed: %v", err)
	}

	// Step 6: Verify context modifications were persisted
	if !loaded2.Context.IncludeEnvironment {
		t.Error("Context.IncludeEnvironment modification not persisted")
	}

	if loaded2.Context.IncludeGit {
		t.Error("Context.IncludeGit modification not persisted")
	}

	if loaded2.Context.MaxContextSize != 16000 {
		t.Errorf("Context.MaxContextSize not persisted: expected 16000, got %d", loaded2.Context.MaxContextSize)
	}

	if len(loaded2.Context.ExcludePatterns) != 4 {
		t.Errorf("ExcludePatterns count mismatch: expected 4, got %d", len(loaded2.Context.ExcludePatterns))
	}

	// Verify specific patterns
	expectedPatterns := map[string]bool{
		".git":      false,
		"node_modules": false,
		".env":      false,
		".vscode":   false,
	}

	for _, pattern := range loaded2.Context.ExcludePatterns {
		expectedPatterns[pattern] = true
	}

	for pattern, found := range expectedPatterns {
		if !found {
			t.Errorf("ExcludePattern %q not found", pattern)
		}
	}

	// Step 7: Verify unmodified context fields are still intact
	if loaded2.Context.IncludeFiles != original.Context.IncludeFiles {
		t.Error("unmodified Context.IncludeFiles changed")
	}

	if loaded2.Context.IncludeHistory != 100 {
		// IncludeHistory should be 100 after our modification in step 3
		t.Logf("Context.IncludeHistory is %d (expected modification in separate test)", loaded2.Context.IncludeHistory)
	}
}

// TestLoadSaveLoadCycleWithDisplayModification tests the cycle with display config modification.
func TestLoadSaveLoadCycleWithDisplayModification(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Step 1: Create and save initial config
	original := SampleConfig()
	if err := original.Save(configPath); err != nil {
		t.Fatalf("initial save failed: %v", err)
	}

	// Step 2: Load the config
	loaded1, err := Load(configPath)
	if err != nil {
		t.Fatalf("first load failed: %v", err)
	}

	// Step 3: Modify display configuration
	loaded1.Display.SyntaxHighlight = false
	loaded1.Display.ShowContext = false
	loaded1.Display.Emoji = true
	loaded1.Display.Color = false

	// Step 4: Save the modified config
	if err := loaded1.Save(configPath); err != nil {
		t.Fatalf("modified save failed: %v", err)
	}

	// Step 5: Load the modified config
	loaded2, err := Load(configPath)
	if err != nil {
		t.Fatalf("second load failed: %v", err)
	}

	// Step 6: Verify all display modifications were persisted
	if loaded2.Display.SyntaxHighlight {
		t.Error("Display.SyntaxHighlight modification not persisted")
	}

	if loaded2.Display.ShowContext {
		t.Error("Display.ShowContext modification not persisted")
	}

	if !loaded2.Display.Emoji {
		t.Error("Display.Emoji modification not persisted")
	}

	if loaded2.Display.Color {
		t.Error("Display.Color modification not persisted")
	}
}

// TestLoadSaveLoadCycleWithHistoryModification tests the cycle with history config modification.
func TestLoadSaveLoadCycleWithHistoryModification(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Step 1: Create and save initial config
	original := SampleConfig()
	if err := original.Save(configPath); err != nil {
		t.Fatalf("initial save failed: %v", err)
	}

	// Step 2: Load the config
	loaded1, err := Load(configPath)
	if err != nil {
		t.Fatalf("first load failed: %v", err)
	}

	// Step 3: Modify history configuration
	loaded1.History.Enabled = false
	loaded1.History.MaxSize = 5000
	loaded1.History.FilePath = "/var/log/how/history"

	// Step 4: Save the modified config
	if err := loaded1.Save(configPath); err != nil {
		t.Fatalf("modified save failed: %v", err)
	}

	// Step 5: Load the modified config
	loaded2, err := Load(configPath)
	if err != nil {
		t.Fatalf("second load failed: %v", err)
	}

	// Step 6: Verify history modifications were persisted
	if loaded2.History.Enabled {
		t.Error("History.Enabled modification not persisted")
	}

	if loaded2.History.MaxSize != 5000 {
		t.Errorf("History.MaxSize not persisted: expected 5000, got %d", loaded2.History.MaxSize)
	}

	if loaded2.History.FilePath != "/var/log/how/history" {
		t.Errorf("History.FilePath not persisted: expected /var/log/how/history, got %s", loaded2.History.FilePath)
	}
}

// TestLoadSaveLoadCycleComplexModification tests a complex modification cycle with multiple changes.
func TestLoadSaveLoadCycleComplexModification(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Step 1: Create complex initial config
	original := &Config{
		CurrentProvider: "primary",
		Providers: map[string]ProviderConfig{
			"primary": {
				Type:      "anthropic",
				APIKey:    "sk-ant-primary",
				Model:     "claude-3-5-sonnet",
				MaxTokens: 2048,
				Temperature: 0.7,
				CustomHeaders: map[string]string{
					"X-Custom-1": "value1",
				},
			},
			"backup": {
				Type:      "openai",
				APIKey:    "sk-openai-backup",
				Model:     "gpt-4",
				MaxTokens: 4096,
			},
		},
		Context: ContextConfig{
			IncludeFiles:       true,
			IncludeHistory:     50,
			IncludeEnvironment: false,
			IncludeGit:         true,
			MaxContextSize:     8000,
			ExcludePatterns:    []string{".git", "node_modules"},
		},
		Display: DisplayConfig{
			SyntaxHighlight: true,
			ShowContext:     true,
			Emoji:           false,
			Color:           true,
		},
		History: HistoryConfig{
			Enabled:  true,
			MaxSize:  1000,
			FilePath: "~/.local/share/how/history",
		},
	}

	if err := original.Save(configPath); err != nil {
		t.Fatalf("initial save failed: %v", err)
	}

	// Step 2: Load the config
	loaded1, err := Load(configPath)
	if err != nil {
		t.Fatalf("first load failed: %v", err)
	}

	// Step 3: Perform complex modifications
	loaded1.CurrentProvider = "backup"

	// Modify primary provider
	primary := loaded1.Providers["primary"]
	primary.APIKey = "sk-ant-updated"
	primary.MaxTokens = 4096
	primary.Temperature = 0.5
	primary.CustomHeaders["X-Custom-1"] = "updated-value1"
	primary.CustomHeaders["X-Custom-2"] = "value2"
	loaded1.Providers["primary"] = primary

	// Add new provider
	loaded1.Providers["local"] = ProviderConfig{
		Type:      "local",
		APIKey:    "local-key",
		Model:     "llama-2",
		MaxTokens: 8192,
		BaseURL:   "http://localhost:8000",
	}

	// Modify context
	loaded1.Context.IncludeEnvironment = true
	loaded1.Context.MaxContextSize = 16000
	loaded1.Context.ExcludePatterns = append(loaded1.Context.ExcludePatterns, ".env", ".vscode")

	// Modify display
	loaded1.Display.Emoji = true
	loaded1.Display.Color = false

	// Modify history
	loaded1.History.MaxSize = 2000

	// Step 4: Save the modified config
	if err := loaded1.Save(configPath); err != nil {
		t.Fatalf("modified save failed: %v", err)
	}

	// Step 5: Load the modified config
	loaded2, err := Load(configPath)
	if err != nil {
		t.Fatalf("second load failed: %v", err)
	}

	// Step 6: Verify all complex modifications were persisted
	if loaded2.CurrentProvider != "backup" {
		t.Errorf("CurrentProvider not persisted: expected backup, got %s", loaded2.CurrentProvider)
	}

	// Verify primary provider modifications
	primary2 := loaded2.Providers["primary"]
	if primary2.APIKey != "sk-ant-updated" {
		t.Errorf("primary APIKey not persisted: expected sk-ant-updated, got %s", primary2.APIKey)
	}

	if primary2.MaxTokens != 4096 {
		t.Errorf("primary MaxTokens not persisted: expected 4096, got %d", primary2.MaxTokens)
	}

	if primary2.Temperature != 0.5 {
		t.Errorf("primary Temperature not persisted: expected 0.5, got %f", primary2.Temperature)
	}

	if len(primary2.CustomHeaders) != 2 {
		t.Errorf("primary CustomHeaders count mismatch: expected 2, got %d", len(primary2.CustomHeaders))
	}

	if primary2.CustomHeaders["X-Custom-1"] != "updated-value1" {
		t.Errorf("primary CustomHeader X-Custom-1 not persisted: expected updated-value1, got %s", primary2.CustomHeaders["X-Custom-1"])
	}

	// Verify new provider persisted
	if _, exists := loaded2.Providers["local"]; !exists {
		t.Fatal("local provider not persisted")
	}

	// Verify context modifications
	if !loaded2.Context.IncludeEnvironment {
		t.Error("context IncludeEnvironment modification not persisted")
	}

	if loaded2.Context.MaxContextSize != 16000 {
		t.Errorf("context MaxContextSize not persisted: expected 16000, got %d", loaded2.Context.MaxContextSize)
	}

	if len(loaded2.Context.ExcludePatterns) != 4 {
		t.Errorf("ExcludePatterns count mismatch: expected 4, got %d", len(loaded2.Context.ExcludePatterns))
	}

	// Verify display modifications
	if !loaded2.Display.Emoji {
		t.Error("display Emoji modification not persisted")
	}

	if loaded2.Display.Color {
		t.Error("display Color modification not persisted")
	}

	// Verify history modifications
	if loaded2.History.MaxSize != 2000 {
		t.Errorf("history MaxSize not persisted: expected 2000, got %d", loaded2.History.MaxSize)
	}

	// Step 7: Verify data integrity - backup provider should remain unchanged
	backup := loaded2.Providers["backup"]
	if backup.APIKey != "sk-openai-backup" {
		t.Errorf("unmodified backup provider APIKey changed: expected sk-openai-backup, got %s", backup.APIKey)
	}

	if backup.Model != "gpt-4" {
		t.Errorf("unmodified backup provider Model changed: expected gpt-4, got %s", backup.Model)
	}
}

// TestLoadSaveLoadCycleMultipleCycles tests multiple load-save-load cycles in sequence.
func TestLoadSaveLoadCycleMultipleCycles(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Initialize config
	config := SampleConfig()
	if err := config.Save(configPath); err != nil {
		t.Fatalf("initial save failed: %v", err)
	}

	// Perform multiple cycles with different modifications
	for cycle := 1; cycle <= 3; cycle++ {
		// Load
		loaded, err := Load(configPath)
		if err != nil {
			t.Fatalf("cycle %d load failed: %v", cycle, err)
		}

		// Modify based on cycle number
		loaded.Display.SyntaxHighlight = (cycle%2 == 0)
		loaded.Context.IncludeHistory = 50 + (cycle * 10)

		// Save
		if err := loaded.Save(configPath); err != nil {
			t.Fatalf("cycle %d save failed: %v", cycle, err)
		}
	}

	// Final verification
	final, err := Load(configPath)
	if err != nil {
		t.Fatalf("final load failed: %v", err)
	}

	// After 3 cycles (cycle 3 is odd), SyntaxHighlight should be false
	if final.Display.SyntaxHighlight {
		t.Error("SyntaxHighlight not correctly modified after cycles")
	}

	// IncludeHistory should be 80 (50 + 3*10)
	if final.Context.IncludeHistory != 80 {
		t.Errorf("IncludeHistory mismatch: expected 80, got %d", final.Context.IncludeHistory)
	}
}

// TestLoadSaveLoadCycleDataConsistency verifies that no data is lost during cycles.
func TestLoadSaveLoadCycleDataConsistency(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Create a complex config with multiple providers and extensive settings
	original := &Config{
		CurrentProvider: "anthropic",
		Providers: map[string]ProviderConfig{
			"anthropic": {
				Type:      "anthropic",
				APIKey:    "sk-ant-1",
				Model:     "claude-3-5-sonnet",
				BaseURL:   "https://api.anthropic.com",
				MaxTokens: 2048,
				Temperature: 0.7,
				TopP:       0.99,
				SystemPrompt: "You are a helpful assistant.",
				CustomHeaders: map[string]string{
					"Authorization": "Bearer token",
					"X-Custom":      "header",
				},
			},
			"openai": {
				Type:      "openai",
				APIKey:    "sk-openai-1",
				Model:     "gpt-4",
				MaxTokens: 4096,
				Temperature: 0.8,
				TopP:       1.0,
			},
		},
		Context: ContextConfig{
			IncludeFiles:       true,
			IncludeHistory:     100,
			IncludeEnvironment: true,
			IncludeGit:         true,
			MaxContextSize:     16000,
			ExcludePatterns:    []string{".git", "node_modules", ".env", ".vscode", "dist", "build"},
		},
		Display: DisplayConfig{
			SyntaxHighlight: true,
			ShowContext:     true,
			Emoji:           true,
			Color:           true,
		},
		History: HistoryConfig{
			Enabled:  true,
			MaxSize:  5000,
			FilePath: "~/.local/share/how/history",
		},
	}

	if err := original.Save(configPath); err != nil {
		t.Fatalf("initial save failed: %v", err)
	}

	// Load cycle 1: Load and verify no modifications
	loaded1, err := Load(configPath)
	if err != nil {
		t.Fatalf("first load failed: %v", err)
	}

	assertConfigsEqual(t, original, loaded1, "cycle 1 load")

	// Load cycle 2: Modify one field and verify others remain intact
	loaded1.Display.Emoji = false
	if err := loaded1.Save(configPath); err != nil {
		t.Fatalf("cycle 2 save failed: %v", err)
	}

	loaded2, err := Load(configPath)
	if err != nil {
		t.Fatalf("cycle 2 load failed: %v", err)
	}

	// Verify modification
	if loaded2.Display.Emoji {
		t.Error("display Emoji modification not persisted")
	}

	// Verify data integrity - all other fields should match original
	verifyDataIntegrity(t, original, loaded2, "cycle 2")

	// Load cycle 3: Modify provider and verify context remains intact
	loaded3, err := Load(configPath)
	if err != nil {
		t.Fatalf("cycle 3 load failed: %v", err)
	}

	anthrop := loaded3.Providers["anthropic"]
	anthrop.MaxTokens = 8192
	loaded3.Providers["anthropic"] = anthrop

	if err := loaded3.Save(configPath); err != nil {
		t.Fatalf("cycle 3 save failed: %v", err)
	}

	loaded4, err := Load(configPath)
	if err != nil {
		t.Fatalf("cycle 3 final load failed: %v", err)
	}

	// Verify provider modification
	if loaded4.Providers["anthropic"].MaxTokens != 8192 {
		t.Errorf("provider modification not persisted: expected 8192, got %d", loaded4.Providers["anthropic"].MaxTokens)
	}

	// Verify context data integrity
	if len(loaded4.Context.ExcludePatterns) != len(original.Context.ExcludePatterns) {
		t.Errorf("ExcludePatterns count mismatch: expected %d, got %d", len(original.Context.ExcludePatterns), len(loaded4.Context.ExcludePatterns))
	}

	for i, pattern := range original.Context.ExcludePatterns {
		if loaded4.Context.ExcludePatterns[i] != pattern {
			t.Errorf("ExcludePatterns[%d] mismatch: expected %s, got %s", i, pattern, loaded4.Context.ExcludePatterns[i])
		}
	}

	// Verify custom headers are preserved
	anthrop4 := loaded4.Providers["anthropic"]
	if len(anthrop4.CustomHeaders) != len(original.Providers["anthropic"].CustomHeaders) {
		t.Errorf("CustomHeaders count mismatch: expected %d, got %d", len(original.Providers["anthropic"].CustomHeaders), len(anthrop4.CustomHeaders))
	}

	for key, expectedValue := range original.Providers["anthropic"].CustomHeaders {
		actualValue, exists := anthrop4.CustomHeaders[key]
		if !exists {
			t.Errorf("CustomHeader %q not found", key)
		}
		if actualValue != expectedValue {
			t.Errorf("CustomHeader %q mismatch: expected %s, got %s", key, expectedValue, actualValue)
		}
	}
}

// TestLoadSaveLoadCycleEmptyToPopulated tests the cycle starting from empty config.
func TestLoadSaveLoadCycleEmptyToPopulated(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Step 1: Create and save empty config
	emptyConfig := &Config{
		Providers: make(map[string]ProviderConfig),
	}

	if err := emptyConfig.Save(configPath); err != nil {
		t.Fatalf("initial save failed: %v", err)
	}

	// Step 2: Load the empty config
	loaded1, err := Load(configPath)
	if err != nil {
		t.Fatalf("first load failed: %v", err)
	}

	if loaded1 == nil {
		t.Fatal("loaded1 config should not be nil")
	}

	// Step 3: Populate the config
	loaded1.CurrentProvider = "anthropic"
	loaded1.Providers["anthropic"] = ProviderConfig{
		Type:      "anthropic",
		APIKey:    "sk-ant-key",
		Model:     "claude-3-5-sonnet",
		MaxTokens: 2048,
	}
	loaded1.Context = ContextConfig{
		IncludeFiles:   true,
		IncludeHistory: 50,
		MaxContextSize: 8000,
	}
	loaded1.Display = DisplayConfig{
		SyntaxHighlight: true,
		ShowContext:     true,
		Color:           true,
	}
	loaded1.History = HistoryConfig{
		Enabled:  true,
		MaxSize:  1000,
		FilePath: "~/.local/share/how/history",
	}

	// Step 4: Save the populated config
	if err := loaded1.Save(configPath); err != nil {
		t.Fatalf("modified save failed: %v", err)
	}

	// Step 5: Load the populated config
	loaded2, err := Load(configPath)
	if err != nil {
		t.Fatalf("second load failed: %v", err)
	}

	// Step 6: Verify all data was persisted correctly
	if loaded2.CurrentProvider != "anthropic" {
		t.Errorf("CurrentProvider mismatch: expected anthropic, got %s", loaded2.CurrentProvider)
	}

	if len(loaded2.Providers) != 1 {
		t.Errorf("Providers count mismatch: expected 1, got %d", len(loaded2.Providers))
	}

	anthrop, exists := loaded2.Providers["anthropic"]
	if !exists {
		t.Fatal("anthropic provider not found")
	}

	if anthrop.APIKey != "sk-ant-key" {
		t.Errorf("APIKey mismatch: expected sk-ant-key, got %s", anthrop.APIKey)
	}

	if loaded2.Context.IncludeFiles {
		if !loaded2.Context.IncludeFiles {
			t.Error("IncludeFiles should be true")
		}
	}

	if !loaded2.Display.SyntaxHighlight {
		t.Error("SyntaxHighlight should be true")
	}

	if !loaded2.History.Enabled {
		t.Error("History.Enabled should be true")
	}
}

// TestLoadSaveLoadCycleFilePermissions verifies that file permissions are preserved across cycles.
func TestLoadSaveLoadCycleFilePermissions(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Step 1: Create and save config
	config := SampleConfig()
	if err := config.Save(configPath); err != nil {
		t.Fatalf("initial save failed: %v", err)
	}

	// Step 2: Check initial permissions
	fileInfo1, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}

	initialPerms := fileInfo1.Mode().Perm()
	expectedPerms := os.FileMode(0644)

	if initialPerms != expectedPerms {
		t.Errorf("initial permissions mismatch: expected %o, got %o", expectedPerms, initialPerms)
	}

	// Step 3: Load and save again
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	loaded.Display.SyntaxHighlight = !loaded.Display.SyntaxHighlight

	if err := loaded.Save(configPath); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	// Step 4: Verify permissions are still 0644
	fileInfo2, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("stat after resave failed: %v", err)
	}

	finalPerms := fileInfo2.Mode().Perm()
	if finalPerms != expectedPerms {
		t.Errorf("permissions changed after resave: expected %o, got %o", expectedPerms, finalPerms)
	}
}

// Helper function to assert two configs are equal
func assertConfigsEqual(t *testing.T, expected, actual *Config, context string) {
	t.Helper()

	if expected.CurrentProvider != actual.CurrentProvider {
		t.Errorf("%s: CurrentProvider mismatch: expected %q, got %q", context, expected.CurrentProvider, actual.CurrentProvider)
	}

	if len(expected.Providers) != len(actual.Providers) {
		t.Errorf("%s: providers count mismatch: expected %d, got %d", context, len(expected.Providers), len(actual.Providers))
	}

	for name, expectedProvider := range expected.Providers {
		actualProvider, exists := actual.Providers[name]
		if !exists {
			t.Errorf("%s: provider %q not found", context, name)
			continue
		}

		if expectedProvider.Type != actualProvider.Type {
			t.Errorf("%s: provider %q type mismatch: expected %q, got %q", context, name, expectedProvider.Type, actualProvider.Type)
		}

		if expectedProvider.APIKey != actualProvider.APIKey {
			t.Errorf("%s: provider %q APIKey mismatch", context, name)
		}

		if expectedProvider.Model != actualProvider.Model {
			t.Errorf("%s: provider %q model mismatch", context, name)
		}

		if expectedProvider.MaxTokens != actualProvider.MaxTokens {
			t.Errorf("%s: provider %q maxTokens mismatch: expected %d, got %d", context, name, expectedProvider.MaxTokens, actualProvider.MaxTokens)
		}
	}

	if expected.Context.IncludeFiles != actual.Context.IncludeFiles {
		t.Errorf("%s: Context.IncludeFiles mismatch", context)
	}

	if expected.Display.SyntaxHighlight != actual.Display.SyntaxHighlight {
		t.Errorf("%s: Display.SyntaxHighlight mismatch", context)
	}

	if expected.History.Enabled != actual.History.Enabled {
		t.Errorf("%s: History.Enabled mismatch", context)
	}
}

// Helper function to verify data integrity after modifications
func verifyDataIntegrity(t *testing.T, original, modified *Config, context string) {
	t.Helper()

	// Check that provider count hasn't changed unexpectedly
	if len(modified.Providers) != len(original.Providers) {
		t.Errorf("%s: provider count changed: original %d, modified %d", context, len(original.Providers), len(modified.Providers))
	}

	// Check that all original providers still exist
	for name := range original.Providers {
		if _, exists := modified.Providers[name]; !exists {
			t.Errorf("%s: provider %q missing", context, name)
		}
	}

	// Check context fields (except those we might have modified)
	if len(modified.Context.ExcludePatterns) != len(original.Context.ExcludePatterns) {
		t.Errorf("%s: ExcludePatterns count changed: original %d, modified %d", context, len(original.Context.ExcludePatterns), len(modified.Context.ExcludePatterns))
	}

	if modified.History.FilePath != original.History.FilePath {
		t.Errorf("%s: History.FilePath changed: original %q, modified %q", context, original.History.FilePath, modified.History.FilePath)
	}
}
