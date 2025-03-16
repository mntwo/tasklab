package yaml_config

type Option interface {
	apply(cfg *optconfig)
}

type option func(cfg *optconfig)

func (fn option) apply(cfg *optconfig) {
	fn(cfg)
}

type optconfig struct {
	configFile string
	configData interface{}
}

func defaultConfig() *optconfig {
	return &optconfig{
		configFile: "config.yaml",
		configData: struct{}{},
	}
}

func WithConfigFile(file string) Option {
	return option(func(cfg *optconfig) {
		cfg.configFile = file
	})
}

func WithConfigData(data interface{}) Option {
	return option(func(cfg *optconfig) {
		cfg.configData = data
	})
}
