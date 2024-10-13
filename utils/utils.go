package utils

func In(s string, arr []string) int {
	for idx, str := range arr {
		if str == s {
			return idx
		}
	}
	return -1
}

func HasPrefix(prefix, str string) bool {
	for idx, char := range prefix {
		if idx >= len(str) || char != rune(str[idx]) {
			return false
		}
	}
	return true
}

func Split(s string, sep string) []string {
	res := []string{}

	if len(s) == 0 {
		return append(res, "")
	}
	if sep == "" {
		for _, char := range s {
			res = append(res, string(char))
		}
		return res
	}

	start := 0
	for idx := 0; idx < len(s); idx++ {
		if idx+len(sep) < len(s) && sep == string(s[idx:idx+len(sep)]) {
			res = append(res, s[start:idx])
			start = idx + len(sep)
			idx += len(sep) - 1
		} else if idx == len(s)-1 {
			res = append(res, s[start:])
		}
	}

	return res
}

func IsNumeric(s string) bool {
	for _, char := range s {
		if char > '9' || char < '0' {
			return false
		}
	}
	return true
}
