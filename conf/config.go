package conf

import "github.com/BurntSushi/toml"

type Config struct {

}

func LoadConfigFromFile(filePath string) (*Config, error) {
	var config Config
	_, err := toml.DecodeFile(filePath, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
