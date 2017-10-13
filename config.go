package main

import (
	"errors"
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
	ConfigFilePath string
}

func NewConfig() (*Config, error) {
	var config Config

	if !existConfigDir() {
		// configディレクトリが無かった場合のエラー
		return nil, errors.New("not found config directory")
	}

	if !existConfigFile() {
		// configファイルが無かった場合のエラー
		return nil, errors.New("not found config file")
	}

	_, err := toml.DecodeFile(configFilePath, &config)
	if err != nil {
		return nil, errors.New("failed decode config file")
	}

	if config.HistoryFilePath == "" {
		return nil, errors.New("not setting history file path")
	}

	config.HistoryFilePath, err = homedir.Expand(config.HistoryFilePath)
	if err != nil {
		return nil, err
	}

	// 設定ファイルの指定がない場合、早期リターン
	if config.ConfigFilePath == "" {
		return &config, nil
	}

	config.ConfigFilePath, err = homedir.Expand(config.ConfigFilePath)
	if err != nil {
		return nil, err
	}

	return &config, nil
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
