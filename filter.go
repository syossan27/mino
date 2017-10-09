package main

import "strings"

type Filter struct {
	SearchQuery []rune
}

func NewFilter() Filter {
	return Filter {}
}

func (f *Filter) FilterResult(commandHistory []Command) []Command {
	var filterResult []Command
	for _, command := range commandHistory {
		filterIndex := f.getIndex(command.Content)
		if filterIndex == -1 {
			continue
		}
		command.FilterIndex = filterIndex
		filterResult = append(filterResult, command)
	}
	return filterResult
}

func (f *Filter) getIndex(command string) int {
	SearchQueryToString := string(f.SearchQuery)
	index := strings.Index(command, SearchQueryToString)
	return index
}

func (f *Filter) Append(c rune) {
	f.SearchQuery = append(f.SearchQuery, c)
}

func (f *Filter) Delete() {
	f.SearchQuery = f.SearchQuery[:len(f.SearchQuery) - 1]
}

func (f *Filter) Length() int {
	return len(f.SearchQuery)
}
