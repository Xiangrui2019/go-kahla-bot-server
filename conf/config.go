package conf

import (
	"github.com/BurntSushi/toml"
	"gopkg.in/macaron.v1"
)

type Config struct {
	Host string `toml:"HOST"`
	Port int `toml:"PORT"`
	Env string `toml:"ENV"`
	BotConfig Bot `toml:"bot"`
	CorsOriginURL string `toml:"CORS_ORIGIN_URL"`
	MySqlDSN string `toml:"MYSQL_DSN"`
	TokenLength int `toml:"TOKEN_LENGTH"`
}

type Bot struct {
	Name string `toml:"name"`
	Email string `toml:"email"`
	Password string `toml:"password"`
	KahlaServer string `toml:"kahlaserver"`
	CallbackServer string `toml:"callbackserver"`
	MessageCallbackEndpoint string `toml:"messagecallbackendpoint"`
}

func (config *Config) ConfigEnvironment() error {
	switch config.Env {
	case "dev":
		macaron.Env = macaron.DEV
	case "test":
		macaron.Env = macaron.TEST
	case "prod":
		macaron.Env = macaron.PROD
	}

	return nil
}

func LoadConfigFromFile(filePath string) (*Config, error) {
	var config Config
	_, err := toml.DecodeFile(filePath, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
