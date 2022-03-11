// Package config provides methods for accesing config file in TOML format
package utils

import (
	"os"

	"github.com/naoina/toml"
)

type Config struct {
	Debug bool
	Kitsu struct {
		Hostname string
		Email    string
		Password string
	}
	Backup struct {
		Threads         int
		PollDuration    int
		LocalStorage    string
		IgnoreExtension []string
		FastDelete      bool
		S3              struct {
			AccessKey        string
			SecretKey        string
			BucketName       string
			Endpoint         string
			Region           string
			S3ForcePathStyle bool
			RootFolderName   string
		}
	}
}

func ConfRead() Config {
	path := "conf.toml"
	if os.Getenv("TEST") == "true" {
		path = os.Getenv("CONF_PATH")
	}

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var config Config
	if err := toml.NewDecoder(f).Decode(&config); err != nil {
		panic(err)
	}

	return config
}
