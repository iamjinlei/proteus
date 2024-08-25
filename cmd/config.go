package main

type Config struct {
	Entry  string            `yaml:"entry"`
	Assets map[string]string `yaml:"assets"`
}
