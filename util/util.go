package util

func PadRight(str string, length int) string {
	for {
		if len(str) >= length {
			return str
		}
		str = str + string(' ')
	}
}
