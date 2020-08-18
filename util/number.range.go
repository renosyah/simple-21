package util

func GetNumberBetweenRange(num, min, step, max int) int {
	lNum := min
	for i := min; i <= max; i += step {
		if num >= lNum && num <= i {
			return i
		}
		lNum = i
	}
	return num
}
