package config

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Log      LogConfig      `yaml:"log"`
	Database DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type LogConfig struct {
	Level string `yaml:"level"`
}

type ServerConfig struct {
	Port   string `yaml:"port"`
	Host   string `yaml:"host"`
	Domain string `yaml:"domain"`
}
