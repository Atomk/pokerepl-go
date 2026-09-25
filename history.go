package main

import "fmt"

type History struct {
	_list         []string
	_currentIndex int
}

func NewHistory() History {
	return History{
		_list: make([]string, 0, 50),
	}
}

func (h History) Add(entry string) {
	h._list = append(h._list, entry)
}

func (h History) Previous() string {
	if len(h._list) == 0 || h._currentIndex == 0 {
		return ""
		//return "", fmt.Errorf("no previous items in history")
	}
	h._currentIndex--
	return h._list[h._currentIndex]
}

func (h History) Next() (string, error) {
	count := h.Count()
	if count == 0 || h._currentIndex+1 >= count {
		return "", fmt.Errorf("already on most recent entry")
	}
	h._currentIndex++
	return h._list[h._currentIndex], nil
}

func (h History) Count() int {
	return len(h._list)
}
