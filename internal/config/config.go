package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ContextEntry holds the connection details for a named context.
type ContextEntry struct {
	APIURL string `json:"api_url"`
	Token  string `json:"token"`
}

// Config is the root config stored at ~/.mikrom/config.json.
// The flat APIURL/Token fields are kept for backward compatibility and always
// reflect the active context. Contexts holds all named contexts.
type Config struct {
	APIURL         string                  `json:"api_url"`
	Token          string                  `json:"token"`
	CurrentContext string                  `json:"current_context,omitempty"`
	Contexts       map[string]ContextEntry `json:"contexts,omitempty"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".mikrom", "config.json"), nil
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{APIURL: "http://localhost:8080"}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// If a context is active, sync APIURL/Token from that context entry so
	// callers can always use cfg.APIURL and cfg.Token directly.
	if cfg.CurrentContext != "" {
		if entry, ok := cfg.Contexts[cfg.CurrentContext]; ok {
			cfg.APIURL = entry.APIURL
			cfg.Token = entry.Token
		}
	}

	return &cfg, nil
}

// Save writes the config to disk. It syncs the current APIURL/Token into the
// active context entry before serialising.
func (c *Config) Save() error {
	// Sync the active context entry.
	if c.CurrentContext != "" {
		if c.Contexts == nil {
			c.Contexts = make(map[string]ContextEntry)
		}
		c.Contexts[c.CurrentContext] = ContextEntry{APIURL: c.APIURL, Token: c.Token}
	}

	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// ActiveContext returns the name of the currently active context. When no
// context has been explicitly set it returns "default".
func (c *Config) ActiveContext() string {
	if c.CurrentContext == "" {
		return "default"
	}
	return c.CurrentContext
}

// AddContext creates or overwrites a named context.
func (c *Config) AddContext(name, apiURL, token string) {
	if c.Contexts == nil {
		c.Contexts = make(map[string]ContextEntry)
	}
	c.Contexts[name] = ContextEntry{APIURL: apiURL, Token: token}
}

// UseContext switches the active context. It persists the current APIURL/Token
// into the outgoing context before switching.
func (c *Config) UseContext(name string) error {
	if c.Contexts == nil || c.Contexts[name] == (ContextEntry{}) {
		return fmt.Errorf("context %q not found — use 'mikrom context add' to create it", name)
	}

	// Persist current values into the outgoing context.
	if c.CurrentContext != "" {
		c.Contexts[c.CurrentContext] = ContextEntry{APIURL: c.APIURL, Token: c.Token}
	}

	entry := c.Contexts[name]
	c.CurrentContext = name
	c.APIURL = entry.APIURL
	c.Token = entry.Token
	return nil
}

// RemoveContext deletes a named context. The active context cannot be removed.
func (c *Config) RemoveContext(name string) error {
	if name == c.ActiveContext() {
		return fmt.Errorf("cannot remove the active context %q — switch first with 'mikrom context use'", name)
	}
	if c.Contexts == nil {
		return fmt.Errorf("context %q not found", name)
	}
	if _, ok := c.Contexts[name]; !ok {
		return fmt.Errorf("context %q not found", name)
	}
	delete(c.Contexts, name)
	return nil
}
