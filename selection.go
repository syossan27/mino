package main

// Selectionは表示されているコマンド履歴の選択の責務を負う
// Offsetは「選択されている位置」を担う

type (
	// Index: カーソルが指しているコマンド履歴配列の添字
	// LatestExecOrdeer: 採番された最新の実行順序番号
	// Selected: 選択済みのコマンド履歴情報配列
	Selection struct {
		Index 			int
		LatestExecOrder int
		Command
		Selected        []Info
	}

	// ExecOrder： 実行順序番号
	Info struct {
		Index     int
		ExecOrder int
		Command   string
	}
)

const (
	// Offsetの初期値
	defaultIndex			= 1
	defaultLatestExecOrder 	= 1
)

func NewSelection() Selection {
	return Selection {
		Index: 				defaultIndex,
		LatestExecOrder: 	defaultLatestExecOrder,
	}
}

func (s *Selection) Select() {
	// 選択済みだった場合、選択を外す
	if s.DeleteIfExistSelected() {
		return
	}

	command := s.Command

	info := Info{
		Index:     command.Index,
		ExecOrder: s.LatestExecOrder,
		Command:   command.Content,
	}
	s.Selected = append(s.Selected, info)
	s.LatestExecOrder++
}

// 返り値がInfoなら選択済み
// nilなら選択されていない
func (s *Selection) GetSelectedNumber(index int) (int, bool) {
	for _, info := range s.Selected {
		if info.Index == index {
			return info.ExecOrder, true
		}
	}

	return 0, false
}

func (s *Selection) DeleteIfExistSelected() bool {
	selectedIndex, exist := s.IsSelected()
	if exist {
		s.Selected = deleteSelected(s.Selected, selectedIndex)
		s.ReNumbering(selectedIndex)
		s.LatestExecOrder--
		return true
	}

	return false
}

func (s *Selection) IsSelected() (int, bool) {
	for key, info := range s.Selected {
		if info.Index == s.Command.Index {
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
	s.LatestExecOrder = defaultLatestExecOrder
}

// InfoのNumberを再採番
// 削除された選択済みコマンド以降のコマンドの選択番号をデクリメントする
func (s *Selection) ReNumbering(selectedIndex int) {
	for i := selectedIndex ; i < len(s.Selected); i++ {
		s.Selected[i].ExecOrder--
	}
}

func (s *Selection) ClearIndex() {
	s.Index = defaultIndex
}
