package core

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Context       string `yaml:"context"`
	HideCompleted bool   `yaml:"hideCompleted"`
}

func (config Config) IsBlank() bool {
	return config.Context == "" && !config.HideCompleted
}

const FILENAME = "config.yaml"
const LOG_FILE_TIME_FORMAT = "20060102"

func BasePath() string {
	usr, _ := user.Current()
	dirname := filepath.Join(usr.HomeDir, ".dv")
	if os.Getenv("DV_BASE_PATH") != "" {
		dirname = os.Getenv("DV_BASE_PATH")
	}

	return dirname
}

func ConfigFilePath() string {
	return fmt.Sprintf("%s/%s", BasePath(), FILENAME)
}

func TodoFilePath(context string) string {
	return fmt.Sprintf("%s/lists/%s.md", BasePath(), context)
}

func ScheduleDirectoryPath() string {
	return fmt.Sprintf("%s/schedules", BasePath())
}

func ScheduleFilePath(d time.Time) string {
	return fmt.Sprintf("%s/%s.md", ScheduleDirectoryPath(), d.Format(time.DateOnly))
}

func LogDirectoryPath() string {
	return filepath.Join(BasePath(), "logs")
}

func LogFileName(t time.Time) string {
	return fmt.Sprintf("%s.md", t.Format(LOG_FILE_TIME_FORMAT))
}

func LogFilePath(t time.Time) string {
	return filepath.Join(LogDirectoryPath(), LogFileName(t))
}

func LogFileExists(d time.Time) bool {
	return FileExists(os.DirFS(LogDirectoryPath()), fmt.Sprintf("%s.md", d.Format(LOG_FILE_TIME_FORMAT)))
}

func ConfigFileExists(filesystem fs.FS) bool {
	return FileExists(filesystem, FILENAME)
}

func FileExists(filesystem fs.FS, filename string) bool {
	files, err := fs.ReadDir(filesystem, ".")
	if err != nil {
		return false
	}

	found := false

	for _, file := range files {
		if file.Name() == filename {
			found = true
		}
	}

	return found
}

func DefaultConfig() Config {
	return Config{Context: "main", HideCompleted: false}
}

func ReadConfig(fileSystem fs.FS) (Config, error) {
	if ConfigFileExists(fileSystem) {
		yamlFile, err := fs.ReadFile(fileSystem, FILENAME)
		if err != nil {
			return Config{}, err
		}

		var config Config
		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			return Config{}, err
		}

		if config.IsBlank() {
			config = DefaultConfig()
		}

		return config, nil
	}

	return DefaultConfig(), nil
}

func SaveConfig(writer io.Writer, config Config) error {
	output, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	_, err = writer.Write(output)
	return err
}

func PersistConfig(config Config) error {
	file, err := os.OpenFile(ConfigFilePath(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = SaveConfig(file, config)
	return err
}
