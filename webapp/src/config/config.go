package config

import "os"

const (
	Environment = "development"
)

//Config struct
type Config struct {
	StaticRoot string
	ViewRoot   string
	Environment string
}


var c Config

func init() {
	setupConfig()
}

func setupConfig() {
	c.StaticRoot = os.Getenv("STATIC_ROOT")
	c.ViewRoot = os.Getenv("VIEW_ROOT")
	c.Environment = os.Getenv("ENV")
}

//Get configs
func Get() Config {
	return c
}
