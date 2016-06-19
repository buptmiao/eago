package eago

func AssertEqual(v bool) {
	if !v {
		panic("not equal!")
	}
}

func AssertNil(v interface{}) {
	if v != nil {
		panic("not Nil")
	}
}

func AssertNotNil(v interface{}) {
	if v == nil {
		panic("Nil")
	}
}
