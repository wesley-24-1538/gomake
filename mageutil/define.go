package mageutil

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

var (
	serviceBinaries             map[string]int
	toolBinaries                []string
	MaxFileDescriptors          int
	subdirectoryOfCmdForCompile []string
)

type Config struct {
	ServiceBinaries    map[string]int `yaml:"serviceBinaries"`
	ToolBinaries       []string       `yaml:"toolBinaries"`
	MaxFileDescriptors int            `yaml:"maxFileDescriptors"`
}

type specifyConfig struct {
	ServiceBinaries []string `yaml:"serviceBinaries"`
}

func InitForSSC() {
	yamlFile, err := ioutil.ReadFile("start-config.yml")
	if err != nil {
		fmt.Printf("error reading YAML file: %v", err)
		os.Exit(1)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Printf("error unmarshalling YAML: %v", err)
		os.Exit(1)
	}

	adjustedBinaries := make(map[string]int)
	for binary, count := range config.ServiceBinaries {
		if !isSpecifyServiceBinary(binary) {
			continue
		}
		if runtime.GOOS == "windows" {
			binary += ".exe"
		}
		adjustedBinaries[binary] = count
	}
	serviceBinaries = adjustedBinaries
	toolBinaries = config.ToolBinaries
	MaxFileDescriptors = config.MaxFileDescriptors
}

func isSpecifyServiceBinary(binary string) bool {
	if len(subdirectoryOfCmdForCompile) > 0 {
		return slices.Contains(subdirectoryOfCmdForCompile, binary)
	} else {
		return true
	}
}

func IfSpecifyDirectoryOfCmd(path string) bool {
	if len(subdirectoryOfCmdForCompile) == 0 {
		return true
	}
	dirPath := filepath.Dir(path)
	if !strings.Contains(path, "cmd") { //不是cmd目录下的目录不处理
		return true
	}
	tmp := strings.Split(dirPath, "/")
	return slices.Contains(subdirectoryOfCmdForCompile, tmp[len(tmp)-1])
}

func init() {
	//获取文件配置
	yamlFile, err := ioutil.ReadFile("specify-servers.yml")
	if err != nil {
		fmt.Printf("file:define func:init() error reading YAML: %v", err)
		return
	}
	var config specifyConfig
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Printf("file:define func:init() error unmarshalling YAML: %v", err)
		os.Exit(1)
	}
	fmt.Printf("%+v \n", config)
	subdirectoryOfCmdForCompile = config.ServiceBinaries
}
