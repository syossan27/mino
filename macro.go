package main

import (
	"io/ioutil"
	"path/filepath"
	"os"
	"github.com/pkg/errors"
)

var (
	configMacroDirPath = filepath.Join(homeDirPath, ".mino", "macro")
)

type (
	Macro struct {
		Name string
		Info []Info
		Shell
	}

	Shell struct {
		Type string
		ConfigFilePath string
	}
)

func NewMacro(selection Selection, shellType string, configFilePath string) *Macro {
	selectedInfo := selection.Selected
	return &Macro {
		Name: os.Args[1],
		Info: selectedInfo,
		Shell: Shell {
			Type: shellType,
			ConfigFilePath: configFilePath,
		},
	}
}

func (m *Macro) SaveFile() error {
	if !existMacroDir() {
		// macroディレクトリが無かった場合のエラー
		return errors.New("not found macro directory")
	}

	content := m.getFileHeader()
	for _, info := range m.Info {
		 content = content + info.Command + "\n"
	}
	ioutil.WriteFile(filepath.Join(configMacroDirPath, m.Name + ".sh"), []byte(content), os.ModePerm)

	return nil
}

func existMacroDir() bool {
	_, err := os.Stat(configMacroDirPath)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func (m *Macro) getFileHeader() string {
	var fileHeader string
	switch m.Shell.Type {
	case "bash":
		fileHeader = "#!/bin/bash\n"
	case "zsh":
		fileHeader = "#!/bin/zsh\n"
	}

	if m.Shell.ConfigFilePath != "" {
		fileHeader += "source " + m.Shell.ConfigFilePath + "\n"
	}

	return fileHeader
}
