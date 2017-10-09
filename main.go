package main

import (
	"github.com/mitchellh/go-homedir"
	"os"
	"fmt"
)

var (
	homeDirPath, _ = homedir.Dir()
)

func main() {
	// 引数の確認
	argsLength := len(os.Args)
	if argsLength != 2 {
		fmt.Println("引数の数が間違っています")
		os.Exit(1)
	}

	// 設定ファイル読み込み
	config := NewConfig()

	// History構造体を生成
	h := History {
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

	// Termboxの表示
	t := NewTermbox(commandHistory)
	err = t.Init()
	if err != nil {
		panic("termboxのinitializeに失敗")
	}
	t.Display()

	// マクロ生成
	macroName := os.Args[1]
	macro := NewMacro(macroName, t.Selection, config.ShellType, config.ConfigFilePath)
	macro.SaveFile()
	fmt.Println("Create Macro!!")
}
