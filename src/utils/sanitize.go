package utils

func SanitizeString(str string) string {
	var sanStr string
	for _, char := range str {
		if len(string(char)) < 4 {
			sanStr = sanStr + string(char)
		}
	}
	return sanStr
}
