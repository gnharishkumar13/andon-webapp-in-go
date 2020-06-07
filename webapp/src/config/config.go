package config

import "os"

type Config struct {
	StaticRoot string
}

var c Config

func init() {
	setupConfig()
}

func setupConfig() {
	c.StaticRoot = os.Getenv("STATIC_ROOT")
}

//Get configs
func Get() Config {
	return c
}
