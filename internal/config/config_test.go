package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_DefaultProfile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("STARSCOPE_CONFIG_DIR", dir)

	cfg, err := New()
	require.NoError(t, err)
	assert.Equal(t, "default", cfg.CurrentProfile())
}

func TestConfig_SetAndGet(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("STARSCOPE_CONFIG_DIR", dir)

	cfg, err := New()
	require.NoError(t, err)

	cfg.Set("profiles.test.api_url", "https://test.example.com")
	assert.Equal(t, "https://test.example.com", cfg.Get("profiles.test.api_url"))
}

func TestConfig_Save(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("STARSCOPE_CONFIG_DIR", dir)

	cfg, err := New()
	require.NoError(t, err)

	cfg.SetProfileField("default", "api_url", "https://starscope.app/api/v1")
	cfg.SetProfileField("default", "workspace_id", 42)

	err = cfg.Save()
	require.NoError(t, err)

	// Verify file was written
	configPath := filepath.Join(dir, "config.yaml")
	_, err = os.Stat(configPath)
	assert.NoError(t, err)

	// Reload and verify
	cfg2, err := New()
	require.NoError(t, err)
	profile := cfg2.GetProfile("default")
	assert.Equal(t, "https://starscope.app/api/v1", profile.APIURL)
	assert.Equal(t, 42, profile.WorkspaceID)
}

func TestConfig_ResolveToken_FlagFirst(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("STARSCOPE_CONFIG_DIR", dir)

	cfg, err := New()
	require.NoError(t, err)

	cfg.SetProfileField("default", "token", "config-token")
	token := cfg.ResolveToken("flag-token")
	assert.Equal(t, "flag-token", token)
}

func TestConfig_ResolveToken_EnvSecond(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("STARSCOPE_CONFIG_DIR", dir)
	t.Setenv("STARSCOPE_TOKEN", "env-token")

	cfg, err := New()
	require.NoError(t, err)

	token := cfg.ResolveToken("")
	assert.Equal(t, "env-token", token)
}

func TestConfig_ResolveAPIURL_Default(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("STARSCOPE_CONFIG_DIR", dir)

	cfg, err := New()
	require.NoError(t, err)

	url := cfg.ResolveAPIURL("")
	assert.Equal(t, DefaultAPIURL, url)
}

func TestConfig_ResolveAPIURL_Override(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("STARSCOPE_CONFIG_DIR", dir)

	cfg, err := New()
	require.NoError(t, err)

	url := cfg.ResolveAPIURL("https://custom.example.com")
	assert.Equal(t, "https://custom.example.com", url)
}

func TestConfig_ProfileSwitch(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("STARSCOPE_CONFIG_DIR", dir)

	cfg, err := New()
	require.NoError(t, err)

	cfg.SetProfileField("staging", "api_url", "https://staging.example.com")
	cfg.SetCurrentProfile("staging")

	assert.Equal(t, "staging", cfg.CurrentProfile())
	profile := cfg.GetProfile("staging")
	assert.Equal(t, "https://staging.example.com", profile.APIURL)
}
