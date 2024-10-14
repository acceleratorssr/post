package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

type Info struct {
	Config *Config
}

func FindFirstYAMLFile() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	fmt.Println(execPath)
	searchPath := ""
	if !strings.Contains(execPath, "output") {
		searchPath = "./user/config"
	} else {
		searchPath = "../user/config"
	}

	var yamlFile string
	err = filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".yaml" {
			yamlFile = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	if yamlFile == "" {
		return "", fmt.Errorf("no yaml file found")
	}

	return yamlFile, nil
}

func InitConfig() *Info {
	c := &Config{}

	yamlFile, err := FindFirstYAMLFile()
	if err != nil {
		panic(fmt.Errorf("find yaml file error: %v", err))
	}

	yamlConf, err := os.ReadFile(yamlFile)
	if err != nil {
		panic(fmt.Errorf("read yaml error: %v\n", err))
	}

	err = yaml.Unmarshal(yamlConf, c)
	if err != nil {
		panic(fmt.Errorf("config init unmarshal: %v\n", err))
	}

	return &Info{
		Config: c,
	}
}
