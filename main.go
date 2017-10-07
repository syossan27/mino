package main

import (
	"github.com/mitchellh/go-homedir"
)

var (
	homeDirPath, _ = homedir.Dir()
)

func main() {
	// 設定ファイル読み込み
	config := NewConfig()

	h := History{
		ShellType: config.ShellType,
		HistoryFilePath: config.HistoryFilePath,
	}

	// historyファイルを読み込み、コマンド履歴を受け取る
	commandHistory, err := h.Load()
	if err != nil {
		// エラー表示
		// historyファイルの読み込みが大前提なので、exitさせる
		panic("historyファイルの読み込み失敗")
	}

	var reverseCommandHistory []string
	for i := len(commandHistory) - 1; i > 0; i-- {
		command := commandHistory[i]
		reverseCommandHistory = append(reverseCommandHistory, command)
	}

	t := NewTermbox()
	t.Init()
	t.Do(reverseCommandHistory)
}
