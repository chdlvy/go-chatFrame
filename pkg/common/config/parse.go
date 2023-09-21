package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"runtime"
)

var (
	//get curFilePath
	_, curFilePath, _, _ = runtime.Caller(0)
	//root is the root directory of this project
	Root = filepath.Join(filepath.Dir(curFilePath), "../../../")
)

const (
	DefaultFolderPath = "./config/"
)

// configFolderPath is relation path and start from root directory of this project
func InitConfig(configName, configFolderPath string) error {
	if configFolderPath == "" {
		configFolderPath = DefaultFolderPath
	}

	configpath := filepath.Join(Root, configFolderPath, configName)
	fmt.Println("configpath:", configpath)
	_, err := os.Stat(configpath)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(configpath)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(data, &Config); err != nil {
		return err
	}
	return nil
}
