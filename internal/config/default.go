package config

import (
	"flag"
	"time"

	"github.com/mntwo/tasklab/internal/configer/yaml_config"
)

var defaultConfig *Config
var configFile = flag.String("i", "config.yaml", "the application config file, eg: -i config.yaml, default: config.yaml")

func init() {
	flag.Parse()
	defaultConfig = New()
}

func New() *Config {
	c := &Config{}
	configer := yaml_config.New(
		yaml_config.WithConfigFile(*configFile),
		yaml_config.WithConfigData(c),
	)
	err := configer.Parse()
	if err != nil {
		panic(err)
	}
	return c
}

func GetApplication() *Application {
	if defaultConfig == nil {
		return nil
	}
	return defaultConfig.Application
}

func GetHttpServers() []*HttpServer {
	if defaultConfig == nil {
		return nil
	}
	return defaultConfig.HttpServer
}

func GetHttpServer(name string) *HttpServer {
	if defaultConfig == nil {
		return nil
	}
	for _, s := range defaultConfig.HttpServer {
		if s.Name == name {
			return s
		}
	}
	return nil
}

func GetLog() *Log {
	if defaultConfig == nil {
		return nil
	}
	return defaultConfig.Log
}

func GetPostgres() []*Postgres {
	if defaultConfig == nil {
		return nil
	}
	return defaultConfig.Postgres
}

type Config struct {
	Application *Application  `json:"application" yaml:"application"`
	HttpServer  []*HttpServer `json:"http_server" yaml:"http_server"`
	Log         *Log          `json:"log" yaml:"log"`
	Postgres    []*Postgres   `json:"postgres" yaml:"postgres"`
}

type Application struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
	Env     string `json:"env" yaml:"env"`
}

type HttpServer struct {
	Name         string        `json:"name" yaml:"name"`
	Addr         string        `json:"addr" yaml:"addr"`
	CloseTimeout time.Duration `json:"close_timeout" yaml:"close_timeout"`
	WriteTimeout time.Duration `json:"write_timeout" yaml:"write_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout" yaml:"read_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout" yaml:"idle_timeout"`
}

type Log struct {
	Type       string `json:"type" yaml:"type"`
	Level      string `json:"level" yaml:"level"`
	MaxSize    int    `json:"max_size" yaml:"max_size"`
	MaxAge     int    `json:"max_age" yaml:"max_age"`
	MaxBackups int    `json:"max_backups" yaml:"max_backups"`
	LocalTime  bool   `json:"local_time" yaml:"local_time"`
	Compress   bool   `json:"compress" yaml:"compress"`
}

type Postgres struct {
	Name            string        `json:"name" yaml:"name"`
	Host            string        `json:"host" yaml:"host"`
	Port            int           `json:"port" yaml:"port"`
	User            string        `json:"user" yaml:"user"`
	Password        string        `json:"password" yaml:"password"`
	Database        string        `json:"database" yaml:"database"`
	MaxOpenConns    int           `json:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns" yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime" yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time" yaml:"conn_max_idle_time"`
}
