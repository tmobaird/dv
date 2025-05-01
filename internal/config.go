package internal

import (
	"fmt"
	"io"
	"io/fs"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Context string `yaml:"context"`
}

func FileExists(fileSystem fs.FS) bool {
	files, err := fs.ReadDir(fileSystem, ".")
	fmt.Println(files)
	if err != nil {
		return false
	}

	found := false

	for _, file := range files {
		if file.Name() == "config.yaml" {
			found = true
		}
	}

	return found
}

func DefaultConfig() Config {
	return Config{Context: "main"}
}

func Read(fileSystem fs.FS) (Config, error) {
	if FileExists(fileSystem) {
		yamlFile, err := fs.ReadFile(fileSystem, "config.yaml")
		if err != nil {
			return Config{}, err
		}

		var config Config
		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			return Config{}, err
		}

		return config, nil
	}
	
	return DefaultConfig(), nil
}

func Save(writer io.Writer, config Config) error {
	output, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	_, err = writer.Write(output)
	return err
}
