package yaml_config

import (
	"errors"
	"os"

	"github.com/mntwo/tasklab/internal/configer"
	"gopkg.in/yaml.v3"
)

var _ configer.Configer = (*YamlConfig)(nil)

type YamlConfig struct {
	configFile string
	data       interface{}
}

func New(opts ...Option) *YamlConfig {
	cfg := defaultConfig()

	for _, opt := range opts {
		opt.apply(cfg)
	}

	return &YamlConfig{
		configFile: cfg.configFile,
		data:       cfg.configData,
	}
}

func (y *YamlConfig) Parse() error {
	if y.configFile == "" {
		return errors.New("config file is empty")
	}
	dataBytes, err := os.ReadFile(y.configFile)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(dataBytes, y.data)
}
