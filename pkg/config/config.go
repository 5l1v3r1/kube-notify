package config

import "github.com/BurntSushi/toml"

type Config struct {
	KubeNotify KubeNotifyConf
	Slack      SlackConf
}

type KubeNotifyConf struct {
	IgnoreNotified bool   `toml:"ignoreNotified,omitempty"`
	LocalMode      bool   `toml:"localMode,omitempty"`
	ConfigPath     string `toml:"configPath,omitempty"`
}

type SlackConf struct {
	HookURL     string   `toml:"hookURL,omitempty"`
	Token       string   `toml:"token,omitempty"`
	Channel     string   `toml:"channel,omitempty"`
	AuthUser    string   `toml:"authUser,omitempty"`
	NotifyUsers []string `toml:"notifyUsers,omitempty"`
}

func Load(pathToToml string) (*Config, error) {
	var conf Config
	if _, err := toml.DecodeFile(pathToToml, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
