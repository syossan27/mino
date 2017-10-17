package main

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli"
)

func ExecMacro(c *cli.Context) error {
	// 引数の確認
	err := ValidateArgs(c.Args())
	if err != nil {
		Fatal(err)
	}
	macro := c.Args().First()

	// 設定ファイル読み込み
	config, err := NewConfig()
	if err != nil {
		Fatal(err)
	}

	// マクロの実行
	out, err := exec.Command(config.ShellType, "-c", "'" + filepath.Join(configMacroDirPath, macro + ".sh") +  "'").Output()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", out)

	return nil
}
