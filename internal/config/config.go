package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

const (
	DefaultAPIURL  = "https://starscope.app/api/v1"
	DefaultProfile = "default"
	AppName        = "starscope-cli"
)

// Config holds the application configuration.
type Config struct {
	v       *viper.Viper
	dir     string
	keyring *Keyring
}

// Profile represents a named configuration profile.
type Profile struct {
	APIURL      string `mapstructure:"api_url"`
	WorkspaceID int    `mapstructure:"workspace_id"`
	Token       string `mapstructure:"token"`
}

// New creates a new Config instance, loading from the config file if it exists.
func New() (*Config, error) {
	dir, err := configDir()
	if err != nil {
		return nil, fmt.Errorf("determining config directory: %w", err)
	}

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(dir)

	v.SetDefault("current_profile", DefaultProfile)

	v.SetEnvPrefix("STARSCOPE")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("reading config: %w", err)
			}
		}
	}

	return &Config{
		v:       v,
		dir:     dir,
		keyring: NewKeyring(),
	}, nil
}

// CurrentProfile returns the name of the active profile.
func (c *Config) CurrentProfile() string {
	return c.v.GetString("current_profile")
}

// SetCurrentProfile sets the active profile.
func (c *Config) SetCurrentProfile(name string) {
	c.v.Set("current_profile", name)
}

// GetProfile returns the profile with the given name.
func (c *Config) GetProfile(name string) Profile {
	key := fmt.Sprintf("profiles.%s", name)
	return Profile{
		APIURL:      c.v.GetString(key + ".api_url"),
		WorkspaceID: c.v.GetInt(key + ".workspace_id"),
		Token:       c.v.GetString(key + ".token"),
	}
}

// SetProfileField sets a field on the given profile.
func (c *Config) SetProfileField(profile, field string, value any) {
	key := fmt.Sprintf("profiles.%s.%s", profile, field)
	c.v.Set(key, value)
}

// ResolveToken returns the API token using the priority chain:
// flag/env override > keychain > config file.
func (c *Config) ResolveToken(flagToken string) string {
	if flagToken != "" {
		return flagToken
	}

	if envToken := os.Getenv("STARSCOPE_TOKEN"); envToken != "" {
		return envToken
	}

	profile := c.CurrentProfile()
	if token, err := c.keyring.Get(profile); err == nil && token != "" {
		return token
	}

	return c.GetProfile(profile).Token
}

// ResolveAPIURL returns the API URL using the priority chain.
func (c *Config) ResolveAPIURL(flagURL string) string {
	if flagURL != "" {
		return flagURL
	}

	profile := c.GetProfile(c.CurrentProfile())
	if profile.APIURL != "" {
		return profile.APIURL
	}

	return DefaultAPIURL
}

// ResolveWorkspaceID returns the workspace ID using the priority chain.
func (c *Config) ResolveWorkspaceID(flagID int) int {
	if flagID != 0 {
		return flagID
	}

	return c.GetProfile(c.CurrentProfile()).WorkspaceID
}

// StoreToken stores the token in the OS keychain.
// Falls back to config file if keychain is unavailable.
func (c *Config) StoreToken(profile, token string) error {
	if err := c.keyring.Set(profile, token); err != nil {
		c.SetProfileField(profile, "token", token)
		return nil
	}
	return nil
}

// RemoveToken removes the token from the OS keychain and config.
func (c *Config) RemoveToken(profile string) error {
	_ = c.keyring.Delete(profile)
	c.SetProfileField(profile, "token", "")
	return nil
}

// Save writes the config to disk.
func (c *Config) Save() error {
	if err := os.MkdirAll(c.dir, 0o700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	configPath := filepath.Join(c.dir, "config.yaml")
	return c.v.WriteConfigAs(configPath)
}

// Get returns a raw config value.
func (c *Config) Get(key string) any {
	return c.v.Get(key)
}

// Set sets a raw config value.
func (c *Config) Set(key string, value any) {
	c.v.Set(key, value)
}

// AllSettings returns all config values.
func (c *Config) AllSettings() map[string]any {
	return c.v.AllSettings()
}

// Dir returns the config directory path.
func (c *Config) Dir() string {
	return c.dir
}

func configDir() (string, error) {
	if dir := os.Getenv("STARSCOPE_CONFIG_DIR"); dir != "" {
		return dir, nil
	}

	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("APPDATA not set")
		}
		return filepath.Join(appData, AppName), nil
	default:
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		xdgConfig := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfig == "" {
			xdgConfig = filepath.Join(home, ".config")
		}
		return filepath.Join(xdgConfig, AppName), nil
	}
}
