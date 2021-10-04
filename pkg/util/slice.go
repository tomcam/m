package util

import (
	"sort"
)

/* Searching a sorted slice is fast.
   This tracks whether the slice has been sorted
   and sorts it on first search.
*/
type SearchInfo struct {
	list   []string
	sorted bool
}

func NewSearchInfo(s []string) *SearchInfo {
	return &SearchInfo{
		list:   s,
		sorted: false,
	}
}

func (s *SearchInfo) Sort() {
	sort.Slice(s.list, func(i, j int) bool {
		s.sorted = true
		return s.list[i] <= s.list[j]
	})
}

func (s *SearchInfo) Contains(searchFor string) bool {
	if !s.sorted {
		s.Sort()
	}
	var pos int
	l := len(s.list)
	pos = sort.Search(l, func(i int) bool {
		return s.list[i] >= searchFor
	})
	return pos < l && s.list[pos] == searchFor

}
