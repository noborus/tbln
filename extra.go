package tbln

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
