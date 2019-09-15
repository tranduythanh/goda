package goda

import "sort"

// List ...
type List []interface{}

// ListString ...
type ListString []string

// UniqueString ...
func (l List) UniqueString() (s ListString) {
	var m = map[string]bool{}

	for _, el := range l {
		str := el.(string)

		if _, ok := m[str]; ok {
			continue
		}

		m[str] = true
		s = append(s, str)
	}
	return s
}

// SortString ...
func (l List) SortString() (s ListString) {
	for _, el := range l {
		s = append(s, el.(string))
	}
	sort.Strings(s)
	return s
}

// Sort ...
func (l ListString) Sort() (s ListString) {
	for _, el := range l {
		s = append(s, el)
	}
	sort.Strings(s)
	return s
}
