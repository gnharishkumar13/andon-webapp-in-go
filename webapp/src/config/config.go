package config

import (
	"fmt"
	"os"
)

const (
	Environment = "development"
)

//Config struct
type Config struct {
	StaticRoot  string
	ViewRoot    string
	Environment string
	Salt        string
}

var c Config

func init() {
	setupConfig()
}

func setupConfig() {
	c.StaticRoot = os.Getenv("STATIC_ROOT")
	c.ViewRoot = os.Getenv("VIEW_ROOT")
	c.Environment = os.Getenv("ENV")
	c.Salt = os.Getenv("salt")
	fmt.Println("salt", c.Salt)
}

//Get configs
func Get() Config {
	return c
}
