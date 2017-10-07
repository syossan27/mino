package main

import(
	"path/filepath"
	"os"
	"bufio"
	"regexp"
)

var (
	zshRegExp = regexp.MustCompile(`:\s[0-9]*:[0-9]*;`)
)

type History struct {
	ShellType string
	HistoryFilePath string
}

func (h *History) Load() ([]string, error) {
	// historyファイルのパスが設定されているなら、その値を使用する
	// 設定されていない場合には、それぞれのshellのデフォルトのパスを使用する
	if h.HistoryFilePath == "" {
		historyFilePath, err := h.getHistoryFilePath()
		if err != nil {
			// shellTypeが読み込めなかったエラー
			return nil, err
		}
		h.HistoryFilePath = historyFilePath
	}

	commandHistory, err := h.decodeHistoryFile()
	if err != nil {
		// historyファイルがデコード出来なかったエラー
		return nil, err
	}

	return commandHistory, nil
}

func (h History) getHistoryFilePath() (string, error) {
	// Shellの種類からhistoryファイルの場所を返す
	var historyFilePath string

	switch h.ShellType {
	case "bash":
		historyFilePath = filepath.Join(homeDirPath, ".local/share/fish/fish_history")
	case "zsh":
		historyFilePath = filepath.Join(homeDirPath, ".zsh_history")
	case "fish":
		historyFilePath = filepath.Join(homeDirPath, ".bash_history")
	}

	return historyFilePath, nil
}

func (h History) decodeHistoryFile() ([]string, error) {
	// コマンド履歴一覧を返す
	var commandHistory []string

	fp, err := os.Open(h.HistoryFilePath)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		// zsh対応
		scanText := zshRegExp.ReplaceAllString(scanner.Text(), "")
		commandHistory = append(commandHistory, scanText)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return commandHistory, nil
}