package main

// Selectionは表示されているコマンド履歴の選択の責務を負う
// Offsetは「選択されている位置」を担う

type Selection struct {
	Offset int
}

// Offsetの初期値
const initSelectionOffset = 1

func NewSelection() Selection {
	return Selection {
		Offset: initSelectionOffset,
	}
}
