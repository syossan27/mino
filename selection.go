package main

// Selectionは表示されているコマンド履歴の選択の責務を負う
// Offsetは「選択されている位置」を担う

// Number: 選択番号
// Index: commandHistoryのindex
type (
	Selection struct {
		Offset                int
		Selected              []Info
		currentSelectedNumber int
	}

	Info struct {
		Index int
		Number int
		Command string
	}
)

const (
	// Offsetの初期値
	initSelectionOffset       = 1
	initCurrentSelectedNumber = 1
)

func NewSelection() Selection {
	return Selection {
		Offset:                initSelectionOffset,
		currentSelectedNumber: initCurrentSelectedNumber,
	}
}

func (s *Selection) Select(commandHistory []string) {
	// 選択済みだった場合、選択を外す
	if s.DeleteIfExistSelected() {
		return
	}

	info := Info{
		Index: s.Offset,
		Number: s.currentSelectedNumber,
		Command: commandHistory[s.Offset - 1],
	}
	s.Selected = append(s.Selected, info)
	s.currentSelectedNumber++
}

// 返り値がInfoなら選択済み
// nilなら選択されていない
func (s *Selection) GetSelectedIndex(commandIndex int) (int, bool) {
	for _, info := range s.Selected {
		if info.Index == commandIndex {
			return info.Number, true
		}
	}

	return 0, false
}

func (s *Selection) DeleteIfExistSelected() bool {
	selectedIndex, exist := s.IsSelected()
	if exist {
		s.Selected = deleteSelected(s.Selected, selectedIndex)
		s.ReNumbering(selectedIndex)
		s.currentSelectedNumber--
		return true
	}

	return false
}

func (s *Selection) IsSelected() (int, bool) {
	for key, info := range s.Selected {
		if info.Index == s.Offset {
			return key, true
		}
	}

	return 0, false
}

func deleteSelected(a []Info, i int) []Info {
	copy(a[i:], a[i+1:])
	a = a[:len(a)-1]
	return a
}

func (s *Selection) ClearSelected() {
	s.Selected = []Info{}
	s.currentSelectedNumber = initCurrentSelectedNumber
}

// InfoのNumberを再採番
// 削除された選択済みコマンド以降のコマンドの選択番号をデクリメントする
func (s *Selection) ReNumbering(selectedIndex int) {
	for i := selectedIndex ; i < len(s.Selected); i++ {
		s.Selected[i].Number--
	}
}
