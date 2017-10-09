package main

import (
	"io/ioutil"
	"path/filepath"
	"os"
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

func NewMacro(name string, s Selection, shellType string, configFilePath string) *Macro {
	selectedInfo := s.Selected
	return &Macro {
		Name: name,
		Info: selectedInfo,
		Shell: Shell {
			Type: shellType,
			ConfigFilePath: configFilePath,
		},
	}
}

func (m *Macro) SaveFile() {
	if !existMacroDir() {
		// macroディレクトリが無かった場合のエラー
		panic("Error: macroディレクトリがない")
	}

	content := m.getFileHeader()
	for _, info := range m.Info {
		 content = content + info.Command + "\n"
	}
	ioutil.WriteFile(filepath.Join(configMacroDirPath, m.Name + ".sh"), []byte(content), os.ModePerm)
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
