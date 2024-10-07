package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	VersionMajor      string `mapstructure:"VER_MAJOR"`
	VersionMinor      string `mapstructure:"VER_MINOR"`
	VersionPatch      string `mapstructure:"VER_PATCH"`
	VersionPre        string `mapstructure:"VER_PRE"`
	SearchRoot        string `mapstructure:"SEARCH_ROOT"`
	GitPullUser       string `mapstructure:"GIT_PULL_USER"`
	SlackWebhook      string `mapstructure:"SLACK_WEBHOOK"`
	RepoGitRemote     string `mapstructure:"REPO_GIT_REMOTE"`
	BeforeRestart     string `mapstructure:"BEFORE_RESTART"`
	AfterRestart      string `mapstructure:"AFTER_RESTART"`
	RestartWithBuild  bool   `mapstructure:"RESTART_WITH_BUILD"`
	UpdateRepoOnly    bool   `mapstructure:"UPDATE_REPO_ONLY"`
	UpdateImageByRepo bool   `mapstructure:"UPDATE_IMAGE_BY_REPO"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	// 環境変数のプレフィックスを設定
	viper.SetEnvPrefix("SAKKYOKU")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}
