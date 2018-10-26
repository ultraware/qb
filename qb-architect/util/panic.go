package util

// PanicOnErr panics if err is not nil
func PanicOnErr(err interface{}) {
	if err != nil {
		panic(err)
	}
}
