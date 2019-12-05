package main

var config *Config

// Config global config
type Config struct {
	Port     int
	Username string
	Password string
	Server   string
}
