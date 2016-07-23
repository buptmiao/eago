package eago

func AssertEqual(v bool) {
	if !v {
		panic("not equal!")
	}
}

//Note:
// there is something trap in comparing interface{} and nil
// for detail: http://golang.org/doc/go_faq.html#nil_error
func AssertErrNil(v interface{}) {
	if v != nil {
		panic("not Nil")
	}
}

func AssertNotNil(v interface{}) {
	if v == nil {
		panic("Nil")
	}
}
