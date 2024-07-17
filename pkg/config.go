package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type ConfigurationData struct {
	APIKey  string     `json:"api_key"`
	BaseURL string     `json:"base_url"`
	OrgID   string     `json:"org_id"`
	Last    *time.Time `json:"last"`
	path    string     `json:"-"`
}

func ReadConfigFromEnv() (*ConfigurationData, error) {
	configs := map[string]string{}
	for _, config := range []string{"API_KEY", "BASE_URL", "ORG_ID"} {
		value := os.Getenv(config)
		if value == "" {
			return nil, fmt.Errorf("%s environment variable not set", config)
		}
		configs[config] = value
	}
	return &ConfigurationData{
		APIKey:  configs["API_KEY"],
		BaseURL: configs["BASE_URL"],
		OrgID:   configs["ORG_ID"],
		Last:    nil,
		path:    "env-variable",
	}, nil
}

func ReadConfigFile(path string) (*ConfigurationData, error) {
	config := ConfigurationData{}
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(contents, &config)
	if err != nil {
		return nil, err
	}
	config.path = path
	return &config, nil
}

func (c *ConfigurationData) UpdateLast(newTime time.Time) error {
	c.Last = &newTime
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = os.WriteFile(c.path, b, 0644)
	return err
}

func (c *ConfigurationData) GetLastTime() time.Time {
	if c.Last == nil {
		return time.Now().Add(-time.Hour * 1)
	}
	return *c.Last
}
