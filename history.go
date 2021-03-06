package main

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"regexp"
)

var (
	zshRegExp = regexp.MustCompile(`:\s[0-9]*:[0-9]*;`)
)

type (
	History struct {
		ShellType string
		HistoryFilePath string
	}

	// Index: コマンド履歴配列の添字
	// （検索条件で再生成した際に元の添字を知るのに必要）
	// Content: コマンド文字列
	// FilterIndex: 検索に引っかかった位置
	Command struct {
		Index 		int
		Content 	string
		FilterIndex int
	}
)

func NewHistory(shellType, historyFilePath string) History {
	return History {
		ShellType: shellType,
		HistoryFilePath: historyFilePath,
	}
}

func (h *History) Load() ([]Command, error) {
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
	default:
		return "", errors.New("invalid shell type")
	}

	return historyFilePath, nil
}

func (h History) decodeHistoryFile() ([]Command, error) {
	// コマンド履歴一覧を返す
	var commandHistory []string

	fp, err := os.Open(h.HistoryFilePath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		command := zshRegExp.ReplaceAllString(scanner.Text(), "") // zsh対応
		command = parseNonAscii([]byte(command))
		commandHistory = append(commandHistory, command)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// コマンド履歴を最新から表示するため逆順にする
	var reverseCommandHistory []Command
	for i := len(commandHistory) - 1; i >= 0; i-- {
		command := Command {
			Index: i,
			Content: commandHistory[i],
		}

		reverseCommandHistory = append(reverseCommandHistory, command)
	}

	return reverseCommandHistory, nil
}

func parseNonAscii(latin1Byte []byte) string {
	isMarking := false
	var byteBuffer []byte

	for _, codePoint := range latin1Byte {
		// 131は0x83の10進数表現
		if codePoint == 131 {
			isMarking = true
			continue
		}

		if isMarking {
			// 6bit目を反転させるために
			// 0x20をXORする
			invertCodePoint := codePoint ^ 32
			byteBuffer = append(byteBuffer, invertCodePoint)
			isMarking = false
		} else {
			byteBuffer = append(byteBuffer, codePoint)
		}
	}

	return string(byteBuffer)
}
