package config

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		HTTP   `yaml:"http"`
		Logger `yaml:"logger"`
		GCP    `yaml:"gcp"`
		Alert  `yaml:"alert"`
	}

	// HTTP -.
	HTTP struct {
		Port             string   `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		AllowedMethods   []string `env-required:"true" yaml:"allowed_methods" env:"HTTP_ALLOWED_METHODS"`
		AllowedHeaders   []string `env-required:"true" yaml:"allowed_headers" env:"HTTP_ALLOWED_HEADERS"`
		ExposeHeaders    []string `env-required:"true" yaml:"expose_headers" env:"HTTP_EXPOSE_HEADERS"`
		AllowCredentials bool     `env-required:"true" yaml:"allow_credentials" env:"HTTP_ALLOW_CREDENTIALS"`
	}

	// Log -.
	Logger struct {
		Level  string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
		Format string `env-required:"true" yaml:"log_format"  env:"LOG_FORMAT"`
	}

	// GCP -.
	GCP struct {
		ProjectID       string  `env-required:"true" yaml:"project_id" env:"GOOGLE_CLOUD_PROJECT_ID"`
		CredentialsFile *string `yaml:"credentials_file,omitempty" env:"GOOGLE_APPLICATION_CREDENTIALS"`
	}

	// Alert -.
	Alert struct {
		Provider string `yaml:"provider" env:"ALERT_PROVIDER"`
		Slack    `yaml:"slack"`
	}

	// Slack -.
	Slack struct {
		WebhookURL string `yaml:"webhook_url" env:"SLACK_WEBHOOK_URL"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("Could not get current file info")
	}
	path := filepath.Join(filepath.Dir(filename), "..", "..")
	err := cleanenv.ReadConfig(filepath.Join(path, "config", "config.yaml"), cfg)
	if err != nil {
		return nil, fmt.Errorf("Config Error: %w", err)
	}
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
