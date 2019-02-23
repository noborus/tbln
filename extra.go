package tbln

import "sort"

// Extra is table definition extra struct.
type Extra struct {
	value      interface{}
	hashTarget bool
}

// NewExtra is return new extra struct.
func NewExtra(value interface{}) Extra {
	return Extra{
		value:      value,
		hashTarget: false,
	}
}

func sortExt(ext map[string]Extra) list {
	l := list{}
	for key, extra := range ext {
		e := entry{key, extra}
		l = append(l, e)
	}
	sort.Sort(l)
	return l
}

type entry struct {
	n string
	v Extra
}
type list []entry

func (l list) Len() int {
	return len(l)
}
func (l list) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
func (l list) Less(i, j int) bool {
	return l[i].n < l[j].n
}
