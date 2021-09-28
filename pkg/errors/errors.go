package errors

type (
	Op string
	// each request is assigned a unique string
	// for distributed tracing.
	ReqID    string
	Severity uint
)

const (
	SeverityError Severity = iota + 1
	SeverityWarn
	SeverityInfo
	SeverityDebug
)

type err struct {
	e        error
	op       Op
	code     Code
	reqID    ReqID
	severity Severity
}

// Error : to support in-build error
// interface.
func (e *err) Error() string {
	return e.e.Error()
}

// Construct a new stack based error.
func E(args ...interface{}) error {
	e := &err{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case error:
			e.e = arg
		case Op:
			e.op = arg
		case Code:
			e.code = arg
		case ReqID:
			e.reqID = arg
		case Severity:
			e.severity = arg
		default:
			panic("bad call to error constructor")
		}
	}

	return e
}

// Ops : return list of operation ids from error stack.
func Ops(e error) []Op {
	er, ok := e.(*err)
	if !ok {
		return nil
	}
	out := []Op{er.op}
	subErr, ok := er.e.(*err)
	if !ok {
		return out
	}
	out = append(out, Ops(subErr)...)

	return out
}

func ErrCode(e error) Code {
	er, ok := e.(*err)
	if !ok {
		return CodeUnexpected
	}
	if er.code > 0 {
		return er.code
	}

	return ErrCode(er.e)
}

func ErrReqID(e error) ReqID {
	er, ok := e.(*err)
	if !ok {
		return ""
	}
	if er.reqID != "" {
		return er.reqID
	}

	return ErrReqID(er.e)
}

func ErrSeverity(e error) Severity {
	er, ok := e.(*err)
	if !ok {
		return SeverityError
	}
	if er.severity != 0 {
		return er.severity
	}

	return ErrSeverity(er.e)
}
