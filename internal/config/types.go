package config

import (
	"encoding/json"
)

type CustomModel struct {
	DisplayName string `json:"model_display_name"`
	Model       string `json:"model"`
	BaseURL     string `json:"base_url"`
	APIKey      string `json:"api_key"`
	Provider    string `json:"provider"`
	MaxTokens   int    `json:"max_tokens"`
}

type ConfigData struct {
	CustomModels []CustomModel
	extra        map[string]json.RawMessage
}

func (c *ConfigData) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if cm, ok := raw["custom_models"]; ok {
		if err := json.Unmarshal(cm, &c.CustomModels); err != nil {
			return err
		}
		delete(raw, "custom_models")
	}

	c.extra = raw
	return nil
}

func (c ConfigData) MarshalJSON() ([]byte, error) {
	result := make(map[string]json.RawMessage)
	for k, v := range c.extra {
		result[k] = v
	}

	cm, err := json.Marshal(c.CustomModels)
	if err != nil {
		return nil, err
	}
	result["custom_models"] = cm

	return json.Marshal(result)
}

var Providers = []string{
	"anthropic",
	"openai",
	"generic-chat-completion-api",
}
