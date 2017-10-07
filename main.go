package main

import (
	"github.com/k0kubun/pp"
	"github.com/mitchellh/go-homedir"
	"fmt"
)

var (
	homeDirPath, _ = homedir.Dir()
)

func main() {
	// 設定ファイル読み込み
	config := NewConfig()

	pp.Println(config)

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
	for _, v := range commandHistory {
		fmt.Println(v)
	}

	t := NewTermbox()
	t.Init()
	t.Do()
}
