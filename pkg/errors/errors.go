package errors

type Op string

type Err struct {
	op   Op
	code Code
	err  error
}

func (e *Err) Error() string {
	return e.err.Error()
}

func E(args ...interface{}) error {
	e := &Err{}

	for _, arg := range args {
		switch arg := arg.(type) {
		case error:
			e.err = arg
		case Op:
			e.op = arg
		case Code:
			e.code = arg
		default:
			panic("bad call to error constructor")
		}
	}

	return e
}

func Ops(err error) []Op {
	e, ok := err.(*Err)
	if !ok {
		return nil
	}
	out := []Op{e.op}
	if _, ok := e.err.(*Err); !ok {
		return out
	}
	out = append(out, Ops(e.err)...)

	return out
}

// return top most error status code present on
// error stack.
func ErrCode(err error) Code {
	e, ok := err.(*Err)
	if !ok {
		return CodeInternal
	}
	if e.code > 0 {
		return e.code
	}

	return ErrCode(e.err)
}
