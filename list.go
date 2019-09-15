package goda

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
