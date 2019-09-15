package goda

import "sort"

// List ...
type List []interface{}

// UniqueString ...
func (l List) UniqueString() (s []string) {
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
func (l List) SortString() (s []string) {
	for _, el := range l {
		s = append(s, el.(string))
	}
	sort.Strings(s)
	return s
}
