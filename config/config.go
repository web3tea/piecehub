package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server ServerConfig `toml:"server"`
	Disks  []DiskConfig `toml:"disks"`
	S3s    []S3Config   `toml:"s3s"`
}

type ServerConfig struct {
	Address      string `toml:"address"`
	ReadTimeout  int    `toml:"read_timeout"`
	WriteTimeout int    `toml:"write_timeout"`
}

type DiskConfig struct {
	Name     string `toml:"name"`
	RootDir  string `toml:"root_dir"`
	MaxSize  int64  `toml:"max_size"`
	DirectIO bool   `toml:"direct_io"`
}

type S3Config struct {
	Name      string `toml:"name"`
	Endpoint  string `toml:"endpoint"`
	Region    string `toml:"region"`
	Bucket    string `toml:"bucket"`
	AccessKey string `toml:"access_key"`
	SecretKey string `toml:"secret_key"`
}

var DefaultConfig = Config{
	Server: ServerConfig{
		Address:      ":8080",
		ReadTimeout:  600,
		WriteTimeout: 600,
	},
}

func LoadConfig(path string) (*Config, error) {
	config := DefaultConfig

	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateConfig(cfg *Config) error {
	names := make(map[string]bool)

	for _, disk := range cfg.Disks {
		if disk.Name == "" {
			return fmt.Errorf("disk name cannot be empty")
		}
		if names[disk.Name] {
			return fmt.Errorf("duplicate storage name: %s", disk.Name)
		}
		names[disk.Name] = true
	}

	for _, s3 := range cfg.S3s {
		if s3.Name == "" {
			return fmt.Errorf("s3 name cannot be empty")
		}
		if names[s3.Name] {
			return fmt.Errorf("duplicate storage name: %s", s3.Name)
		}
		names[s3.Name] = true
	}

	return nil
}
