package polling

type Filter interface {
	Check(event interface{}) bool
}

// Built-in filter
type f struct {
	checkFn func(interface{}) bool
}

func (f f) Check(event interface{}) bool {
	return f.checkFn(event)
}

func makeDefaultFilterFromFunc(checkFn func(interface{}) bool) Filter {
	return f{checkFn: checkFn}
}
