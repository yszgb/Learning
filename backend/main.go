package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"learning.com/learning/app"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config-dev.yaml", "Server config file path")
	flag.Parse() // flag 函数作用：解析命令行参数，使命令行参数赋值给 configPath

	currnetDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	if !strings.HasPrefix(configPath, "/") {
		// configPath 不以 '/' 开头时，说明是相对路径，需要加上当前目录
		configPath = filepath.Join(currnetDir, configPath)
	}

	app, err := app.New(configPath, currnetDir)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	err = app.Run()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// configPath := filepath.Join(currnetDir, "config-dev.yaml")
	// config, err := models.NewConfig(configPath, currnetDir)
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// 	return
	// }
	// fmt.Println(config)
}
