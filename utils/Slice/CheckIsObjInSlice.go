package Slice

func CheckIsStringInSlice(s string, list []string) bool {
	for _, k := range list {
		if k == s {
			return true
		}
	}
	return false
}
