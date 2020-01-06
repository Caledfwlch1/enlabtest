package tools

func IIF(cond bool, outTrue, outFalse interface{}) interface{} {
	if cond {
		return outTrue
	}
	return outFalse
}
