package main

import (
	"github.com/mitchellh/go-homedir"
)

var (
	homeDirPath, _ = homedir.Dir()
)

const (
	ExitCodeError = 1
)

func main() {
	// 引数の確認
	err := ValidateArgs()
	if err != nil {
		Fatal(err)
	}

	// 設定ファイル読み込み
	config, err := NewConfig()
	if err != nil {
		Fatal(err)
	}

	// historyファイルを読み込み、コマンド履歴を受け取る
	history := NewHistory(config.ShellType, config.HistoryFilePath)
	commands, err := history.Load()
	if err != nil {
		Fatal(err)
	}

	// Termboxの表示
	t := NewTermbox(commands)
	err = t.Init()
	if err != nil {
		Fatal(err)
	}
	t.Display()

	// マクロ生成
	macro := NewMacro(t.Selection, config.ShellType, config.ConfigFilePath)
	macro.SaveFile()
	Success("Create Macro!")
}
