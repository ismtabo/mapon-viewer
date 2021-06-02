package main

import "github.com/ismtabo/mapon-viewer/pkg/cfg"

// Config for the application configuration.
type Config struct {
	Server struct {
		Host string `yaml:"host" envconfig:"SERVER_HOST"`
		Port string `yaml:"port" envconfig:"SERVER_PORT"`
	} `yaml:"server"`
	Log struct {
		Level string `yaml:"level" envconfig:"LOG_LEVEL"`
	} `yaml:"log"`
	Mongo       cfg.MongoConfig   `yaml:"mongo"`
	Session     cfg.SessionConfig `yaml:"session"`
	MaponConfig cfg.MaponConfig   `yaml:"mapon"`
}
