package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_ValidFullYAML(t *testing.T) {

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "full.yaml")

	// This YAML exercises every field in your Config struct
	yamlContent := `
supervisor:
  name: "main-supervisor"
  log_dir: "/var/log/gosuper"
  restart_delay: "5s"
  stop_timeout: "10s"
  env:
    PATH: "/usr/bin"
    GLOBAL: "true"

services:
  web:
    command: ["/usr/bin/sleep", "3600"]
    dir: "/var/www"
    restart_window: "30s"
    stdout: "/var/log/web.out"
    stderr: "/var/log/web.err"
    env:
      PORT: "8080"
      MODE: "production"
    restart_limit: 3
    autostart: true
    autorestart: true

  cleaner:
    command: ["/bin/rm", "-rf", "/tmp/cache"]
    dir: "/tmp"
    restart_window: "0s"
    stdout: ""
    stderr: ""
    env: {}
    restart_limit: 0
    autostart: false
    autorestart: false
`
	err := os.WriteFile(configPath, []byte(yamlContent), 0644)
	assert.NoError(t, err)

	cfg, err := LoadConfig(configPath)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// ----------------- Supervisor -----------------
	assert.NotNil(t, cfg.Supervisor)
	assert.Equal(t, "main-supervisor", cfg.Supervisor.Name)
	assert.Equal(t, "/var/log/gosuper", cfg.Supervisor.LogDir)
	assert.Equal(t, "5s", cfg.Supervisor.RestartDelay)
	assert.Equal(t, "10s", cfg.Supervisor.StopTimeout)
	assert.NotNil(t, cfg.Supervisor.Env)
	assert.Equal(t, 2, len(cfg.Supervisor.Env))
	assert.Equal(t, "/usr/bin", cfg.Supervisor.Env["PATH"])
	assert.Equal(t, "true", cfg.Supervisor.Env["GLOBAL"])

	// ----------------- Services -----------------
	assert.Equal(t, 2, len(cfg.Services))

	// --- web service ---
	webSvc, exists := cfg.Services["web"]
	assert.True(t, exists)
	assert.Equal(t, []string{"/usr/bin/sleep", "3600"}, webSvc.Command)
	assert.Equal(t, "/var/www", webSvc.Dir)
	assert.Equal(t, "30s", webSvc.RestartWindow)
	assert.Equal(t, "/var/log/web.out", webSvc.Stdout)
	assert.Equal(t, "/var/log/web.err", webSvc.Stderr)
	assert.Equal(t, map[string]string{"PORT": "8080", "MODE": "production"}, webSvc.Env)
	assert.Equal(t, 3, webSvc.RestartLimit)
	assert.True(t, webSvc.Autostart)
	assert.True(t, webSvc.Autorestart)

	// --- cleaner service (minimal fields) ---
	cleaner, exists := cfg.Services["cleaner"]
	assert.True(t, exists)
	assert.Equal(t, []string{"/bin/rm", "-rf", "/tmp/cache"}, cleaner.Command)
	assert.Equal(t, "/tmp", cleaner.Dir)
	assert.Equal(t, "0s", cleaner.RestartWindow)
	assert.Equal(t, "", cleaner.Stdout)
	assert.Equal(t, "", cleaner.Stderr)
	// Env defined as {} -> should be an empty map, not nil
	assert.NotNil(t, cleaner.Env)
	assert.Empty(t, cleaner.Env)
	assert.Equal(t, 0, cleaner.RestartLimit)
	assert.False(t, cleaner.Autostart)
	assert.False(t, cleaner.Autorestart)
}

func TestLoadConfig_FileNotFound(t *testing.T) {

	cfg, err := LoadConfig("/nonexistent/config.yaml")
	assert.Nil(t, cfg)
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}

func TestLoadConfig_InvalidYAML(t *testing.T) {

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "broken.yaml")

	broken := `supervisor: [this is not a valid yaml structure`
	err := os.WriteFile(configPath, []byte(broken), 0644)
	assert.NoError(t, err)

	cfg, err := LoadConfig(configPath)
	assert.Nil(t, cfg)
	assert.Error(t, err)
}

func TestLoadConfig_EmptyFile(t *testing.T) {

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "empty.yaml")
	err := os.WriteFile(configPath, []byte{}, 0644)
	assert.NoError(t, err)

	cfg, err := LoadConfig(configPath)
	assert.Error(t, err) // yaml.Unmarshal with empty input gives no error
	assert.Nil(t, cfg)
}

func TestLoadConfig_OnlySupervisor(t *testing.T) {

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "supervisor_only.yaml")
	yamlContent := `
supervisor:
  name: "bare"
  log_dir: "/dev/null"
  restart_delay: "1s"
  stop_timeout: "5s"
`
	err := os.WriteFile(configPath, []byte(yamlContent), 0644)
	assert.NoError(t, err)

	cfg, err := LoadConfig(configPath)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.NotNil(t, cfg.Supervisor)
	assert.Equal(t, "bare", cfg.Supervisor.Name)
	assert.Equal(t, 0, len(cfg.Services)) // no services defined
}

func TestLoadConfig_ServiceWithDefaultsMissing(t *testing.T) {

	// When you omit optional fields, YAML should set them to zero values.
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "minimal.yaml")
	yamlContent := `
services:
  ping:
    command: ["/bin/ping", "localhost"]
`
	err := os.WriteFile(configPath, []byte(yamlContent), 0644)
	assert.NoError(t, err)

	cfg, err := LoadConfig(configPath)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(cfg.Services))
	pingSvc := cfg.Services["ping"]
	assert.Equal(t, []string{"/bin/ping", "localhost"}, pingSvc.Command)
	// Zero values
	assert.Equal(t, "", pingSvc.Dir)
	assert.Equal(t, "", pingSvc.RestartWindow)
	assert.Equal(t, "", pingSvc.Stdout)
	assert.Equal(t, "", pingSvc.Stderr)
	assert.Nil(t, pingSvc.Env) // note: yaml will read env: (no value) as nil, not empty map
	assert.Equal(t, 0, pingSvc.RestartLimit)
	assert.False(t, pingSvc.Autostart)
	assert.False(t, pingSvc.Autorestart)
}
