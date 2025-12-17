package config

import (
	"encoding/json"
	"testing"
)

func TestConfigDataMarshalUnmarshal(t *testing.T) {
	input := `{
		"custom_models": [
			{
				"model_display_name": "Test Model",
				"model": "test-model",
				"base_url": "https://api.test.com",
				"api_key": "sk-test",
				"provider": "openai",
				"max_tokens": 4096
			}
		],
		"other_field": "should be preserved",
		"nested": {"key": "value"}
	}`

	var cfg ConfigData
	if err := json.Unmarshal([]byte(input), &cfg); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(cfg.CustomModels) != 1 {
		t.Fatalf("Expected 1 custom model, got %d", len(cfg.CustomModels))
	}

	m := cfg.CustomModels[0]
	if m.DisplayName != "Test Model" {
		t.Errorf("Expected DisplayName 'Test Model', got '%s'", m.DisplayName)
	}
	if m.Model != "test-model" {
		t.Errorf("Expected Model 'test-model', got '%s'", m.Model)
	}
	if m.MaxTokens != 4096 {
		t.Errorf("Expected MaxTokens 4096, got %d", m.MaxTokens)
	}

	output, err := json.Marshal(&cfg)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}

	if _, ok := result["other_field"]; !ok {
		t.Error("other_field was not preserved")
	}
	if _, ok := result["nested"]; !ok {
		t.Error("nested field was not preserved")
	}
	if _, ok := result["custom_models"]; !ok {
		t.Error("custom_models was not included")
	}
}

func TestEmptyConfig(t *testing.T) {
	input := `{}`

	var cfg ConfigData
	if err := json.Unmarshal([]byte(input), &cfg); err != nil {
		t.Fatalf("Failed to unmarshal empty config: %v", err)
	}

	if cfg.CustomModels == nil {
		cfg.CustomModels = []CustomModel{}
	}

	if len(cfg.CustomModels) != 0 {
		t.Errorf("Expected 0 custom models, got %d", len(cfg.CustomModels))
	}
}

func TestProviders(t *testing.T) {
	expected := []string{"anthropic", "openai", "generic-chat-completion-api"}
	
	if len(Providers) != len(expected) {
		t.Fatalf("Expected %d providers, got %d", len(expected), len(Providers))
	}

	for i, p := range expected {
		if Providers[i] != p {
			t.Errorf("Expected provider[%d] to be '%s', got '%s'", i, p, Providers[i])
		}
	}
}
