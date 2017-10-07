package main

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
)

var (
	configDirPath = filepath.Join(homeDirPath, ".mino")
	configFilePath = filepath.Join(configDirPath, "config.toml")
)

type Config struct {
	ShellType string
	HistoryFilePath string
}

func NewConfig() *Config {
	var config Config

	if !existConfigDir() {
		// configディレクトリが無かった場合のエラー
		panic("Error: .minoがない")
	}

	if !existConfigFile() {
		// configファイルが無かった場合のエラー
		panic("Error: config.tomlがない")
	}

	_, err := toml.DecodeFile(configFilePath, &config)
	if err != nil {
		// configファイルのデコードに失敗
		panic("config.tomlのデコードに失敗")
	}

	config.HistoryFilePath, err = homedir.Expand(config.HistoryFilePath)
	if err != nil {
		panic("homedirが含まれている場合の処理に失敗")
	}

	return &config
}

func existConfigDir() bool {
	_, err := os.Stat(configDirPath)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func existConfigFile() bool {
	_, err := os.Stat(configFilePath)
	if os.IsNotExist(err) {
		return false
	}

	return true
}
