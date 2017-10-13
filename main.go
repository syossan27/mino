package main

import (
	"errors"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/constabulary/gb/testdata/src/h"
)

var (
	homeDirPath, _ = homedir.Dir()
)

const (
	// ExitCodeOK is exit code for OK
	ExitCodeOK = iota
	// ExitCodeError is exit code for Error
	ExitCodeError
)

func main() {
	// 引数の確認
	argsLength := len(os.Args)
	if argsLength != 2 {
		Fatal(errors.New("missing argument"))
	}

	// 設定ファイル読み込み
	config, err := NewConfig()
	if err != nil {
		Fatal(err)
	}

	// History構造体を生成
	h := History{
		ShellType:       config.ShellType,
		HistoryFilePath: config.HistoryFilePath,
	}

	// historyファイルを読み込み、コマンド履歴を受け取る
	commandHistory, err := h.Load()
	if err != nil {
		Fatal(err)
	}

	// Termboxの表示
	t := NewTermbox(commandHistory)
	err = t.Init()
	if err != nil {
		Fatal(errors.New("failed initialize termbox"))
	}
	t.Display()

	// マクロ生成
	macroName := os.Args[1]
	macro := NewMacro(macroName, t.Selection, config.ShellType, config.ConfigFilePath)
	macro.SaveFile()
	Success("Create Macro!")

	os.Exit(ExitCodeOK)
}
